package services

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/decagonhq/meddle-api/config"
	"github.com/decagonhq/meddle-api/db"
	"github.com/decagonhq/meddle-api/errors"
	"github.com/decagonhq/meddle-api/models"
	"github.com/go-co-op/gocron"
	"google.golang.org/api/option"
	"log"
	"strconv"
	"time"
)

//go:generate mockgen -destination=../mocks/auth_mock.go -package=mocks github.com/decagonhq/meddle-api/services PushNotification

type PushNotifier interface {
	AuthorizeNotification(request *models.AddNotificationTokenArgs) (*models.FCMNotificationToken, *errors.Error)
	CheckIfThereIsNextMedication()
	SendPushNotification(registrationTokens []string, payload *models.PushPayload) (*messaging.MulticastMessage, *errors.Error)
	NotificationsCronJob()
	GetSingleUserDeviceTokens(userId int) ([]string, *errors.Error)
}

type notificationService struct {
	Conf             *config.Config
	notificationRepo db.NotificationRepository
	Client           *messaging.Client
}

// NewFirebaseCloudMessaging instantiates an FCM service
func NewFirebaseCloudMessaging(notificationRepo db.NotificationRepository, conf *config.Config) (PushNotifier, error) {
	firebaseApp, err := firebase.NewApp(context.Background(), nil, option.WithCredentialsFile(conf.GoogleApplicationCredentials))
	if err != nil {
		log.Println(err)
		return &notificationService{}, err
	}

	var fcm notificationService

	fcm.Client, err = firebaseApp.Messaging(context.Background())
	if err != nil {
		log.Println(err)
		return &notificationService{}, err
	}

	return &notificationService{
		notificationRepo: notificationRepo,
		Conf:             conf,
		Client:           fcm.Client,
	}, nil
}

func (fcm *notificationService) AuthorizeNotification(request *models.AddNotificationTokenArgs) (*models.FCMNotificationToken, *errors.Error) {
	token, err := fcm.notificationRepo.AddNotificationToken(request)
	if err != nil {
		return nil, errors.ErrInternalServerError
	}
	return token, nil
}

func (fcm *notificationService) GetSingleUserDeviceTokens(userid int) ([]string, *errors.Error) {
	tokens, err := fcm.notificationRepo.GetSingleUserDeviceTokens(userid)
	if err != nil {
		return nil, errors.ErrInternalServerError
	}
	return tokens, nil
}

// CheckIfThereIsNextMedication cron job
//check all currently due medication in db
func (fcm *notificationService) CheckIfThereIsNextMedication() {
	medicationNotifications, err := fcm.notificationRepo.GetAllNextMedicationsToSendNotifications()
	if err != nil {
		log.Println("could not get medications from db", err)
		return
	}

	//check db for all the time of notifications
	for i := 0; i < len(medicationNotifications); i++ {
		go func(i int) {
			userId := medicationNotifications[i].UserID
			deviceTokens, err := fcm.notificationRepo.GetSingleUserDeviceTokens(int(userId))
			if err != nil {
				log.Printf("error retrieving device notification tokens: %v\n", err)
				return
			}

			if len(deviceTokens) == 0 {
				log.Printf("empty token list: %v\n", err)
				return
			}

			notification, err := fcm.SendPushNotification(deviceTokens, &models.PushPayload{
				Body:  "'" + (medicationNotifications)[i].Name + "' is due in___",
				Title: (medicationNotifications)[i].Name,
				Data: map[string]string{
					"link": "/user/medication/id?=" + strconv.Itoa(int((medicationNotifications)[i].ID)),
				},
				ClickAction: "/user/medication/id?=" + strconv.Itoa(int((medicationNotifications)[i].ID)),
			})
			if err != nil {
				log.Println("error sending notification", err)
				return
			}

			log.Println("logging notifications", notification)
		}(i)

	}
}

func (fcm *notificationService) SendPushNotification(registrationTokens []string, payload *models.PushPayload) (*messaging.MulticastMessage, *errors.Error) {

	notification := &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title:    payload.Title,
			Body:     payload.Body,
			ImageURL: "https://imgur.com/a/hmt6Mx2",
		},
		Data: payload.Data,
		Webpush: &messaging.WebpushConfig{
			Data: payload.Data,
			Notification: &messaging.WebpushNotification{
				Title:   payload.Title,
				Body:    payload.Body,
				Icon:    "https://imgur.com/a/hmt6Mx2",
				Vibrate: []int{200, 100, 200},
			},
		},
		Android: &messaging.AndroidConfig{
			Notification: &messaging.AndroidNotification{
				Color:                 "#4C51BF",
				ClickAction:           payload.ClickAction,
				DefaultSound:          true,
				DefaultVibrateTimings: true,
				DefaultLightSettings:  true,
			},
			Data: payload.Data,
		},
		APNS: &messaging.APNSConfig{
			Headers: map[string]string{
				"apns-priority": "10",
			},
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					Alert: &messaging.ApsAlert{
						Title: payload.Title,
						Body:  payload.Body,
					},
					Sound:            "default",
					ContentAvailable: true,
				},
			},
			FCMOptions: nil,
		},
		Tokens: registrationTokens,
	}

	_, err := fcm.Client.SendMulticast(context.Background(), notification)
	if err != nil {
		log.Fatalln(err)
		return &messaging.MulticastMessage{}, errors.ErrInternalServerError
	}

	return notification, nil
}

func (fcm *notificationService) NotificationsCronJob() {
	_, presentMinute, presentSecond := time.Now().UTC().Clock()
	waitTime := time.Duration(60-presentMinute)*time.Minute + time.Duration(60-presentSecond)*time.Second
	scheduler := gocron.NewScheduler(time.UTC)
	time.Sleep(waitTime)
	scheduler.Every(1).Hour().Do(func() {
		fcm.CheckIfThereIsNextMedication()
	})
	scheduler.StartBlocking()
}
