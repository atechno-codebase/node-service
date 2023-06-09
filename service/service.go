package service

import (
	"log"
	"node/config"
	"node/models"
)

func Init() {
	mongoUrl := config.Get("mongoUrl").String()
	dbName := config.Get("dbName").String()

	models.Init(mongoUrl, dbName)
	log.Println("service initialised")
}
