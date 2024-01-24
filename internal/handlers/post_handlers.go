package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/UniversityOfGdanskProjects/projektprogramistyczny-Yoolayn/internal/msgs"
	"github.com/UniversityOfGdanskProjects/projektprogramistyczny-Yoolayn/internal/types"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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
		ID     string `json:"id"`
	}{
		Code:   http.StatusCreated,
		Status: "OK",
		ID:     result.InsertedID.(primitive.ObjectID).String(),
	})
}

func GetPost(c *gin.Context, posts *mongo.Collection) {
	boardId, postId, err := postId(c)
	if err != nil {
		return
	}

	filter := bson.M{
		"_id":   postId,
		"board": boardId,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	var post types.Post
	err = posts.FindOne(ctx, filter).Decode(&post)
	if err == mongo.ErrNoDocuments {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrNotFound,
			"post not found!",
		))
		return
	} else if err != nil {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrInternal,
			"an internal error has accured",
			"decode error", err,
		))
		return
	}

	c.JSON(http.StatusOK, post)
}

func GetPosts(c *gin.Context, posts *mongo.Collection) {
	boardId, err := idFromParams(c)
	if err != nil {
		return
	}

	filter := bson.M{
		"board": boardId,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	cursor, err := posts.Find(ctx, filter)
	if err != nil {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrInternal,
			"skill issue",
		))
		return
	}

	var results []types.Post
	err = cursor.All(ctx, &results)
	if err != nil {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrInternal,
			"Failed decoding cursor",
			"GetPosts cursor", err,
		))
		return
	}

	c.JSON(http.StatusOK, results)
}
