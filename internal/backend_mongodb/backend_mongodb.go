package backend_mongodb

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbName    = "counter"
	collName  = "counter"
	counterID = "counter"
)

func DoCountMongoDB(mongodbURI, hostname string) (int, error) {
	clientOptions := options.Client().ApplyURI(mongodbURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Error().
			Str("hostname", hostname).
			Msg(fmt.Sprintf("error=%s", err))
		return -1, err
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			log.Error().
				Str("hostname", hostname).
				Msg(fmt.Sprintf("error=%s", err))
		}
	}()

	collection := client.Database(dbName).Collection(collName)

	// Query the current counter value
	filter := bson.M{"_id": counterID}
	update := bson.M{"$inc": bson.M{"count": 1}} // Increment the counter

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var result bson.M
	err = collection.FindOneAndUpdate(context.TODO(), filter, update, opts).Decode(&result)
	if err != nil {
		log.Error().
			Str("hostname", hostname).
			Msg(fmt.Sprintf("error=%s", err))
		return -1, err
	}

	count := result["count"].(int32)
	return int(count), nil
}

func GetCountMongoDB(mongodbURI, hostname string) (int, error) {

	clientOptions := options.Client().ApplyURI(mongodbURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Error().
			Str("hostname", hostname).
			Msg(fmt.Sprintf("error=%s", err))
		return -1, err
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			log.Error().
				Str("hostname", hostname).
				Msg(fmt.Sprintf("error=%s", err))
		}
	}()

	collection := client.Database(dbName).Collection(collName)

	// Query the current counter value
	filter := bson.M{"_id": counterID}

	var result bson.M
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Error().
			Str("hostname", hostname).
			Msg(fmt.Sprintf("error=%s", err))
		return -1, err
	}

	count := result["count"].(int32)
	return int(count), nil
}
