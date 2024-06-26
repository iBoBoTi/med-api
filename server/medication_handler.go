package server

import (
	"github.com/decagonhq/meddle-api/errors"
	"github.com/decagonhq/meddle-api/models"
	"github.com/decagonhq/meddle-api/server/response"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func (s *Server) handleCreateMedication() gin.HandlerFunc {
	return func(c *gin.Context) {
		var medicationRequest models.MedicationRequest
		_, user, err := GetValuesFromContext(c)
		if err != nil {
			err.Respond(c)
			return
		}
		userId := user.ID
		if err := decode(c, &medicationRequest); err != nil {
			response.JSON(c, "", http.StatusBadRequest, nil, err)
			return
		}
		medicationRequest.UserID = userId
		createdMedication, err := s.MedicationService.CreateMedication(&medicationRequest)
		if err != nil {
			err.Respond(c)
			return
		}
		response.JSON(c, "medication created successful", http.StatusCreated, createdMedication, nil)
	}
}

func (s *Server) handleGetMedDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, user, err := GetValuesFromContext(c)
		if err != nil {
			err.Respond(c)
			return
		}
		id := c.Param("id")
		userId, errr := strconv.ParseUint(id, 10, 32)
		if errr != nil {
			response.JSON(c, "error parsing id", http.StatusBadRequest, nil, errr)
			return
		}
		medication, err := s.MedicationService.GetMedicationDetail(uint(userId), user.ID)
		if err != nil {
			response.JSON(c, "", http.StatusInternalServerError, nil, errors.New("internal server error", http.StatusInternalServerError))
			return
		}
		response.JSON(c, "retrieved medications successfully", http.StatusOK, gin.H{"medication": medication}, nil)
	}
}

func (s *Server) handleGetAllMedications() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, user, err := GetValuesFromContext(c)
		if err != nil {
			err.Respond(c)
			return
		}
		medications, err := s.MedicationService.GetAllMedications(user.ID)
		if err != nil {
			err.Respond(c)
			return
		}
		response.JSON(c, "medications retrieved successfully", http.StatusOK, medications, nil)
	}
}

func (s *Server) handleGetNextMedication() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, user, err := GetValuesFromContext(c)
		if err != nil {
			err.Respond(c)
			return
		}

		medication, err := s.MedicationService.GetNextMedications(user.ID)
		if err != nil {
			err.Respond(c)
			return
		}
		response.JSON(c, "medication retrieved successfully", http.StatusOK, medication, nil)
	}
}

func (s *Server) handleUpdateMedication() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, user, err := GetValuesFromContext(c)
		if err != nil {
			err.Respond(c)
			return
		}
		medicationID, errr := strconv.ParseUint(c.Param("medicationID"), 10, 32)
		if errr != nil {
			response.JSON(c, "invalid ID", http.StatusBadRequest, nil, errr)
			return
		}
		var updateMedicationRequest models.UpdateMedicationRequest
		if err := decode(c, &updateMedicationRequest); err != nil {
			response.JSON(c, "", http.StatusBadRequest, nil, err)
			return
		}
		err = s.MedicationService.UpdateMedication(&updateMedicationRequest, uint(medicationID), user.ID)
		if err != nil {
			err.Respond(c)
			return
		}
		response.JSON(c, "medication updated successfully", http.StatusOK, nil, nil)
	}
}

func (s *Server) handleFindMedication() gin.HandlerFunc {
	return func(c *gin.Context) {

		medicationName := c.Query("name")
		medicationDosage := c.Query("dosage")
		medicationDuration := c.Query("duration")
		medicationPrescribedBy := c.Query("medication_prescribed_by")
		medicationPurpose := c.Query("purpose_of_medication")

		dosage,_ := strconv.Atoi(medicationDosage)
		medDuration,_ := strconv.Atoi(medicationDuration)

		medications, err := s.MedicationService.FindMedication(medicationName, medicationPrescribedBy, medicationPurpose, medDuration, dosage)
		if err != nil {
			log.Printf("Error1: %v", err.Error())
			response.JSON(c, "", http.StatusInternalServerError, nil, errors.New("internal server error", http.StatusInternalServerError))
			return
		}
		response.JSON(c, "medications retrieved successfully", http.StatusOK, medications, nil)
	}
}

