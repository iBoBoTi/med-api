package main

import (
	"log"
	"net/http"
	"time"

	"github.com/decagonhq/meddle-api/config"
	"github.com/decagonhq/meddle-api/db"
	"github.com/decagonhq/meddle-api/server"
	"github.com/decagonhq/meddle-api/services"
)

func main() {
	http.DefaultClient.Timeout = time.Second * 10
	conf, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	gormDB := db.GetDB(conf)
	authRepo := db.NewAuthRepo(gormDB)
	mail := services.NewMailService(conf)
	authService := services.NewAuthService(authRepo, conf, mail)
	medicationRepo := db.NewMedicationRepo(gormDB)
	medicationService := services.NewMedicationService(medicationRepo, conf)

	s := &server.Server{
		Config:            conf,
		AuthRepository:    authRepo,
		AuthService:       authService,
		MedicationService: medicationService,
	}
	go services.UpdateMedicationCronJob(medicationService)
	s.Start()
}
