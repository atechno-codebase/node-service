package service

import (
	"context"
	"errors"
	"log"
	"node/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetReadingsByUid(ctx context.Context, uid string) ([]models.Reading, error) {
	dbResult, err := models.Get(ctx, models.READING_COLLECTION, bson.M{
		"uid": uid,
	}, &options.FindOptions{
		Sort: bson.M{"datetime": -1},
	})

	if err != nil {
		log.Println("error while fetching readings from DB: ", err)
		return nil, err
	}

	readings, ok := dbResult.([]models.Reading)
	if !ok {
		err := errors.New("error while casting to readings")
		log.Println("error while fetching readings from DB: ", err)
		return nil, err

	}

	return readings, nil
}
