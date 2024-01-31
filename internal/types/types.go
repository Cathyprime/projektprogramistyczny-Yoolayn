package types

import (
	"context"
	"errors"
	"redoot/internal/msgs"
	"slices"
	"time"

	"github.com/charmbracelet/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var (
	Administrators []User
	Collections    = struct {
		Admins []User
		Users  *mongo.Collection
		Client *mongo.Client
	}{}
)

var ErrPostEditPermissions = errors.New("User cannot edit the post!")

type ContentType int

const (
	Text ContentType = iota
	Image
	Link
)

type User struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name"`
	Bio      string             `json:"bio" bson:"bio"`
	Avatar   string             `json:"avatar" bson:"avatar"`
	Pronouns string             `json:"pronouns" bson:"pronouns"`
	Password string             `json:"password" bson:"password"`
	Email    string             `json:"email" bson:"email"`
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
	return true
}

func (u User) IsTaken() bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	err := Collections.Users.FindOne(ctx, bson.M{
		"name": u.Name,
	}).Decode(nil)

	return err == mongo.ErrNoDocuments
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

type Comment struct {
	ID     primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Author primitive.ObjectID `json:"author" bson:"author"`
	Post   primitive.ObjectID `json:"post" bson:"post"`
	Body   string             `json:"body" bson:"body"`
	Votes  int                `json:"votes" bson:"votes"`
}

type NicePost struct {
	Title       string      `json:"title" bson:"title"`
	BodyType    ContentType `json:"bodyType" bson:"bodyType"`
	BodyContent string      `json:"bodyContent" bson:"bodyContent"`
	Author      string      `json:"author" bson:"author"`
	Votes       int         `json:"votes" bson:"votes"`
	Board       string      `json:"board" bson:"board"`
}

type Credentials struct {
	Name       string `json:"name"`
	Password   string `json:"password"`
	authorized bool
}

func (c *Credentials) Authorize() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	if Collections.Users == nil {
		return errors.New("missing users collection in Collections struct; types package")
	}

	var usr User
	err := Collections.Users.FindOne(ctx, bson.M{"name": c.Name}).Decode(&usr)
	if err != nil {
		return err
	}

	log.Debug("password", "usr", usr.Password)
	log.Debug("password", "  c", c.Password)
	err = bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(c.Password))
	if err != nil {
		return err
	}

	c.authorized = true

	return nil
}

func (c Credentials) ToUser() (User, error) {
	if !c.authorized {
		return User{}, msgs.ErrNotAuthorized
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	var usr User
	err := Collections.Users.FindOne(ctx, bson.M{"name": c.Name}).Decode(&usr)
	if err != nil {
		return User{}, err
	}

	return usr, nil
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
