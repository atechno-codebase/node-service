package service

import (
	"log"
	"node/config"
	"node/models"
)

func Init() {
	models.Init(config.Configuration.MongoUrl, config.Configuration.DatabaseName)
	log.Println("service initialised")
}
