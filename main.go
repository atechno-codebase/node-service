package main

import (
	"io"
	"log"
	"node/config"
	"node/server"
	"node/service"
	"os"
	"path/filepath"
)

func init() {
	config.Init()

	initLogging()
	log.SetPrefix("node-service: ")
}

func main() {
	service.Init()
	server.Init()
}

func initLogging() {
	logFolderPath := config.Get("log_path").String()
	logFilePath := filepath.Join(filepath.Clean(logFolderPath), "node-service.log")
	logFile, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	logDest := io.MultiWriter(logFile, os.Stdout)
	if err != nil {
		log.Println("could not open log folder path")
		return
	}
	log.SetOutput(logDest)
}
