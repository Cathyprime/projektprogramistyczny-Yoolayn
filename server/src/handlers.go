package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	Name string    `json:"name"`
	Id   uuid.UUID `json:"id"`
}

type Post struct {
	User  User   `json:"user"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

func newPost(c *gin.Context, coll *mongo.Collection) {
	post := struct {
		Post Post `json:"post"`
	}{}
	err := c.BindJSON(&post)
	if err != nil {
		log.Fatal("binding ", err)
	}

	jsonPost, err := json.MarshalIndent(post, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(jsonPost))
	c.IndentedJSON(200, post)
}
