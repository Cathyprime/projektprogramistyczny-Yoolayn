package handlers

import (
	"context"
	"net/http"
	"redoot/internal/msgs"
	"redoot/internal/types"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateComment(c *gin.Context, comments *mongo.Collection) {
	var body struct {
		Comment   types.Comment      `json:"comment"`
		Requester primitive.ObjectID `json:"requester"`
	}
	err := decodeBody(c, &body)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	result, err := comments.InsertOne(ctx, body.Comment)
	if err != nil {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrBadOptions,
			"Bad options provided in the InsertOne",
			err,
		))
		return
	}
	c.JSON(http.StatusCreated, struct {
		Code   int    `json:"code"`
		Status string `json:"status"`
		ID     string `json:"id"`
	}{
		Code:   http.StatusCreated,
		Status: "OK",
		ID:     result.InsertedID.(primitive.ObjectID).String(),
	})
}

// func GetComment(c *gin.Context, comments *mongo.Collection) {
// 	boardId, postId, commentId, err := commentIdParams(c)
// 	if err != nil {
// 		return
// 	}
// }
