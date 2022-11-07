package main

import (
	"company-service/config"
	"company-service/db"
	"company-service/server"
	"company-service/services"
	"log"
	"net/http"
	"time"
)

func main() {
	http.DefaultClient.Timeout = time.Second * 10
	conf, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	gormDB := &db.GormDB{}
	gormDB.Init(conf)
	companyRepo := db.NewAuthRepo(db.NewDB())
	kafka := db.NewKafkaRepo(gormDB)
	companyService := services.NewCompanyService(companyRepo, kafka, conf)
	s := &server.Server{
		Config:            conf,
		CompanyRepository: companyRepo,
		CompanyService:    companyService,
	}
	s.Start()
}
