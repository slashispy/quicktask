package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/slashispy/quicktask/config"
	"github.com/slashispy/quicktask/internal/database"
	"github.com/slashispy/quicktask/internal/handlers"
	"github.com/slashispy/quicktask/internal/middleware"
	"github.com/slashispy/quicktask/internal/repositories"
	"github.com/slashispy/quicktask/internal/usecases"
)

func main() {

	cfg := config.LoadConfig()

	client, db, err := database.ConnectMongoDB(cfg.MongoURI, cfg.MongoDB)
	if err != nil {
		log.Fatalf("❌ Failed to connect to MongoDB: %v", err)
	}

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = client.Disconnect(ctx)
	}()

	userRepo := repositories.NewMongoUserRepository(db, "users")
	userUseCase := &usecases.UserUseCase{
		UserRepo:  userRepo,
		Encrypter: &usecases.BcryptEncrypter{},
	}

	r := chi.NewRouter()

	userHandler := handlers.RegisterUserHandler(userUseCase)
	
	r.Group(func(protected chi.Router) {
		protected.Use(middleware.JWTAuth(cfg.JWTSecret))
		protected.Post("/v1/users", userHandler)
	})

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		log.Printf("🚀 Starting server on port %s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ ListenAndServe: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("❌ Server Shutdown Failed:%+v", err)
	}

}
