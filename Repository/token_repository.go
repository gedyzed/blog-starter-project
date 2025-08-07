package repository

import (
	"context"
	"errors"
	"time"

	domain "github.com/gedyzed/blog-starter-project/Domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrTokenNotFound = errors.New("token not found")
	ErrInternalServer = errors.New("internal server error")
)

type mongoTokenRepo struct {
	coll *mongo.Collection
}

func NewMongoTokenRepository(coll *mongo.Collection) domain.ITokenRepo {
	return &mongoTokenRepo{
		coll: coll,
	}
}

func (r *mongoTokenRepo) Save(ctx context.Context, tokens *domain.Token) error {

    oid, err := primitive.ObjectIDFromHex(tokens.UserID)
	if err != nil {
		return domain.ErrIncorrectUserID
	}

	filter := bson.M{"user_id": oid}
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

	_, err = r.coll.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	if err != nil {

	}

	return nil
}

func (r *mongoTokenRepo) FindByUserID(ctx context.Context, userID string) (*domain.Token, error) {
	var tokens domain.Token

	err := r.coll.FindOne(ctx, bson.M{"user_id": userID}).Decode(&tokens)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrTokenNotFound
		}
		return nil, err
	}

	return &tokens, nil
}

func (r *mongoTokenRepo) DeleteByUserID(ctx context.Context, userID string) error {
	result, err := r.coll.DeleteOne(ctx, bson.M{"user_id": userID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return ErrTokenNotFound
	}

	return nil
}

type mongoVTokenRepo struct {
	coll *mongo.Collection
}

func NewMongoVTokenRepository(coll *mongo.Collection) domain.IVTokenRepo {
	return &mongoVTokenRepo{
		coll: coll,
	}
}

func (r *mongoVTokenRepo) CreateVCode(ctx context.Context, token *domain.VToken) error {


	filter := bson.M{"user_id": token.UserID}
	result := r.coll.FindOne(ctx, filter)

	if result.Err() == nil {
		// Token already exists, replace it with the new one
		update := bson.M{
			"$set": bson.M{
				"token_type": token.TokenType,
				"token":      token.Token,
				"expires_at": token.ExpiresAt,
			},
		}

		_, err := r.coll.UpdateOne(ctx, filter, update)
		if err != nil {
			return err
		}
		return nil
	}

	// No existing token found, insert the new one
	_, err := r.coll.InsertOne(ctx, token)
	if err != nil {
		return err
	}

	return nil
}

func (r *mongoVTokenRepo) DeleteVCode(ctx context.Context, id string) error {

	filter := bson.M{"user_id": id}
	result, err := r.coll.DeleteOne(ctx, filter)

	if err != nil {
		return errors.New("internal server error")
	}

	if result.DeletedCount == 0 {
		return errors.New("token not found")
	}

	return nil
}

func (r *mongoVTokenRepo) GetVCode(ctx context.Context, token string) (*domain.VToken, error) {

	filter := bson.M{"token": token}
	result := r.coll.FindOne(ctx, filter)
	if errors.Is(result.Err(), mongo.ErrNoDocuments) {
		return nil, ErrTokenNotFound
	} else if result.Err() != nil {
		return nil, ErrInternalServer
	}

	var existingToken *domain.VToken
	err := result.Decode(&existingToken)
	if err != nil {
		return nil, ErrInternalServer
	}

	return existingToken, nil
}








