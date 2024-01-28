package handlers

import (
	"context"
	"fmt"
	"net/http"
	"redoot/internal/msgs"
	"redoot/internal/types"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateComment(c *gin.Context, comments *mongo.Collection) {
	var body struct {
		Comment   types.Comment     `json:"comment"`
		Requester types.Credentials `json:"requester"`
	}
	err := decodeBody(c, &body)
	if err != nil {
		return
	}

	log.Debug(msgs.DebugStruct, "requester", body.Requester)
	if err := body.Requester.Authorize(); err != nil {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrNotAuthorized,
			"user not authorized",
			"error", err,
		))
		return
	}

	usr, err := body.Requester.ToUser()
	if err != nil {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrInternal,
			"skill issue",
			"error", err,
		))
		return
	}

	if body.Comment.Author != usr.ID {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrForbidden,
			"can't create comment for someone else",
		))
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

func GetComment(c *gin.Context, comments *mongo.Collection) {
	_, _, commentId, err := commentIdParams(c)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	var comment types.Comment
	err = comments.FindOne(ctx, bson.M{"_id": commentId}).Decode(&comment)
	if err == mongo.ErrNoDocuments {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrNotFound,
			"comment not found",
			"GetComment", err,
		))
		return
	} else if err != nil {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrInternal,
			"failed parsing documents, skill issue",
			"msgs", err,
		))
		return
	}

	log.Debug(msgs.DebugStruct, "comment", fmt.Sprintf("%#v\n", comment))
	c.JSON(http.StatusOK, comment)
}

func GetComments(c *gin.Context, commentsColl *mongo.Collection) {
	_, postId, err := postId(c)
	if err != nil {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrInternal,
			"skill issue",
		))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	cursor, err := commentsColl.Find(ctx, bson.M{"post": postId})
	if err != nil {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrInternal,
			"bad options provided for find",
			"reason", "bad options provided for GetComments",
		))
		return
	}

	var comments []types.Comment
	err = cursor.All(ctx, &comments)
	if err != nil {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrInternal,
			"failed decoding cursor",
			"GetComments cursor", err,
		))
		return
	}
	log.Debug(msgs.DebugStruct, "users", fmt.Sprintf("%#v\n", comments))
	c.JSON(http.StatusOK, comments)
}
