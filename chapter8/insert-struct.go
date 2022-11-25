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

	doc := Category{
		primitive.NewObjectID(),
		"Open Source",
		"Tasks for open-source projects",
	}
	//insert a category object
	_, err = collection.InsertOne(context.Background(), &doc)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted doc")

	var count int64
	count, err = collection.CountDocuments(context.Background(), bson.D{})
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%d records inserted", count)
	}
}
