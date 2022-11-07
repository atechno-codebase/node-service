package main

import (
	"log"
	"net/http"
	"node/config"
	"node/service"
)

var PORT string

func init() {
	config.Init()
	log.SetPrefix("node-service: ")
	service.Init()
	PORT = config.Get("port").String()
}

func main() {
	log.Printf("started server on port: %s\n", PORT)
	log.Fatalln(http.ListenAndServe(":"+PORT, nil))
}
