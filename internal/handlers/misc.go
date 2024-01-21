package handlers

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/log"
	"go.mongodb.org/mongo-driver/mongo"
)

func Interrupt(s *http.Server, c *mongo.Collection) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down the Server...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := c.Drop(ctx)
	if err != nil {
		log.Fatal("Failed to drop users ", "reason:", err)
	}

	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Error Shutting down: ", "reason:", err)
	}
}
