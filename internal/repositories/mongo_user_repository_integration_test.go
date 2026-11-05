//go:build integration

// Integration tests for the MongoDB-backed user repository.
//
// These tests spin up a real MongoDB instance in a throwaway Docker container
// using testcontainers-go, so they require a working Docker daemon. They are
// excluded from the default `go test ./...` run via the `integration` build
// tag. Run them with:
//
//	go test -tags=integration ./internal/repositories/...
package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/slashispy/quicktask/internal/models"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// startMongo launches a MongoDB container and returns a connected database and
// a cleanup function.
func startMongo(t *testing.T) (*mongo.Database, func()) {
	t.Helper()
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "mongo:7",
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor:   wait.ForListeningPort("27017/tcp").WithStartupTimeout(60 * time.Second),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("failed to start mongo container: %v", err)
	}

	endpoint, err := container.Endpoint(ctx, "mongodb")
	if err != nil {
		t.Fatalf("failed to get container endpoint: %v", err)
	}

	connectCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(connectCtx, options.Client().ApplyURI(endpoint))
	if err != nil {
		t.Fatalf("failed to connect to mongo: %v", err)
	}
	if err := client.Ping(connectCtx, nil); err != nil {
		t.Fatalf("failed to ping mongo: %v", err)
	}

	cleanup := func() {
		disconnectCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = client.Disconnect(disconnectCtx)
		_ = container.Terminate(context.Background())
	}
	return client.Database("quicktask_test"), cleanup
}

func TestMongoUserRepository_SaveAndFindByEmail(t *testing.T) {
	db, cleanup := startMongo(t)
	defer cleanup()

	repo := NewMongoUserRepository(db, "users")

	user := models.User{Name: "Ada", Email: "ada@example.com", Password: "hashed"}
	if err := repo.Save(user); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	got, err := repo.FindByEmail("ada@example.com")
	if err != nil {
		t.Fatalf("FindByEmail failed: %v", err)
	}
	if got == nil {
		t.Fatal("expected a user, got nil")
	}
	if got.Email != "ada@example.com" || got.Name != "Ada" {
		t.Errorf("unexpected user: %+v", got)
	}
}

func TestMongoUserRepository_FindByEmail_NotFound(t *testing.T) {
	db, cleanup := startMongo(t)
	defer cleanup()

	repo := NewMongoUserRepository(db, "users")

	_, err := repo.FindByEmail("missing@example.com")
	// The repository surfaces mongo.ErrNoDocuments for a missing user; the
	// RegisterUser use case relies on this exact behavior.
	if err != mongo.ErrNoDocuments {
		t.Fatalf("expected mongo.ErrNoDocuments, got %v", err)
	}
}
