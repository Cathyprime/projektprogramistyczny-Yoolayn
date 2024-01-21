package handlers

import (
	"github.com/UniversityOfGdanskProjects/projektprogramistyczny-Yoolayn/internal/types"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewPost(c *gin.Context, coll *mongo.Collection) {
	body := struct {
		Post types.Post `json:"post"`
	}{}

	err := c.BindJSON(&body)
	if err != nil {
		_ = c.AbortWithError(500, err)
		return
	}

	c.IndentedJSON(200, body)
}
