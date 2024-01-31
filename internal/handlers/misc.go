package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"redoot/internal/msgs"
	"redoot/internal/types"
	"sync"
	"syscall"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type findResultUsers struct {
	users []types.User
	err   error
}

type findResultBoards struct {
	boards []types.Board
	err    error
}

type findResultPosts struct {
	posts []types.Post
	err   error
}

func decodeBody(c *gin.Context, bdy interface{}) error {
	log.Debug(msgs.DebugStruct, "bdy", bdy)
	err := json.NewDecoder(c.Request.Body).Decode(bdy)
	if err != nil {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrWrongFormat,
			"malformed data",
			"body", fmt.Sprintf("%#v", bdy),
		))
		return err
	}
	return nil
}

func Interrupt(s *http.Server, collections ...*mongo.Collection) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down the Server...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	for _, c := range collections {
		err := c.Drop(ctx)
		if err != nil {
			log.Fatal("Failed to drop collection ", "reason:", err)
		}
	}

	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Error Shutting down: ", "reason:", err)
	}
}

func commentIdParams(c *gin.Context) (boardObjId primitive.ObjectID, postObjId primitive.ObjectID, commentObjId primitive.ObjectID, err error) {
	boardId, ok := c.Params.Get("id")
	if !ok {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrInternal,
			"failed getting board ID",
		))
		return primitive.NilObjectID, primitive.NilObjectID, primitive.NilObjectID, msgs.ErrFailedToGetParams
	}

	boardObjId, err = primitive.ObjectIDFromHex(boardId)
	if err != nil {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrObjectIDConv,
			"wrong id",
			"message", err,
		))
		return primitive.NilObjectID, primitive.NilObjectID, primitive.NilObjectID, msgs.ErrObjectIDConv
	}

	postId, ok := c.Params.Get("postId")
	if !ok {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrInternal,
			"failed getting post id",
		))
		return primitive.NilObjectID, primitive.NilObjectID, primitive.NilObjectID, msgs.ErrFailedToGetParams
	}

	postObjId, err = primitive.ObjectIDFromHex(postId)
	if err != nil {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrObjectIDConv,
			"wrong id",
			"message", err,
		))
		return primitive.NilObjectID, primitive.NilObjectID, primitive.NilObjectID, msgs.ErrObjectIDConv
	}

	commentId, ok := c.Params.Get("commentId")
	if !ok {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrInternal,
			"failed getting comment id",
		))
		return primitive.NilObjectID, primitive.NilObjectID, primitive.NilObjectID, msgs.ErrFailedToGetParams
	}

	commentObjId, err = primitive.ObjectIDFromHex(commentId)
	if err != nil {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrObjectIDConv,
			"wrong id",
			"message", err,
		))
		return primitive.NilObjectID, primitive.NilObjectID, primitive.NilObjectID, msgs.ErrObjectIDConv
	}

	return boardObjId, postObjId, commentObjId, nil
}

func postId(c *gin.Context) (primitive.ObjectID, primitive.ObjectID, error) {
	boardId, ok := c.Params.Get("id")
	if !ok {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrInternal,
			"failed getting board ID",
		))
		return primitive.NilObjectID, primitive.NilObjectID, msgs.ErrFailedToGetParams
	}

	boardObjId, err := primitive.ObjectIDFromHex(boardId)
	if err != nil {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrObjectIDConv,
			"wrong id",
			"message", err,
		))
		return primitive.NilObjectID, primitive.NilObjectID, msgs.ErrObjectIDConv
	}

	postId, ok := c.Params.Get("postId")
	if !ok {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrInternal,
			"failed getting post id",
		))
		return primitive.NilObjectID, primitive.NilObjectID, msgs.ErrFailedToGetParams
	}

	postObjId, err := primitive.ObjectIDFromHex(postId)
	if err != nil {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrObjectIDConv,
			"wrong id",
			"message", err,
		))
		return primitive.NilObjectID, primitive.NilObjectID, msgs.ErrObjectIDConv
	}

	return boardObjId, postObjId, nil
}

func idFromParams(c *gin.Context) (primitive.ObjectID, error) {
	id, ok := c.Params.Get("id")
	if !ok {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrInternal,
			"Failed getting id parameter",
		))
		return primitive.NilObjectID, msgs.ErrFailedToGetParams
	}

	objid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrObjectIDConv,
			"wrong id",
			"message", err, id,
		))
		return primitive.NilObjectID, msgs.ErrObjectIDConv
	}
	return objid, nil
}

