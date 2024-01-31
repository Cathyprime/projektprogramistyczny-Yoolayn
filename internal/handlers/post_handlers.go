package handlers

import (
	"context"
	"net/http"
	"redoot/internal/msgs"
	"redoot/internal/types"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func MostPopular(c *gin.Context, posts *mongo.Collection) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	lookupAuthorStage := bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "users"},
			{Key: "localField", Value: "author"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "authorInfo"},
		}},
	}

	lookupBoardStage := bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "boards"},
			{Key: "localField", Value: "board"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "boardInfo"},
		}},
	}

	projectFieldsStage := bson.D{
		{Key: "$project", Value: bson.D{
			{Key: "title", Value: 1},
			{Key: "bodyType", Value: 1},
			{Key: "bodyContent", Value: 1},
			{Key: "votes", Value: 1},
			{Key: "author", Value: bson.D{{Key: "$arrayElemAt", Value: bson.A{"$authorInfo.name", 0}}}},
			{Key: "board", Value: bson.D{{Key: "$arrayElemAt", Value: bson.A{"$boardInfo.name", 0}}}},
		}},
	}

	sortStage := bson.D{{Key: "$sort", Value: bson.D{{Key: "votes", Value: -1}}}}

	pipeline := mongo.Pipeline{lookupAuthorStage, lookupBoardStage, projectFieldsStage, sortStage}

	cursor, err := posts.Aggregate(ctx, pipeline)
	if err != nil {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrInternal,
			"connection failed",
		))
		return
	}

	var postssss []types.NicePost
	err = cursor.All(ctx, &postssss)
	if err != nil {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrInternal,
			"connection failed",
		))
		return
	}

	c.JSON(http.StatusOK, postssss)
}

func NewPost(c *gin.Context, posts *mongo.Collection) {
	body := struct {
		Post      types.Post        `json:"post"`
		Requester types.Credentials `json:"requester"`
	}{}

	err := decodeBody(c, &body)
	if err != nil {
		return
	}

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
			"message", err,
		))
		return
	}

	body.Post.Author = usr.ID

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
		Post      types.Post        `json:"post"`
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

	var post types.Post
	err = getAndConvert(posts, postId, &post)
	if err != nil {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrNotFound,
			"post finding skill issue",
		))
		return
	}

	var board types.Board
	err = getAndConvert(boards, boardId, &board)
	if err != nil {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrNotFound,
			"board finding skill issue",
		))
		return
	}

	usr, err := bdy.Requester.ToUser()
	if err != nil {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrNotFound,
			"skill issue finding user",
			"error", err,
		))
		return
	}

	if !(types.IsAdmin(usr) || types.IsModerator(board, usr) || post.Author == usr.ID) {
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
		Code:   http.StatusAccepted,
		Status: "OK",
	})
}

func DeletePost(c *gin.Context, posts *mongo.Collection, boards *mongo.Collection, users *mongo.Collection) {
	boardId, postId, err := postId(c)
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
		return
	}

	var post types.Post
	err = getAndConvert(posts, postId, &post)
	if err != nil {
		log.Fatal("HERE")
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrInternal,
			err.Error(),
		))
		return
	}

	var board types.Board
	err = getAndConvert(boards, boardId, &board)
	if err != nil {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrNotFound,
			"skill issue finding board",
		))
		return
	}

	usr, err := bdy.Requester.ToUser()
	if err != nil {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrNotFound,
			"skill issue finding user",
			"error", err,
		))
		return
	}

	if !(types.IsAdmin(usr) || types.IsModerator(board, usr) || post.Author == usr.ID) {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrForbidden,
			"action is forbidden!",
			"UpdatePost", "is neither an admin, moderator nor owner",
		))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	deleteResult, err := posts.DeleteOne(ctx, bson.M{"_id": postId})
	if err != nil {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrBadOptions,
			"skill issue",
		))
		return
	}

	if deleteResult.DeletedCount != 1 {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrNotFound,
			"post failed to delete",
			"DeleteUser", deleteResult.DeletedCount != 1,
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

func SearchPost(c *gin.Context, posts *mongo.Collection) {
	var length int
	for _, v := range c.Request.URL.Query() {
		length += len(v)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	var wg sync.WaitGroup
	ch := make(chan findResultPosts, length)

	for k, s := range c.Request.URL.Query() {
		for _, v := range s {
			wg.Add(1)
			go findByFieldPosts(ctx, posts, k, v, ch, &wg)
		}
	}

	wg.Wait()
	close(ch)

	var values []types.Post
	for v := range ch {
		if err := v.err; err != nil {
			log.Debug(msgs.DebugSkippedLoop, "struct", v)
			continue
		}
		log.Debug("appending", "values +", v)
		values = append(values, v.posts...)
	}

	if len(values) == 0 {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrNotFound,
			"no posts found with provided parameters",
			"values", values,
		))
		return
	}

	c.JSON(http.StatusOK, values)
}
