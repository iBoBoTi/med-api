package services

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/decagonhq/meddle-api/config"
	"github.com/decagonhq/meddle-api/db"
	"github.com/decagonhq/meddle-api/errors"
	"github.com/decagonhq/meddle-api/models"
	"github.com/go-co-op/gocron"
)

//go:generate mockgen -destination=../mocks/medication_mock.go -package=mocks github.com/decagonhq/meddle-api/services MedicationService

type MedicationService interface {
	CreateMedication(request *models.MedicationRequest) (*models.MedicationResponse, *errors.Error)
	GetNextMedications(userID uint) ([]models.MedicationResponse, *errors.Error)
	GetMedicationDetail(id uint, userId uint) (*models.MedicationResponse, *errors.Error)
	GetAllMedications(userID uint) ([]models.MedicationResponse, *errors.Error)
	CronUpdateMedicationForNextTime() error
	UpdateMedication(request *models.UpdateMedicationRequest, medicationID uint, userID uint) *errors.Error
	FindMedication(medicationName string, by string, purpose string, duration int, dosage int) (*[]models.Medication, error)
}

// medicationService struct
type medicationService struct {
	Config                *config.Config
	medicationRepo        db.MedicationRepository
	medicationHistoryRepo db.MedicationHistoryRepository
}

// NewMedicationService instantiate an authService
func NewMedicationService(medicationRepo db.MedicationRepository, medicationHistoryRepo db.MedicationHistoryRepository, conf *config.Config) MedicationService {
	return &medicationService{
		Config:                conf,
		medicationRepo:        medicationRepo,
		medicationHistoryRepo: medicationHistoryRepo,
	}
}

func (m *medicationService) CreateMedication(request *models.MedicationRequest) (*models.MedicationResponse, *errors.Error) {
	startDate, err := time.Parse(time.RFC3339, request.MedicationStartDate)
	if err != nil {
		return nil, errors.New("wrong date format", http.StatusBadRequest)
	}
	startTime, err := time.Parse(time.RFC3339, request.MedicationStartTime)
	if err != nil {
		return nil, errors.New("wrong time format", http.StatusBadRequest)
	}

	medication := request.ReqToMedicationModel()
	medication.CreatedAt = time.Now().Unix()
	medication.UpdatedAt = time.Now().Unix()
	medication.MedicationStartDate = startDate
	medication.MedicationStartTime = startTime
	var nextTime time.Time
	if medication.MedicationStartTime.Unix() > time.Now().Unix() {
		nextTime = medication.MedicationStartTime
	} else {
		nextTime = medication.MedicationStartTime.Add(time.Hour * time.Duration(medication.TimeInterval))
	}

	medication.MedicationStopDate = medication.MedicationStartTime.AddDate(0, 0, medication.Duration)
	medication.NextDosageTime = GetNextDosageTime(nextTime, medication.MedicationStartTime)

	response, err := m.medicationRepo.CreateMedication(medication)
	if err != nil {
		return nil, errors.ErrInternalServerError
	}
	return response.MedicationToResponse(), nil
}

func (m *medicationService) GetMedicationDetail(id uint, userId uint) (*models.MedicationResponse, *errors.Error) {
	medic, err := m.medicationRepo.GetMedicationDetail(id, userId)
	if err != nil {
		return nil, errors.ErrInternalServerError
	}
	return medic.MedicationToResponse(), nil
}

func (m *medicationService) GetAllMedications(userID uint) ([]models.MedicationResponse, *errors.Error) {
	var medicationResponses []models.MedicationResponse

	medications, err := m.medicationRepo.GetAllMedications(userID)
	if err != nil {
		return nil, errors.ErrInternalServerError
	}

	for _, medication := range medications {
		medicationResponses = append(medicationResponses, *medication.MedicationToResponse())
	}
	return medicationResponses, nil
}

