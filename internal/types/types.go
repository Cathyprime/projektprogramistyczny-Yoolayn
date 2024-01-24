package types

import (
	"context"
	"errors"
	"slices"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	Administrators []User
	Collections    = struct {
		Admins *mongo.Collection
		Client *mongo.Client
	}{}
)

type PrivacyLevel int

var ErrPostEditPermissions = errors.New("User cannot edit the post!")

const (
	Private PrivacyLevel = iota
	Followers
	Public
)

type ContentType int

const (
	Text ContentType = iota
	Image
	Link
)

type User struct {
	ID             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name           string             `json:"name" bson:"name"`
	Bio            string             `json:"bio" bson:"bio"`
	Avatar         string             `json:"avatar" bson:"avatar"`
	Pronouns       string             `json:"pronouns" bson:"pronouns"`
	Password       string             `json:"password" bson:"password"`
	Email          string             `json:"email" bson:"email"`
	ProfilePrivacy PrivacyLevel       `json:"profilePrivacy" bson:"profilePrivacy"`
	PostPrivacy    PrivacyLevel       `json:"postPrivacy" bson:"postPrivacy"`
}

func (u User) Equal(o User) bool {
	if u.Name != o.Name {
		return false
	}
	if u.Bio != o.Bio {
		return false
	}
	if u.Avatar != o.Avatar {
		return false
	}
	if u.Pronouns != o.Pronouns {
		return false
	}
	if u.Password != o.Password {
		return false
	}
	if u.ProfilePrivacy != o.ProfilePrivacy {
		return false
	}
	if u.PostPrivacy != o.PostPrivacy {
		return false
	}
	return true
}

type Board struct {
	ID         primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Name       string               `json:"name" bson:"name"`
	Bio        string               `json:"bio" bson:"bio"`
	Moderators []primitive.ObjectID `json:"moderators" bson:"moderators"`
	Owner      primitive.ObjectID   `json:"owner" bson:"owner"`
	Rules      string               `json:"rules" bson:"rules"`
}

type Post struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title"`
	BodyType    ContentType        `json:"bodyType" bson:"bodyType"`
	BodyContent string             `json:"bodyContent" bson:"bodyContent"`
	Votes       int                `json:"votes" bson:"votes"`
	Author      primitive.ObjectID `json:"author" bson:"author"`
	Board       primitive.ObjectID `json:"board" bson:"board"`
}

func AddAdministrators(newAdmins ...User) {
	Administrators = append(Administrators, newAdmins...)
}

func RemoveAdministrators(remove ...User) {
	for _, v := range remove {
		indexToRemove := -1
		for i, admin := range Administrators {
			if admin.Equal(v) {
				indexToRemove = i
				break
			}
		}
		if indexToRemove != -1 {
			Administrators = append(Administrators[:indexToRemove], Administrators[indexToRemove+1:]...)
		}
	}
}

func IdToStruct(id *primitive.ObjectID, c *mongo.Collection) *mongo.SingleResult {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	return c.FindOne(ctx, bson.M{"_id": id})
}

func IsAdmin(u User) bool {
	return slices.Contains(Administrators, u)
}

func IsModerator(b Board, u User) bool {
	return slices.Contains(b.Moderators, u.ID)
}

func (p Post) CanEditPost(b Board, u User) bool {
	return IsAdmin(u) && IsModerator(b, u) && p.Author == u.ID
}
