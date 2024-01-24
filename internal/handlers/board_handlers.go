package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/UniversityOfGdanskProjects/projektprogramistyczny-Yoolayn/internal/msgs"
	"github.com/UniversityOfGdanskProjects/projektprogramistyczny-Yoolayn/internal/types"
	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// type Board struct {
// 	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
// 	Name       string             `json:"name" bson:"name"`
// 	Bio        string             `json:"bio" bson:"bio"`
// 	Moderators []User             `json:"moderators" bson:"moderators"`
// 	Owner      User               `json:"owner" bson:"owner"`
// 	Rules      string             `json:"rules" bson:"rules"`
// }

func NewBoard(c *gin.Context, boards *mongo.Collection) {
	body := struct {
		Board types.Board `json:"board"`
	}{}

	err := decodeBody(c, &body)
	if err != nil {
		return
	}

	board := body.Board
	log.Debug(msgs.DebugStruct, "board", fmt.Sprintf("%#v", board))

	if log.GetLevel() == log.DebugLevel {
		debugJSON, _ := json.MarshalIndent(board, "", "\t")
		log.Debug(msgs.DebugJSON, "board", string(debugJSON))
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	result, optionsErr := boards.InsertOne(ctx, board)
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

func GetBoards(c *gin.Context, boardsColl *mongo.Collection) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	cursor, err := boardsColl.Find(ctx, bson.M{})
	if err != nil {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrInternal,
			"Bad options provided for Find",
			"reason", "bad options provided for GetBoards",
		))
		return
	}

	var boards []types.Board
	err = cursor.All(ctx, &boards)
	if err != nil {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrInternal,
			"Failed decoding cursor",
			"GetBoards cursor", err,
		))
		return
	}
	log.Debug(msgs.DebugStruct, "users", fmt.Sprintf("%#v\n", boards))
	c.JSON(http.StatusOK, boards)
}

func GetBoard(c *gin.Context, boards *mongo.Collection) {
	objid, err := idFromParams(c)
	if err != nil {
		return
	}

	filter := bson.M{"_id": objid}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	result := boards.FindOne(ctx, filter)

	var board types.Board
	err = result.Decode(&board)
	if err == mongo.ErrNoDocuments {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrNotFound,
			"board not found",
			"GetBoard", err,
		))
	} else if err != nil {
		c.AbortWithStatusJSON(msgs.ReportError(
			msgs.ErrInternal,
			"failed parsing documents, skill issue",
			"msg", err,
		))
	}

	log.Debug(msgs.DebugStruct, "board", fmt.Sprintf("%#v\n", board))
	c.JSON(http.StatusOK, board)
}
