package main

import (
	"context"
	"net/http"
	"os"
	"redoot/internal/handlers"
	"redoot/internal/msgs"
	"redoot/internal/types"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

const mongoUri = "mongodb://localhost:27017"

var auth = options.Credential{
	Username: "root",
	Password: "example",
}

type connection struct {
	con *mongo.Client
	err error
}

func setupMongo(ch chan<- connection) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUri).SetAuth(auth))
	if err != nil {
		log.Fatal("invalid options: ", err)
	}

	err = client.Ping(ctx, nil)
	ch <- connection{
		con: client,
		err: err,
	}
}

func newStyle() (style *log.Styles) {
	style = log.DefaultStyles()
	pinkText := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ffc0cb"))

	grayText := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#808080"))

	style.Key = pinkText
	style.Value = grayText
	return
}

func newMods() (types.User, types.User, types.User, types.User, types.User) {
	id1, err := primitive.ObjectIDFromHex("65b9521f08488450adcbd92d")
	if err != nil {
		panic(err)
	}
	id2, err := primitive.ObjectIDFromHex("65b9521f08488450adcbd92e")
	if err != nil {
		panic(err)
	}
	id3, err := primitive.ObjectIDFromHex("65b9521f08488450adcbd92f")
	if err != nil {
		panic(err)
	}
	id4, err := primitive.ObjectIDFromHex("65b954c547c4f420dc911a6c")
	if err != nil {
		panic(err)
	}
	id5, err := primitive.ObjectIDFromHex("65b954c547c4f420dc911a6d")
	if err != nil {
		panic(err)
	}

	hash1, err := bcrypt.GenerateFromPassword([]byte("password1"), bcrypt.MinCost)
	if err != nil {
		panic(err)
	}
	hash2, err := bcrypt.GenerateFromPassword([]byte("password2"), bcrypt.MinCost)
	if err != nil {
		panic(err)
	}
	hash3, err := bcrypt.GenerateFromPassword([]byte("password3"), bcrypt.MinCost)
	if err != nil {
		panic(err)
	}
	hash4, err := bcrypt.GenerateFromPassword([]byte("password4"), bcrypt.MinCost)
	if err != nil {
		panic(err)
	}
	hash5, err := bcrypt.GenerateFromPassword([]byte("password5"), bcrypt.MinCost)
	if err != nil {
		panic(err)
	}

	mod1 := types.User{
		ID: id1,
		Name: "Mod1",
		Bio: "Dictator",
		Avatar: "base64encodedfile",
		Pronouns: "over/lord",
		Password: string(hash1),
		Email: "mail@mail.com",
	}
	mod2 := types.User{
		ID: id2,
		Name: "Mod2",
		Bio: "Dictator",
		Avatar: "base64encodedfile",
		Pronouns: "over/lord",
		Password: string(hash2),
		Email: "mail@mail.com",
	}
	mod3 := types.User{
		ID: id3,
		Name: "Mod3",
		Bio: "Dictator",
		Avatar: "base64encodedfile",
		Pronouns: "over/lord",
		Password: string(hash3),
		Email: "mail@mail.com",
	}
	user := types.User{
		ID: id4,
		Name: "regular_user",
		Bio: "Dictator",
		Avatar: "base64encodedfile",
		Pronouns: "over/lord",
		Password: string(hash4),
		Email: "mail@mail.com",
	}
	user2 := types.User{
		ID: id5,
		Name: "regular_user2",
		Bio: "Dictator",
		Avatar: "base64encodedfile",
		Pronouns: "over/lord",
		Password: string(hash5),
		Email: "mail@mail.com",
	}

	return mod1, mod2, mod3, user, user2
}

const (
	LevelsDebug   = "debug"
	LevelsInfo    = "info"
	LevelsWarning = "warn"
	LevelsError   = "error"
	LevelsFatal   = "fatal"
)

