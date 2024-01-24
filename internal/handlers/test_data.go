package handlers

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateHelloWorld(ctx context.Context, c *mongo.Collection) *mongo.InsertOneResult {
	result, err := c.InsertOne(ctx, bson.M{
		"message": "hello_world!",
	})
	if err != nil {
		log.Fatal("Failed to add to collection", err)
	}
	return result
}
