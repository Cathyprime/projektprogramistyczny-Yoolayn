package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoUri = "mongodb://localhost:27017"
	auth     = options.Credential{
		Username: "root",
		Password: "example",
	}
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUri).SetAuth(auth))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := client.Disconnect(ctx)
		if err != nil {
			panic(err)
		}
	}()

	db := client.Database("redoot")
	users := db.Collection("users")
	result, err := users.InsertOne(ctx, bson.M{
		"message": "hello_world!",
	})
	if err != nil {
		log.Fatal("Failed to add to collection", err)
	}

	filter := bson.M{"_id": result.InsertedID}
	var resultFind struct {
		Message string `bson:"message"`
	}
	err = users.FindOne(ctx, filter).Decode(&resultFind)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resultFind)

	err = users.Drop(ctx)
	if err != nil {
		log.Fatal("Failed to drop users", err)
	}
}