func setLevel() {
	switch level := os.Getenv("LOG"); level {
	case LevelsDebug:
		log.SetLevel(log.DebugLevel)
	case LevelsInfo:
		log.SetLevel(log.InfoLevel)
	case LevelsWarning:
		log.SetLevel(log.WarnLevel)
	case LevelsError:
		log.SetLevel(log.ErrorLevel)
	case LevelsFatal:
		log.SetLevel(log.FatalLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}

func main() {
	log.SetStyles(newStyle())
	setLevel()
	log.Info("starting")

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)

	ch := make(chan connection)
	defer close(ch)
	go setupMongo(ch)

	var client *mongo.Client
	connectionResult := <-ch
	if connectionResult.err != nil {
		log.Fatal(msgs.ErrTypeConn, "database connection", connectionResult.err)
	}

	client = connectionResult.con
	defer func() {
		err := client.Disconnect(ctx)
		if err != nil {
			log.Fatal(msgs.ErrTypeConn, "database disconnect", err)
		}
	}()

	db := client.Database("redoot")
	users := db.Collection("users")
	boards := db.Collection("boards")
	posts := db.Collection("posts")
	comments := db.Collection("comments")

	types.Collections.Users = users

	r := gin.Default()

	id, err := primitive.ObjectIDFromHex("65b94ef156e6d7c59f478392")
	if err != nil {
		panic(err)
	}

	admin := types.User{
		ID: id,
		Name: "Administrator",
		Bio: "Dictator",
		Avatar: "base64encodedfile",
		Pronouns: "over/lord",
		Password: "passsword",
		Email: "mail@mail.com",
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.MinCost)
	if err != nil {
		panic(err)
	}
	admin.Password = string(hash)
	types.AddAdministrators(admin)

	mod1, mod2, mod3, user, user2 := newMods()

	toAdd := []interface{}{
		admin,
		mod1,
		mod2,
		mod3,
		user,
		user2,
	}

	_, err = users.InsertMany(ctx, toAdd)
	if err != nil {
		panic(err)
	}

	r.GET("/", func(c *gin.Context) { handlers.MostPopular(c, posts) })

	r.GET("/users", func(c *gin.Context) { handlers.GetUsers(c, users) })
	r.POST("/users", func(c *gin.Context) { handlers.NewUser(c, users) })
	r.GET("/users/:id", func(c *gin.Context) { handlers.GetUser(c, users) })
	r.PUT("/users/:id", func(c *gin.Context) { handlers.UpdateUser(c, users) })
	r.DELETE("/users/:id", func(c *gin.Context) { handlers.DeleteUser(c, users) })
	r.GET("/users/search", func(c *gin.Context) { handlers.SearchUser(c, users) })
	r.GET("/users/popular", func(c *gin.Context) { handlers.MostPopularUsers(c, users, posts, comments) })

	r.GET("/boards", func(c *gin.Context) { handlers.GetBoards(c, boards) })
	r.POST("/boards", func(c *gin.Context) { handlers.NewBoard(c, boards) })
	r.GET("/boards/:id", func(c *gin.Context) { handlers.GetBoard(c, boards) })
	r.PUT("/boards/:id", func(c *gin.Context) { handlers.UpdateBoard(c, boards, users) })
	r.DELETE("/boards/:id", func(c *gin.Context) { handlers.DeleteBoard(c, boards, users) })
	r.GET("/boards/search", func(c *gin.Context) { handlers.SearchBoard(c, boards) })

	r.POST("/boards/:id/posts", func(c *gin.Context) { handlers.NewPost(c, posts) })
	r.GET("/boards/:id/posts/:postId", func(c *gin.Context) { handlers.GetPost(c, posts) })
	r.GET("/boards/:id/posts", func(c *gin.Context) { handlers.GetPosts(c, posts) })
	r.PUT("/boards/:id/posts/:postId", func(c *gin.Context) { handlers.UpdatePost(c, posts, boards, users) })
	r.DELETE("/boards/:id/posts/:postId", func(c *gin.Context) { handlers.DeletePost(c, posts, boards, users) })
	r.GET("/boards/:id/posts/search", func(c *gin.Context) { handlers.SearchPost(c, posts) })

	r.POST("/boards/:id/posts/:postId/comments", func(c *gin.Context) { handlers.CreateComment(c, comments) })
	r.GET("/boards/:id/posts/:postId/comments/:commentId", func(c *gin.Context) { handlers.GetComment(c, comments) })
	r.GET("/boards/:id/posts/:postId/comments", func(c *gin.Context) { handlers.GetComments(c, comments) })
	r.PUT("/boards/:id/posts/:postId/comments/:commentId", func(c *gin.Context) { handlers.UpdateComment(c, boards, comments) })
	r.DELETE("/boards/:id/posts/:postId/comments/:commentId", func(c *gin.Context) { handlers.DeleteComment(c, boards, comments) })

	r.POST("/export", func(c *gin.Context) { handlers.ExportToFile(c, users, boards, posts, comments) })
	r.POST("/import", func(c *gin.Context) { handlers.ImportFromFile(c, users, boards, posts, comments) })

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go handlers.Interrupt(srv, users, boards, posts, comments)

	cancel()
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}