func (m *medicationService) UpdateMedication(request *models.UpdateMedicationRequest, medicationID uint, userID uint) *errors.Error {
	startDate, err := time.Parse(time.RFC3339, request.MedicationStartDate)
	if err != nil {
		return errors.New("wrong date format", http.StatusBadRequest)
	}
	startTime, err := time.Parse(time.RFC3339, request.MedicationStartTime)
	if err != nil {
		return errors.New("wrong time format", http.StatusBadRequest)
	}
	medication := models.Medication{
		Name:                   request.Name,
		Dosage:                 request.Dosage,
		TimeInterval:           request.TimeInterval,
		Duration:               request.Duration,
		MedicationPrescribedBy: request.MedicationPrescribedBy,
		PurposeOfMedication:    request.PurposeOfMedication,
		MedicationIcon:         request.MedicationIcon,
		MedicationStartDate:    startDate,
		MedicationStartTime:    startTime,
	}

	nextTime := medication.MedicationStartTime.Add(time.Hour * time.Duration(medication.TimeInterval))
	medication.MedicationStopDate = medication.MedicationStartTime.AddDate(0, 0, medication.Duration)

	medication.NextDosageTime = GetNextDosageTime(nextTime, medication.MedicationStartTime)

	//get medication where user and medication id is defined above then send it for updating
	err = m.medicationRepo.UpdateMedication(&medication, medicationID, userID)
	if err != nil {
		return errors.ErrInternalServerError
	}
	return nil
}

func (m *medicationService) GetNextMedications(userID uint) ([]models.MedicationResponse, *errors.Error) {
	var nextMedicationResponses []models.MedicationResponse

	medications, err := m.medicationRepo.GetNextMedications(userID)
	if err != nil {
		return nil, errors.ErrInternalServerError
	}

	for _, medication := range medications {
		nextMedicationResponses = append(nextMedicationResponses, *medication.MedicationToResponse())
	}
	return nextMedicationResponses, nil

}

func (m *medicationService) CronUpdateMedicationForNextTime() error {
	medications, err := m.medicationRepo.GetAllNextMedicationsToUpdate()
	if err != nil {
		return fmt.Errorf("could not get next medications while running update next dosage cron job")
	}

	//create medication history for each medication
	if medications != nil {
		go m.CreateMedicationHistory(medications)
	}

	for _, medication := range medications {
		timeSumation := medication.NextDosageTime.Add(time.Hour * time.Duration(medication.TimeInterval))
		nextDosageTime := GetNextDosageTime(timeSumation, medication.NextDosageTime)

		if nextDosageTime.Unix() < medication.MedicationStopDate.Unix() {
			err = m.medicationRepo.UpdateNextMedicationTime(&medication, nextDosageTime)
			if err != nil {
				return fmt.Errorf("could not update next medication time while running update next dosage cron job")
			}
		} else {
			err = m.medicationRepo.UpdateMedicationDone(&medication)
			if err != nil {
				return fmt.Errorf("could not update is medication done while running update next dosage cron job")
			}
		}
	}
	return nil
}

func UpdateMedicationCronJob(medicationService MedicationService) {
	// _, presentMinute, _ := time.Now().UTC().Clock()
	// if presentMinute%15 != 0 {
	// 	time.Sleep(time.Duration(presentMinute+(presentMinute%15)) * time.Minute)
	// }
	// waitTime := time.Duration(60-presentMinute)*time.Minute + time.Duration(60-presentSecond)*time.Second
	s := gocron.NewScheduler(time.UTC)
	s.Every(1).Minute().Do(func() {
		err := medicationService.CronUpdateMedicationForNextTime()
		if err != nil {
			log.Printf("cron job error: %v", err)
		}
	})
	s.StartBlocking()
}

func GetNextDosageTime(t1, t2 time.Time) time.Time {
	if t1.Day()-t2.Day() <= 0 {
		return time.Date(t1.Year(), t1.Month(), t1.Day(), t1.Hour(), t1.Minute(), 0, 0, time.UTC)
	}
	return time.Date(t2.Year(), t2.Month(), t2.Day()+1, 9, 0, 0, 0, time.UTC)
}

func (m *medicationService) CreateMedicationHistory(medications []models.Medication) {
	for _, medication := range medications {
		medicationHistory := models.NewMedicationHistory(medication)
		_, err := m.medicationHistoryRepo.CreateMedicationHistory(medicationHistory)
		if err != nil {
			log.Printf("error creating medication history for %v for %v : %v", medication.ID, medication.NextDosageTime, err)
		}
	}
}

func (m *medicationService) FindMedication(medicationName, by, purpose string, duration int, dosage int) (*[]models.Medication, error) {
	var medicationResponses []models.MedicationResponse
	medications, err := m.medicationRepo.FindMedication(medicationName, by, purpose, duration, dosage)
	if err != nil {
		return nil, errors.ErrInternalServerError
	}
	for _, medication := range *medications {
		medicationResponses = append(medicationResponses, *medication.MedicationToResponse())
	}
	return medications, nil
}
