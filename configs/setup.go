package configs

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(EnvMongoURI()))
	if err != nil {
		log.Fatal(err)
	}
	context, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(context)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connect to MongoDb")
	return client

}

var Db *mongo.Client = ConnectDB()

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("baruApi").Collection(collectionName)
	return collection
}
