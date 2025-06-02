package repositories

import (
	"context"

	"github.com/slashispy/quicktask/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoUserRepository struct {
	Collection *mongo.Collection
}

func (r *MongoUserRepository) Save(user models.User) error {
	_, err := r.Collection.InsertOne(context.TODO(), user)
	return err
}

func (r *MongoUserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.Collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func NewMongoUserRepository(db *mongo.Database, collectionName string) *MongoUserRepository {
	return &MongoUserRepository{
		Collection: db.Collection(collectionName),
	}
}
