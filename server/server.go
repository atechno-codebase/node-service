package server

import (
	"log"
	"net/http"
	"node/config"

	"github.com/gorilla/mux"
)

var PORT string

func Init() {
	PORT = config.Get("port").String()
	r := mux.NewRouter()

	initRoutes(r)

	log.Printf("started server on port: %s\n", PORT)
	log.Fatalln(http.ListenAndServe(PORT, r))
}

func initRoutes(r *mux.Router) {

	r.Use(authMiddleware)

	r.HandleFunc("/", GetNodes).
		Methods(http.MethodGet)

	r.HandleFunc("/readings/{uid}", GetReadingForUid).
		Methods(http.MethodGet)

	r.HandleFunc("/archived", GetArchivedNodes).
		Methods(http.MethodGet)

	r.HandleFunc("/add", AddNode).
		Methods(http.MethodPost)

	r.HandleFunc("/modify", ModifyNode).
		Methods(http.MethodPost)

	r.HandleFunc("/{uid}", DeleteNode).
		Methods(http.MethodDelete)

	r.HandleFunc("/readings/all/", GetAllReadings).
		Methods(http.MethodPost)

	r.HandleFunc("/search/{uid}", GetNodeByUid).
		Methods(http.MethodGet)
}
