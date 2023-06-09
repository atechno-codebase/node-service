package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"node/models"
	"node/service"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const TIMEZONE_OFFSET = 19800

func getPrivilegeLevel(designation string) int {
	switch designation {
	case "superadmin":
		return 0
	case "admin":
		return 1
	case "user":
		return 2
	case "maintenance":
		return 3
	case "observer":
		return 4
	default:
		return 5
	}
}

func GetReadingForUid(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	uid, ok := vars["uid"]
	if uid == "" || !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(NewErrorResponse("invalid uid").ByteResponse())
		return
	}

	readings, err := service.GetReadingsByUid(r.Context(), uid)
	if err != nil {
		log.Println("failed to get readings from DB: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	response, err := json.Marshal(readings)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
	return
}

func GetAllReadings(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Header)
	bodyContent, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	var reqBody map[string]any
	err = json.Unmarshal(bodyContent, &reqBody)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	uid, ok := reqBody["uid"].(string)
	if uid == "" || !ok {
		err := errors.New("cannot decode uid from request")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	to, ok := reqBody["to"].(string)
	if to == "" || !ok {
		err := errors.New("cannot decode `to` from request")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	from, ok := reqBody["from"].(string)
	if from == "" || !ok {
		err := errors.New("cannot decode `from` from request")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	fromTs, err := strconv.ParseInt(from, 10, 64)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	toTs, err := strconv.ParseInt(to, 10, 64)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	toTime := time.Unix(toTs, 0)
	fromTime := time.Unix(fromTs, 0)

	// Manually adjust GMT+5:30 offset
	// fromTime = fromTime.Add(-time.Hour * 5).Add(-time.Minute * 30)
	// toTime = toTime.Add(-time.Hour * 5).Add(-time.Minute * 30)
	fmt.Println("getting all readings for", uid)
	log.Println(fromTime, fromTime.Unix())
	log.Println(toTime, toTime.Unix())
	result, err := models.Aggregate(r.Context(), models.READING_COLLECTION, bson.A{
		bson.M{
			"$match": bson.M{
				"$and": bson.A{
					bson.M{"uid": uid},
					bson.M{
						"datetime": bson.M{
							"$gte": fromTs,
							"$lte": toTs,
						},
					},
				},
			},
		},
		bson.M{
			"$sample": bson.M{
				"size": 500,
			},
		},
		bson.M{
			"$sort": bson.M{
				"datetime": 1,
			},
		},
	}, &options.AggregateOptions{})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	response, err := json.Marshal(result)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	log.Println(string(response))

	w.WriteHeader(http.StatusOK)
	w.Write(response)
	return

}
