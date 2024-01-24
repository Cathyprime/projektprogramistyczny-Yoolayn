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

func UpdatePost(c *gin.Context, posts *mongo.Collection, boards *mongo.Collection, users *mongo.Collection) {
	boardId, postId, err := postId(c)
	if err != nil {
		return
	}

	var bdy struct {
		Post      types.Post         `json:"post"`
		Requester primitive.ObjectID `json:"requester"`
	}
	err = decodeBody(c, &bdy)
	if err != nil {
		return
	}

	var post types.Post
	err = getAndConvert(posts, postId, &post)
	if err != nil {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrInternal,
			"post making skill issue",
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

	var user types.User
	err = getAndConvert(users, bdy.Requester, &user)
	if err != nil {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrInternal,
			"user making skill issue",
		))
		return
	}

	if !(types.IsAdmin(user) || types.IsModerator(board, user) || post.Author == user.ID) {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrForbidden,
			"action is forbidden!",
			"UpdatePost", "is neither an admin, moderator nor owner",
		))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	update := bson.M{"$set": bdy.Post}

	updateResult, err := posts.UpdateByID(ctx, postId, update)
	if err != nil {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrBadOptions,
			"options failure",
			"UpdatePost", err,
		))
		return
	}

	if updateResult.ModifiedCount == 0 {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrUpdateFailed,
			"failed to update the board",
			"UpdateBoard", updateResult,
		))
		return
	}
	c.JSON(http.StatusAccepted, struct {
		Code   int    `json:"code"`
		Status string `json:"status"`
	}{
		Code:   http.StatusCreated,
		Status: "OK",
	})
}
