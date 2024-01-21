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

type respError struct {
	Code    int    `json:"code"`
	Error   string `json:"error"`
	Content string `json:"reason"`
}

func NewUser(c *gin.Context, users *mongo.Collection) {
	body := struct {
		User types.User `json:"user"`
	}{}
	err := json.NewDecoder(c.Request.Body).Decode(&body)
	if err != nil {
		log.Warn(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, respError{
			Code:    http.StatusBadRequest,
			Error:   msgs.ErrUserCreation.Error(),
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
		c.AbortWithStatusJSON(http.StatusInternalServerError, respError{
			Code:    http.StatusInternalServerError,
			Error:   msgs.ErrBadOptions.Error(),
			Content: "Bad options provided in the InsertOne",
		})
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

func GetUsers(c *gin.Context, usersColl *mongo.Collection) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	cursor, err := usersColl.Find(ctx, bson.M{})
	if err != nil {
		log.Error(msgs.ErrBadOptions, "reason", "bad options provided for GetUsers")
		c.AbortWithStatusJSON(http.StatusInternalServerError, respError{
			Code:    http.StatusInternalServerError,
			Error:   msgs.ErrInternal.Error(),
			Content: "Bad options provided for Find",
		})
		return
	}

	var users []types.User
	err = cursor.All(ctx, &users)
	if err != nil {
		log.Error(msgs.ErrDecode, "GetUsers cursor", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, respError{
			Code:    http.StatusInternalServerError,
			Error:   msgs.ErrInternal.Error(),
			Content: "Failed decoding cursor",
		})
		return
	}
	log.Debug(msgs.DebugStruct, "users", fmt.Sprintf("%#v\n", users))
	c.JSON(http.StatusOK, users)
}

func GetUser(c *gin.Context, users *mongo.Collection) {
	id, ok := c.Params.Get("id")
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, respError{
			Code:    http.StatusBadRequest,
			Error:   msgs.ErrInternal.Error(),
			Content: "Failed getting id parameter",
		})
		return
	}

	objid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Error(msgs.ErrObjectIDConv, "message", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, respError{
			Code:    http.StatusBadRequest,
			Error:   msgs.ErrDecode.Error(),
			Content: "wrong id",
		})
		return
	}
	filter := struct {
		ID primitive.ObjectID `bson:"_id"`
	}{
		ID: objid,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	result := users.FindOne(ctx, filter)

	var user types.User
	err = result.Decode(&user)
	if err == mongo.ErrNoDocuments {
		log.Warn(mongo.ErrNoDocuments, "getuser", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, respError{
			Code:    http.StatusBadRequest,
			Error:   msgs.ErrNotFound.Error(),
			Content: "user not found",
		})
		return
	} else if err != nil {
		log.Error(msgs.ErrInternal, "msg", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, respError{
			Code:    http.StatusInternalServerError,
			Error:   msgs.ErrInternal.Error(),
			Content: "failed parsing documents, skill issue",
		})
		return
	}

	log.Debug(msgs.DebugStruct, "user", fmt.Sprintf("%#v\n", user))
	c.JSON(http.StatusOK, user)
}

func UpdateUser(c *gin.Context, users *mongo.Collection) {
	id, ok := c.Params.Get("id")
	if !ok {
		log.Error(msgs.ErrInternal, "UpdateUser", "failed to get params")
		c.AbortWithStatusJSON(http.StatusInternalServerError, respError{
			Code:    http.StatusInternalServerError,
			Error:   msgs.ErrInternal.Error(),
			Content: "failed to get id",
		})
		return
	}

	objid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Error(msgs.ErrDecode, "failed to decode objid", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, respError{
			Code:    http.StatusBadRequest,
			Error:   msgs.ErrDecode.Error(),
			Content: "internal error",
		})
		return
	}

	body := c.Request.Body
	var bdy struct {
		User      types.User `json:"user"`
		Requester types.User `json:"requester"`
	}

	err = json.NewDecoder(body).Decode(&bdy)
	if err != nil {
		log.Warn(msgs.ErrDecode, "body", "wrong format")
		c.AbortWithStatusJSON(http.StatusBadRequest, respError{
			Error:   msgs.ErrInternal.Error(),
			Content: "malformed data",
		})
		return
	}

	log.Debug("ids", objid.String(), bdy.Requester.ID.String())
	if objid != bdy.Requester.ID {
		log.Info(msgs.ErrForbidden, "UpdateUser", "ids aren't equal")
		c.AbortWithStatusJSON(http.StatusForbidden, respError{
			Code:    http.StatusForbidden,
			Error:   msgs.ErrForbidden.Error(),
			Content: "action is forbidden!",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	update := bson.M{"$set": bdy.User}

	updateResult, err := users.UpdateByID(ctx, objid, update)
	if err != nil {
		log.Error(msgs.ErrBadOptions, "UpdateUser", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, respError{
			Code:    http.StatusInternalServerError,
			Error:   msgs.ErrBadOptions.Error(),
			Content: "options failure",
		})
		return
	}

	if updateResult.ModifiedCount == 0 {
		log.Warn(msgs.ErrUpdateFailed, "UpdateUser", updateResult)
		c.AbortWithStatusJSON(http.StatusBadRequest, respError{
			Code: http.StatusBadRequest,
			Error: msgs.ErrUpdateFailed.Error(),
			Content: "failed to find the user",
		})
		return
	}

	c.JSON(http.StatusAccepted, struct {
		Code   int    `json:"code"`
		Status string `json:"status"`
		ID     int64  `json:"id"`
	}{
		Code:   http.StatusCreated,
		Status: "OK",
		ID:     updateResult.ModifiedCount,
	})
}