func findByFieldUsers(ctx context.Context, coll *mongo.Collection, key, value string, ch chan<- findResultUsers, wg *sync.WaitGroup) {
	log.Debug("Search started for", key, value)
	resp := findResultUsers{}

	pipeline := mongo.Pipeline{
		bson.D{
			{Key: "$match", Value: bson.D{
				{Key: key, Value: primitive.Regex{Pattern: value, Options: "i"}},
			}},
		},
	}

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		log.Debug("Error occurred in Find", key, value)
		resp.users = nil
		resp.err = err
		ch <- resp
		wg.Done()
		return
	}

	var users []types.User
	err = cursor.All(ctx, &users)
	if err != nil {
		log.Debug("Error occured in cursor.All", key, value)
		resp.users = nil
		resp.err = err
		ch <- resp
		wg.Done()
		return
	}
	resp.users = users
	resp.err = err
	ch <- resp
	wg.Done()
	log.Debug("No errors for", key, value)
}

func findByFieldBoards(ctx context.Context, coll *mongo.Collection, key, value string, ch chan<- findResultBoards, wg *sync.WaitGroup) {
	log.Debug("Search started for", key, value)
	resp := findResultBoards{}

	pipeline := mongo.Pipeline{
		bson.D{
			{Key: "$match", Value: bson.D{
				{Key: key, Value: primitive.Regex{Pattern: value, Options: "i"}},
			}},
		},
	}

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		log.Debug("Error occurred in Find", key, value)
		resp.boards = nil
		resp.err = err
		ch <- resp
		wg.Done()
		return
	}

	var boards []types.Board
	err = cursor.All(ctx, &boards)
	if err != nil {
		log.Debug("Error occured in cursor.All", key, value)
		resp.boards = nil
		resp.err = err
		ch <- resp
		wg.Done()
		return
	}
	resp.boards = boards
	resp.err = err
	ch <- resp
	wg.Done()
	log.Debug("No errors for", key, value)
}

func getAndConvert(c *mongo.Collection, id primitive.ObjectID, result any) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	err := c.FindOne(ctx, bson.M{"_id": id}).Decode(result)
	if err == mongo.ErrNoDocuments {
		return err
	} else if err != nil {
		return errors.New("sister, that's skill issue right there")
	}

	return nil
}

func findByFieldPosts(ctx context.Context, coll *mongo.Collection, key, value string, ch chan<- findResultPosts, wg *sync.WaitGroup) {
	log.Debug("Search started for", key, value)
	resp := findResultPosts{}

	pipeline := mongo.Pipeline{
		bson.D{
			{Key: "$match", Value: bson.D{
				{Key: key, Value: primitive.Regex{Pattern: value, Options: "i"}},
			}},
		},
	}

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		log.Debug("Error occurred in Find", key, value)
		resp.posts = nil
		resp.err = err
		ch <- resp
		wg.Done()
		return
	}

	var posts []types.Post
	err = cursor.All(ctx, &posts)
	if err != nil {
		log.Debug("Error occurred in cursor.All", key, value)
		resp.posts = nil
		resp.err = err
		ch <- resp
		wg.Done()
		return
	}
	resp.posts = posts
	resp.err = err
	ch <- resp
	wg.Done()
	log.Debug("No errors for", key, value)
}

func ExportToFile(c *gin.Context, users, boards, posts, comments *mongo.Collection) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	var creds types.Credentials
	err := decodeBody(c, &creds)
	if err != nil {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrWrongFormat,
			"malformed body",
		))
		return
	}

	err = creds.Authorize()
	if err != nil {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrNotAuthorized,
			"user not authorized",
		))
		return
	}

	usr, err := creds.ToUser()
	if err != nil {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrInternal,
			"internal error has occurred",
		))
		return
	}

	if !types.IsAdmin(usr) {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrNotAuthorized,
			"only admins can access this feature",
		))
		return
	}

	result := struct {
		Users    []types.User    `json:"users"`
		Boards   []types.Board   `json:"boards"`
		Posts    []types.Post    `json:"posts"`
		Comments []types.Comment `json:"comments"`
	}{}

	cursor, err := users.Find(ctx, bson.M{})
	if err != nil {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrInternal,
			"failed getting users",
		))
		return
	}

	err = cursor.All(ctx, &result.Users)
	if err != nil {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrInternal,
			"failed decoding users",
		))
		return
	}

	cursor, err = boards.Find(ctx, bson.M{})
	if err != nil {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrInternal,
			"failed getting boards",
		))
		return
	}

	err = cursor.All(ctx, &result.Boards)
	if err != nil {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrInternal,
			"failed decoding boards",
		))
		return
	}

	cursor, err = posts.Find(ctx, bson.M{})
	if err != nil {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrInternal,
			"failed getting posts",
		))
		return
	}

	err = cursor.All(ctx, &result.Posts)
	if err != nil {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrInternal,
			"failed decoding posts",
		))
		return
	}

	cursor, err = comments.Find(ctx, bson.M{})
	if err != nil {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrInternal,
			"failed getting comments",
		))
		return
	}

	err = cursor.All(ctx, &result.Comments)
	if err != nil {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrInternal,
			"failed decoding comments",
		))
		return
	}

	c.JSON(http.StatusOK, result)
}

