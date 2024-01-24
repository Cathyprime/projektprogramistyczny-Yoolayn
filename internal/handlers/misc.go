package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/UniversityOfGdanskProjects/projektprogramistyczny-Yoolayn/internal/msgs"
	"github.com/UniversityOfGdanskProjects/projektprogramistyczny-Yoolayn/internal/types"
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
			"message", err,
		))
		return primitive.NilObjectID, msgs.ErrObjectIDConv
	}
	return objid, nil
}

func findByFieldUsers(ctx context.Context, coll *mongo.Collection, key, value string, ch chan<- findResultUsers, wg *sync.WaitGroup) {
	log.Debug("Search started for", key, value)
	resp := findResultUsers{}
	filter := bson.D{
		{Key: key, Value: value},
	}

	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		log.Debug("Error occured in Find", key, value)
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
	filter := bson.D{
		{Key: key, Value: value},
	}

	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		log.Debug("Error occured in Find", key, value)
		resp.boards = nil
		resp.err = err
		ch <- resp
		wg.Done()
		return
	}

	var users []types.Board
	err = cursor.All(ctx, &users)
	if err != nil {
		log.Debug("Error occured in cursor.All", key, value)
		resp.boards = nil
		resp.err = err
		ch <- resp
		wg.Done()
		return
	}
	resp.boards = users
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
