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

func decodeBody(c *gin.Context, bdy interface{}) error {
	err := json.NewDecoder(c.Request.Body).Decode(bdy)
	if err != nil {
		log.Warn(msgs.ErrDecode, "body", "wrong format")
		c.AbortWithStatusJSON(http.StatusBadRequest, respError{
			Error:   msgs.ErrInternal.Error(),
			Content: "malformed data",
		})
		return err
	}
	return nil
}

func idFromParams(c *gin.Context) (primitive.ObjectID, error) {
	id, ok := c.Params.Get("id")
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, respError{
			Code:    http.StatusInternalServerError,
			Error:   msgs.ErrInternal.Error(),
			Content: "Failed getting id parameter",
		})
		return primitive.NilObjectID, msgs.ErrFailedToGetParams
	}

	objid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Error(msgs.ErrObjectIDConv, "message", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, respError{
			Code:    http.StatusBadRequest,
			Error:   msgs.ErrDecode.Error(),
			Content: "wrong id",
		})
		return primitive.NilObjectID, msgs.ErrObjectIDConv
	}
	return objid, nil
}

func NewUser(c *gin.Context, users *mongo.Collection) {
	body := struct {
		User types.User `json:"user"`
	}{}

	err := decodeBody(c, &body)
	if err != nil {
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
	objid, err := idFromParams(c)
	if err != nil {
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
		c.AbortWithStatusJSON(http.StatusNotFound, respError{
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
	objid, err := idFromParams(c)
	if err != nil {
		return
	}

	var bdy struct {
		User      types.User `json:"user"`
		Requester types.User `json:"requester"`
	}
	err = decodeBody(c, &bdy)
	if err != nil {
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
			Code:    http.StatusBadRequest,
			Error:   msgs.ErrUpdateFailed.Error(),
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

func DeleteUser(c *gin.Context, users *mongo.Collection) {
	objid, err := idFromParams(c)
	if err != nil {
		return
	}

	body := struct {
		Requester types.User `json:"requester"`
	}{}

	err = decodeBody(c, &body)
	if err != nil {
		return
	}

	if body.Requester.ID != objid {
		log.Warn(msgs.ErrUpdateFailed, "UpdateUser", body.Requester.ID == objid)
		c.AbortWithStatusJSON(http.StatusForbidden, respError{
			Code:    http.StatusForbidden,
			Error:   msgs.ErrForbidden.Error(),
			Content: "action forbidden",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	filter := bson.M{
		"_id": objid,
	}

	deleteResult, err := users.DeleteOne(ctx, filter)
	if err != nil {
		log.Warn(msgs.ErrInternal, "DeleteUser", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, respError{
			Code:    http.StatusInternalServerError,
			Error:   msgs.ErrBadOptions.Error(),
			Content: "internal error",
		})
		return
	}

	if deleteResult.DeletedCount != 1 {
		log.Warn(msgs.ErrNotFound, "DeleteUser", deleteResult.DeletedCount != 1)
		c.AbortWithStatusJSON(http.StatusNotFound, respError{
			Code:    http.StatusNotFound,
			Error:   msgs.ErrNotFound.Error(),
			Content: "user failed to delete",
		})
		return
	}
	c.JSON(http.StatusOK, struct {
		Code   int    `json:"code"`
		Status string `json:"status"`
		ID     int64  `json:"id"`
	}{
		Code:   http.StatusCreated,
		Status: "OK",
		ID:     deleteResult.DeletedCount,
	})
}