func ImportFromFile(c *gin.Context, users, boards, posts, comments *mongo.Collection) {
	body := struct {
		Requester types.Credentials `json:"requester"`
		Users     []types.User      `json:"users"`
		Boards    []types.Board     `json:"boards"`
		Posts     []types.Post      `json:"posts"`
		Comments  []types.Comment   `json:"comments"`
	}{}

	err := decodeBody(c, &body)
	if err != nil {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrInternal,
			"couldn't decode the body",
		))
		return
	}

	log.Debug(msgs.DebugStruct, "body", body)

	err = body.Requester.Authorize()
	if err != nil {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrNotAuthorized,
			"user not authorized",
		))
		log.Debug(msgs.DebugStruct, "requester", body.Requester)
		return
	}

	usr, err := body.Requester.ToUser()
	if err != nil {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrInternal,
			"failed to get the user",
		))
		return
	}

	if !types.IsAdmin(usr) {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrNotAuthorized,
			"only admins have access to this feature",
		))
		return
	}

	var wg sync.WaitGroup
	errs := make(chan error, 4)

	wg.Add(4)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
		defer cancel()

		if len(body.Users) == 0 {
			wg.Done()
			return
		}
		inters := make([]interface{}, len(body.Users))
		for i, p := range body.Users {
			log.Debug("name", "action", "adding")
			inters[i] = p
		}
		log.Debug("name", "action", "insert many")
		_, err := users.InsertMany(ctx, inters)
		if err != nil {
			log.Debug("name", "action", "waiting")
			errs <- err
		}
		wg.Done()
		log.Debug("name", "action", "done")
	}()

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
		defer cancel()

		if len(body.Boards) == 0 {
			wg.Done()
			return
		}
		inters := make([]interface{}, len(body.Boards))
		for i, p := range body.Boards {
			log.Debug("name", "action", "adding")
			inters[i] = p
		}
		log.Debug("name", "action", "insert many")
		_, err := boards.InsertMany(ctx, inters)
		if err != nil {
			log.Debug("name", "action", "waiting")
			errs <- err
		}
		wg.Done()
		log.Debug("name", "action", "done")
	}()

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
		defer cancel()

		if len(body.Posts) == 0 {
			wg.Done()
			return
		}
		inters := make([]interface{}, len(body.Posts))
		for i, p := range body.Posts {
			log.Debug("posts", "action", "adding")
			inters[i] = p
		}
		log.Debug("posts", "action", "insert many")
		_, err := posts.InsertMany(ctx, inters)
		if err != nil {
			log.Debug("posts", "action", "waiting")
			errs <- err
		}
		wg.Done()
		log.Debug("posts", "action", "done")
	}()

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
		defer cancel()

		if len(body.Comments) == 0 {
			wg.Done()
			return
		}
		inters := make([]interface{}, len(body.Comments))
		for i, p := range body.Comments {
			log.Debug("comments", "action", "adding")
			inters[i] = p
		}
		log.Debug("comments", "action", "insert many")
		_, err := comments.InsertMany(ctx, inters)
		if err != nil {
			log.Debug("comments", "action", "waiting")
			errs <- err
		}
		wg.Done()
		log.Debug("comments", "action", "done")
	}()
	log.Debug("waiting")
	wg.Wait()
	close(errs)

	if len(errs) > 0 {
		for e := range errs {
			c.AbortWithStatusJSON(msgs.ReportError(
				msgs.ErrInternal,
				e.Error(),
			))
		}
		return
	}

	c.JSON(http.StatusCreated, struct {
		Code    int    `json:"code"`
		Content string `json:"content"`
	}{
		Code:    http.StatusCreated,
		Content: "added",
	})
}
