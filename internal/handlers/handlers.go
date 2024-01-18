package handlers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Post struct {
	User  User   `json:"user"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

func NewPost(c *gin.Context, coll *mongo.Collection) {
	body := struct {
		Post Post `json:"post"`
	}{}

	err := c.BindJSON(&body)
	if err != nil {
		_ = c.AbortWithError(500, err)
		return
	}

	c.IndentedJSON(200, body)
}

func GetUsers(c *gin.Context, users *mongo.Collection) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	createSampleUsers(ctx, users)
	cursor, err := users.Find(ctx, bson.D{})
	if err != nil {
		fmt.Println("find err")
		log.Fatal(err)
		_ = c.AbortWithError(500, err)
		return
	}

	results := []bson.M{}
	if err = cursor.All(context.TODO(), &results); err != nil {
		fmt.Println("results err")
		log.Fatal(err)
		_ = c.AbortWithError(500, err)
		return
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
		_ = c.AbortWithError(500, err)
		fmt.Println("cursor err")
		return
	}

	fmt.Println(results)
	c.IndentedJSON(200, results)
	if err := cursor.Close(ctx); err != nil {
		log.Fatal(err)
	}
}

func NewUser(c *gin.Context, users *mongo.Collection) {
	body := struct {
		User User `json:"user"`
	}{}

	err := c.BindJSON(&body)
	if err != nil {
		_ = c.AbortWithError(500, err)
		return
	}
	c.IndentedJSON(200, body)
}
