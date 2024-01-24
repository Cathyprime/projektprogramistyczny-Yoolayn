package handlers

import (
	"context"
	"encoding/json"
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

type findResult struct {
	users []types.User
	err error
}

func decodeBody(c *gin.Context, bdy interface{}) error {
	log.Debug(msgs.DebugStruct, "bdy", bdy)
	log.Debug(msgs.DebugStruct, "request body", c.Request.Body)
	err := json.NewDecoder(c.Request.Body).Decode(bdy)
	if err != nil {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrWrongFormat,
			"malformed data",
			"body",
			"wrong format",
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
		log.Error(msgs.ErrObjectIDConv)
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrDecode,
			"wrong id",
			"message",
			err,
		))
		return primitive.NilObjectID, msgs.ErrObjectIDConv
	}
	return objid, nil
}

func findByFieldUsers(ctx context.Context, coll *mongo.Collection, key, value string, ch chan<-findResult, wg *sync.WaitGroup) {
	log.Debug("Search started for", key, value)
	resp := findResult{}
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
