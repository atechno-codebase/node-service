package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type QueryFunc func(client *mongo.Client) (interface{}, error)

var url string

func Init(mongoUrl string) {
	url = mongoUrl
}

func Connect(mongoUrl string) (*mongo.Client, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoUrl))
	if err != nil {
		fmt.Println("Error while creating mongo connection: ", err)
		return nil, err
	}
	return client, nil
}

func RunQuery(dbCallback QueryFunc) (interface{}, error) {
	client, err := Connect(url)
	if err != nil {
		fmt.Println("error while connecting: ", err)
		return nil, err
	}

	result, err := dbCallback(client)
	if err != nil {
		fmt.Println("error in callback: ", err)
		return nil, err
	}

	err = client.Disconnect(context.Background())
	if err != nil {
		fmt.Println("error while closing the connection: ", err)
		return nil, err
	}

	return result, nil
}
