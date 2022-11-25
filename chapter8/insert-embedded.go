package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

const dbUrl = "mongodb://localhost:27018"

type Task struct {
	Description string
	Due         time.Time
}

type Category struct {
	Id          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string
	Description string
	Tasks       []Task
}

func main() {
	clientOptions := options.Client().ApplyURI(dbUrl)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to mongo")

	collection := client.Database("taskdb").Collection("categories")

	doc := Category{
		primitive.NewObjectID(),
		"Open-Source",
		"Tasks for open-source projects",
		[]Task{
			Task{"Create project in mgo", time.Date(2015, time.August, 10, 0, 0, 0, 0, time.UTC)},
			Task{"Create REST API", time.Date(2015, time.August, 20, 0, 0, 0, 0, time.UTC)},
		},
	}

	//insert a Category object with embedded Tasks
	_, err = collection.InsertOne(context.Background(), &doc)
	if err != nil {
		log.Fatal(err)
	}

	var count int64
	count, err = collection.CountDocuments(context.Background(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d documents inserted", count)

	cursor, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}
	for _, result := range results {
		fmt.Println(result)
	}

}
