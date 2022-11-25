package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

const dbUrl = "mongodb://localhost:27018"

type Category struct {
	Id          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string
	Description string
}

func main() {
	clientOptions := options.Client().ApplyURI(dbUrl)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Next, let’s ensure that your MongoDB server was found and connected to successfully using the Ping method.
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to mongo")

	// create a database
	collection := client.Database("taskdb").Collection("categories")

	docM := map[string]string{
		"name":        "Open Source",
		"description": "Tasks for open-sourse projects",
	}

	//insert a map object
	_, err = collection.InsertOne(context.Background(), docM)
	if err != nil {
		log.Fatal(err)
	}

	docD := bson.D{
		{"name", "Project"},
		{"description", "Project Tasks"},
	}

	//insert a document slice
	_, err = collection.InsertOne(context.Background(), docD)
	if err != nil {
		log.Fatal(err)
	}

	var count int64
	count, err = collection.CountDocuments(context.Background(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%d inserted documents", count)
}
