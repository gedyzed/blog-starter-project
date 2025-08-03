package repository

import (
	"context"
	"errors"
	"time"

	"github.com/gedyzed/blog-starter-project/Domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrTokenNotFound = errors.New("token not found")
)

type MongoTokenRepository struct {
	collection *mongo.Collection
}

func NewMongoTokenRepository(coll *mongo.Collection) *MongoTokenRepository {
	return &MongoTokenRepository{
		collection: coll,
	}
}

func (r *MongoTokenRepository) Save(ctx context.Context, tokens domain.Token) error {
	filter := bson.M{"user_id": tokens.UserID}

	update := bson.M{
		"$set": bson.M{
			"access_token":   tokens.AccessToken,
			"refresh_token":  tokens.RefreshToken,
			"access_expiry":  tokens.AccessExpiry,
			"refresh_expiry": tokens.RefreshExpiry,
			"updated_at":     time.Now(),
		},
		"$setOnInsert": bson.M{
			"createdAt": time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	if err != nil {

	}

	return nil
}

func (r *MongoTokenRepository) FindByUserID(ctx context.Context, userID string) (*domain.Token, error) {
	var tokens domain.Token

	err := r.collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&tokens)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrTokenNotFound
		}
		return nil, err
	}

	return &tokens, nil
}

func (r *MongoTokenRepository) DeleteByUserID(ctx context.Context, userID string) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"user_id": userID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return ErrTokenNotFound
	}

	return nil
}
