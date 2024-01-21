package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/UniversityOfGdanskProjects/projektprogramistyczny-Yoolayn/internal/msgs"
	"github.com/UniversityOfGdanskProjects/projektprogramistyczny-Yoolayn/internal/types"
	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type respError struct {
	Error   error  `json:"error"`
	Content string `json:"reason"`
}

func NewUser(c *gin.Context, users *mongo.Collection) {
	body := struct {
		User types.User `json:"user"`
	}{}
	err := json.NewDecoder(c.Request.Body).Decode(&body)
	if err != nil {
		log.Error(err)
		c.AbortWithStatusJSON(400, respError{
			Error:   msgs.ErrUserCreation,
			Content: "malformed data",
		})
		return
	}
	usr := body.User
	log.Debug("usr", "struct", fmt.Sprintf("%#v", usr))

	debugJSON, _ := json.Marshal(usr)
	log.Debug("usr", "json", string(debugJSON))

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	result, optionsErr := users.InsertOne(ctx, usr)
	if optionsErr != nil {
		log.Error(optionsErr)
		c.AbortWithStatusJSON(500, respError{
			Error:   msgs.ErrBadOptions,
			Content: "Bad options provided in the InsertOne",
		})
		return
	}
	c.JSON(201, struct {
		Status string `json:"status"`
		ID     string `json:"id"`
	}{
		Status: "OK",
		ID:     result.InsertedID.(primitive.ObjectID).String(),
	})
}

func GetUsers(c *gin.Context, usersColl *mongo.Collection) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	cursor, err := usersColl.Find(ctx, bson.M{})
	if err != nil {
		log.Error(msgs.ErrBadOptions, "reason", "bad options provided for GetUsers")
		c.AbortWithStatusJSON(500, respError{
			Error: msgs.ErrInternal,
			Content: "Bad options provided for Find",
		})
		return
	}

	var users []types.User
	err = cursor.All(ctx, &users)
	if err != nil {
		log.Error(msgs.ErrDecode, "GetUsers cursor", err)
		c.JSON(500, respError{
			Error: msgs.ErrInternal,
			Content: "Failed decoding cursor",
		})
		return
	}
	log.Debug(msgs.DebugStruct, "users", fmt.Sprintf("%#v\n", users))
	c.JSON(200, users)
}
