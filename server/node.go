package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"node/models"
	"time"

	"github.com/gorilla/mux"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MsgOk = []byte(`{"msg": "ok"}`)

func mergeResponses(nodeResponse, readingsResponse []byte) (string, error) {
	nodeJson := string(nodeResponse)

	readings := gjson.ParseBytes(readingsResponse).Array()
	// fmt.Println("readings", readings)
	// fmt.Println("nodes", nodeJson)

	nodes := gjson.Parse(nodeJson).Array()
	for i := 0; i < len(nodes); i++ {
		var err error
		q := fmt.Sprintf(`%d.`, i)
		nodeJson, err = sjson.Set(nodeJson, q+"reading", []interface{}{})
		if err != nil {
			log.Println(err)
			return "", err
		}
		nodeJson, err = sjson.Set(nodeJson, q+"datetime", 0)
		if err != nil {
			log.Println(err)
			return "", err
		}
	}

	for _, reading := range readings {
		var err error
		uid := reading.Get("uid").String()
		fmt.Println("uid:", uid)
		uidQuery := fmt.Sprintf(`#[uid=="%s"].reading`, uid)

		nodeJson, err = sjson.Set(nodeJson, uidQuery, reading.Get("values").Value())
		if err != nil {
			log.Println(err)
			return "", err
		}
		uidQuery = fmt.Sprintf(`#[uid=="%s"].datetime`, uid)
		nodeJson, err = sjson.Set(nodeJson, uidQuery, reading.Get("datetime").Value())
		if err != nil {
			log.Println(err)
			return "", err
		}
		fmt.Println(nodeJson)
	}
	return nodeJson, nil

}

func GetNodes(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("claims").(map[string]any)
	if !ok {
		err := errors.New("invalid claims found in context")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	username, ok := claims["username"].(string)
	if username == "" || !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(NewErrorResponse("invalid username").ByteResponse())
		return
	}

	designation, ok := claims["designation"].(string)
	if designation == "" || !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(NewErrorResponse("invalid or empty designation value").ByteResponse())
		return
	}

	privilege := getPrivilegeLevel(designation)
	if privilege > 5 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(NewErrorResponse("invalid or empty designation value").ByteResponse())
		return
	}

	var search bson.M
	switch privilege {
	case 0, 1, 3:
		search = bson.M{
			"isArchived": false,
		}
	case 2:
		search = bson.M{
			"isArchived": false,
			"user":       username,
		}
	case 4:
		search = bson.M{
			"isArchived": false,
			"createdBy":  username,
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write(NewErrorResponse("invalid user").ByteResponse())
		return
	}

	nodes, err := models.Get(r.Context(), models.NODE_COLLECTION, search, &options.FindOptions{})
	if err != nil {
		log.Println("error while gettings nodes: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(ErrorToResponse(err))
		return
	}

	nodeResponse, err := json.Marshal(nodes)
	if err != nil {
		log.Println("error while marshallilng nodes: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(ErrorToResponse(err))
		return
	}

	// fetch readings
	readings, err := models.Aggregate(r.Context(), models.READING_COLLECTION, bson.A{
		// bson.M{
		// 	"$match": bson.M{"user": username},
		// },
		bson.M{
			"$group": bson.M{
				"_id":    "$uid",
				"date":   bson.M{"$last": "$datetime"},
				"values": bson.M{"$last": "$values"},
			},
		},
		bson.M{
			"$project": bson.M{
				"uid":      "$_id",
				"datetime": "$date",
				"values":   "$values",
			},
		},
		bson.M{
			"$sort": bson.M{
				"datetime": -1,
			},
		},
	}, &options.AggregateOptions{})
	if err != nil {
		log.Println("error while getting readings: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(ErrorToResponse(err))
		return
	}

	readingsResponse, err := json.Marshal(readings)
	if err != nil {
		log.Println("error while marshalling readings: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(ErrorToResponse(err))
		return
	}
	finalResponse, err := mergeResponses(nodeResponse, readingsResponse)
	if err != nil {
		log.Println("error while merging readings: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(ErrorToResponse(err))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(finalResponse))
	return
}

func GetArchivedNodes(w http.ResponseWriter, r *http.Request) {
	search := bson.M{
		"isArchived": true,
	}
	nodes, err := models.Get(r.Context(), models.NODE_COLLECTION, search, &options.FindOptions{})
	if err != nil {
		log.Println("error while gettings nodes: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(ErrorToResponse(err))
		return
	}

	nodeResponse, err := json.Marshal(nodes)
	if err != nil {
		log.Println("error while marshallilng nodes: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(ErrorToResponse(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(nodeResponse)
	return
}

func AddNode(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("claims").(map[string]any)
	if !ok {
		err := errors.New("invalid claims found in context")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	username, ok := claims["username"].(string)
	if username == "" || !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(NewErrorResponse("invalid username").ByteResponse())
		return
	}

	designation, ok := claims["designation"].(string)
	if designation == "" || !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(NewErrorResponse("invalid or empty designation value").ByteResponse())
		return
	}

	privilege := getPrivilegeLevel(designation)
	if privilege > 5 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(NewErrorResponse("invalid or empty designation value").ByteResponse())
		return
	}

	if privilege > 2 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(NewErrorResponse("unauthorized to add node").ByteResponse())
		return
	}

	nodeContent, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("error while reading request: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(ErrorToResponse(err))
		return
	}

	var node models.Node
	err = json.Unmarshal(nodeContent, &node)
	if err != nil {
		log.Println("error while unmarhsaling node data: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(ErrorToResponse(err))
		return
	}

	now := time.Now().Unix()
	node.User = username
	node.CreatedBy = username
	node.ModifiedBy = username
	node.CreatedOn = now
	node.ModifiedOn = now

	_, err = models.Save(r.Context(), node)
	if err != nil {
		log.Println("could not save node: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(ErrorToResponse(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(MsgOk)
	return
}

func ModifyNode(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("could not read node body: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(ErrorToResponse(err))
		return
	}

	var node map[string]any
	err = json.Unmarshal(body, &node)
	if err != nil {
		log.Println("could not decode node body: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(ErrorToResponse(err))
		return
	}

	uid, ok := node["uid"].(string)
	if uid == "" || !ok {
		err := errors.New("uid does not exist in request body")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(ErrorToResponse(err))
		return
	}

	search := bson.M{"uid": uid}
	update := bson.M{"$set": node}

	err = models.Update(r.Context(), models.NODE_COLLECTION, search, update)
	if err != nil {
		log.Println("could not update node: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(ErrorToResponse(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(MsgOk)
	return
}

func DeleteNode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	uid, ok := vars["uid"]
	if uid == "" || !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(NewErrorResponse("invalid uid").ByteResponse())
		return
	}

	search := bson.M{"uid": uid}
	deletedCount, err := models.Delete(r.Context(), models.NODE_COLLECTION, search)
	if err != nil {
		log.Println("could not delete node: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(ErrorToResponse(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"deleteCount": %d}`, deletedCount)))
	return
}

func GetNodeByUid(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	uid, ok := vars["uid"]
	if uid == "" || !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(NewErrorResponse("invalid uid").ByteResponse())
		return
	}

	search := bson.M{"uid": uid}
	result, err := models.Get(r.Context(), models.NODE_COLLECTION, search, &options.FindOptions{})
	if err != nil {
		log.Println("could not delete node: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(ErrorToResponse(err))
		return
	}

	node, err := json.Marshal(result)
	if err != nil {
		log.Println("could not marshal node: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(ErrorToResponse(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(node)
	return
}
