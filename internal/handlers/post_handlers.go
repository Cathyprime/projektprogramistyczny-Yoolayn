package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/UniversityOfGdanskProjects/projektprogramistyczny-Yoolayn/internal/msgs"
	"github.com/UniversityOfGdanskProjects/projektprogramistyczny-Yoolayn/internal/types"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewPost(c *gin.Context, posts *mongo.Collection) {
	body := struct {
		Post types.Post `json:"post"`
	}{}

	err := decodeBody(c, &body)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	result, optionsErr := posts.InsertOne(ctx, body.Post)
	if optionsErr != nil {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrBadOptions,
			"Bad options provided in the InsertOne",
			optionsErr,
		))
		return
	}

	c.JSON(http.StatusCreated, struct {
		Code   int    `json:"code"`
		Status string `json:"status"`
		ID     string `json:"iD"`
	}{
		Code:   http.StatusCreated,
		Status: "OK",
		ID:     result.InsertedID.(primitive.ObjectID).String(),
	})
}
