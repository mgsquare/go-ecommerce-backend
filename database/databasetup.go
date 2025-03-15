package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBSet() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOption := options.Client().ApplyURI("mongodb://localhost:27017")

	client, err := mongo.Connect(ctx, clientOption)

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("could not connect to mongo db %v:", err)
	}

	fmt.Println("Succesfully connected to mongoDB")
	return client
}

var Client *mongo.Client = DBSet()

func UserData(client *mongo.Client, collectionName string) *mongo.Collection {

	var Collection *mongo.Collection = client.Database("Ecommerce").Collection(collectionName)
	return Collection

}

func ProductData(client *mongo.Client, collectionName string) *mongo.Collection {
	var ProductCollection *mongo.Collection = client.Database("Ecommerce").Collection(collectionName)
	return ProductCollection
}
