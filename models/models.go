package models

import (
	"context"
	"errors"
	"log"
	"node/database"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrUnknownCollection = errors.New("unknown type of collection or document")
)

var databaseName string

func Init(mongoUrl, dbName string) {
	databaseName = dbName
	database.Init(mongoUrl)
}

func getCollectionName(doc any) (string, error) {
	switch doc.(type) {
	case Reading, *Reading:
		return READING_COLLECTION, nil
	case Node, *Node:
		return NODE_COLLECTION, nil
	default:
		return "", ErrUnknownCollection
	}
}

func Save(ctx context.Context, document any) (any, error) {
	collectionName, err := getCollectionName(document)
	if err != nil {
		return nil, err
	}

	res, err := database.RunQuery(func(client *mongo.Client) (interface{}, error) {
		collection := client.Database(databaseName).Collection(collectionName)

		_, err := collection.InsertOne(ctx, document)
		if err != nil {
			log.Printf("error while inserting document for %s\n", collectionName)
			return nil, err
		}
		return nil, nil
	})

	return res, nil
}

func DecodeIntoNode(ctx context.Context, cursor *mongo.Cursor) ([]Node, error) {
	nodes := []Node{}
	err := cursor.All(ctx, &nodes)
	if err != nil {
		log.Printf("error while decoding into node\n")
		return nil, err
	}
	return nodes, nil
}

func DecodeIntoReading(ctx context.Context, cursor *mongo.Cursor) ([]Reading, error) {
	readings := []Reading{}
	err := cursor.All(ctx, &readings)
	if err != nil {
		log.Printf("error while decoding into reading\n")
		return nil, err
	}
	return readings, nil
}

func Get(ctx context.Context, collectionName string, search any, findOptions *options.FindOptions) (any, error) {
	res, err := database.RunQuery(func(client *mongo.Client) (interface{}, error) {
		var err error
		collection := client.Database(databaseName).Collection(collectionName)

		cursor, err := collection.Find(ctx, search, findOptions)
		if err != nil {
			log.Printf("error while finding document for %s\n", collectionName)
			return nil, err
		}

		var returnValue any
		switch collectionName {
		case NODE_COLLECTION:
			returnValue, err = DecodeIntoNode(ctx, cursor)
		case READING_COLLECTION:
			returnValue, err = DecodeIntoReading(ctx, cursor)
		default:
			return nil, ErrUnknownCollection
		}
		if err != nil {
			log.Printf("error while decoding document for %s\n", collectionName)
			return nil, err
		}

		return returnValue, nil
	})

	if err != nil {
		log.Printf("error while running database query for %s\n", collectionName)
		return nil, err
	}
	return res, err
}

func Update(ctx context.Context, collectionName string, search, update any) error {
	_, err := database.RunQuery(func(client *mongo.Client) (interface{}, error) {
		var err error
		collection := client.Database(databaseName).Collection(collectionName)

		_, err = collection.UpdateMany(ctx, search, update)
		if err != nil {
			log.Printf("error while deleting document for %s\n", collectionName)
			return nil, err
		}

		return nil, nil
	})
	if err != nil {
		log.Printf("error while running database query for %s\n", collectionName)
		return err
	}
	return err
}

func Delete(ctx context.Context, collectionName string, search any) (int64, error) {
	res, err := database.RunQuery(func(client *mongo.Client) (interface{}, error) {
		var err error
		collection := client.Database(databaseName).Collection(collectionName)

		delRes, err := collection.DeleteMany(ctx, search)
		if err != nil {
			log.Printf("error while deleting document for %s\n", collectionName)
			return nil, err
		}

		return delRes.DeletedCount, nil
	})
	if err != nil {
		log.Printf("error while running database query for %s\n", collectionName)
		return 0, err
	}
	return res.(int64), err
}

func Aggregate(ctx context.Context, collectionName string, search any, aggregateOptions *options.AggregateOptions) (any, error) {
	res, err := database.RunQuery(func(client *mongo.Client) (interface{}, error) {
		var err error
		collection := client.Database(databaseName).Collection(collectionName)

		cursor, err := collection.Aggregate(ctx, search, aggregateOptions)
		if err != nil {
			log.Printf("error while finding document for %s\n", collectionName)
			return nil, err
		}

		var returnValue any
		switch collectionName {
		case NODE_COLLECTION:
			returnValue, err = DecodeIntoNode(ctx, cursor)
		case READING_COLLECTION:
			returnValue, err = DecodeIntoReading(ctx, cursor)
		default:
			return nil, ErrUnknownCollection
		}
		if err != nil {
			log.Printf("error while decoding document for %s\n", collectionName)
			return nil, err
		}

		return returnValue, nil
	})

	if err != nil {
		log.Printf("error while running database query for %s\n", collectionName)
		return nil, err
	}
	return res, err
}
