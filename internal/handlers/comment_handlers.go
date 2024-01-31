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

func UpdateComment(c *gin.Context, boards, comments *mongo.Collection) {
	boardId, _, commentId, err := commentIdParams(c)
	if err != nil {
		return
	}

	var bdy struct {
		Comment types.Comment `json:"comment"`
		Requester types.Credentials `json:"requester"`
	}
	err = decodeBody(c, &bdy)
	if err != nil {
		return
	}

	if err := bdy.Requester.Authorize(); err != nil {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrNotAuthorized,
			"user not authorized",
			"error", err,
		))
		return
	}

	usr, err := bdy.Requester.ToUser()
	if err != nil {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrInternal,
			"failed getting user",
		))
		return
	}

	var board types.Board
	err = getAndConvert(boards, boardId, &board)
	if err != nil {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrInternal,
			"board making skill issue",
		))
		return
	}

	var comment types.Comment
	err = getAndConvert(comments, commentId, &comment)
	if err != nil {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrInternal,
			"comment making skill issue",
		))
		return
	}

	if !(types.IsAdmin(usr) || types.IsModerator(board, usr) || comment.Author == usr.ID) {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrForbidden,
			"action is forbidden!",
			"UpdateBoard", "is neither an admin, moderator nor owner",
		))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	update := bson.M{"$set": bdy.Comment}

	updateResult, err := comments.UpdateByID(ctx, commentId, update)
	if err != nil {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrBadOptions,
			"options failure",
			"UpdateComment", err,
		))
		return
	}

	if updateResult.ModifiedCount == 0 {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrUpdateFailed,
			"failed to update the comment",
			"UpdateComment", updateResult,
		))
		return
	}
	c.JSON(http.StatusAccepted, struct {
		Code   int    `json:"code"`
		Status string `json:"status"`
	}{
		Code:   http.StatusOK,
		Status: "OK",
	})
}

func DeleteComment(c *gin.Context, boards, comments *mongo.Collection) {
	boardId, _, commentId, err := commentIdParams(c)
	if err != nil {
		return
	}

	var bdy struct {
		Requester types.Credentials `json:"requester"`
	}

	err = decodeBody(c, &bdy)
	if err != nil {
		return
	}

	if err := bdy.Requester.Authorize(); err != nil {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrNotAuthorized,
			"user not authorized",
			"error", err,
		))
	}

	usr, err := bdy.Requester.ToUser()
	if err != nil {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrInternal,
			"failed getting user",
		))
	}

	var board types.Board
	err = getAndConvert(boards, boardId, &board)
	if err != nil {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrInternal,
			"comment making skill issue",
		))
		return
	}

	var comment types.Comment
	err = getAndConvert(comments, commentId, &comment)
	if err != nil {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrInternal,
			"comment making skill issue",
		))
		return
	}

	if !(types.IsAdmin(usr) || types.IsModerator(board, usr) || comment.Author == usr.ID) {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrForbidden,
			"action is forbidden!",
			"UpdateBoard", "is neither an admin, moderator nor owner",
		))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	deleteResult, err := comments.DeleteOne(ctx, bson.M{"_id": commentId})
	if err != nil {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrBadOptions,
			"options failure",
			"DeleteComment", err,
		))
		return
	}

	if deleteResult.DeletedCount != 1 {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrDeleteFailed,
			"failed to update the comment",
			"DeleteComment", deleteResult,
		))
		return
	}

	c.JSON(http.StatusOK, struct {
		Code   int    `json:"code"`
		Status string `json:"status"`
	}{
		Code:   http.StatusOK,
		Status: "OK",
	})
}
