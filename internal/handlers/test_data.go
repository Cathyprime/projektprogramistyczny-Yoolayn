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

func CreateSampleUsers(ctx context.Context, c *mongo.Collection) {
	_, err := c.InsertMany(ctx, []interface{}{
		bson.M{
			"name":     "name1",
			"login":    "login1",
			"password": "password1",
		},
		bson.M{
			"name":     "name2",
			"login":    "login2",
			"password": "password2",
		},
		bson.M{
			"name":     "name3",
			"login":    "login3",
			"password": "password3",
		},
		bson.M{
			"name":     "name4",
			"login":    "login4",
			"password": "password4",
		},
		bson.M{
			"name":     "name5",
			"login":    "login5",
			"password": "password5",
		},
	})
	if err != nil {
		log.Fatal("Failed to add to collection", err)
	}
}
