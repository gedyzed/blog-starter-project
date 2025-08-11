package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	domain "github.com/gedyzed/blog-starter-project/Domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, domain.ErrInvalidUserID
	}

	var tokens domain.Token
	err = r.coll.FindOne(ctx, bson.M{"user_id": oid}).Decode(&tokens)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrTokenNotFound
		}
		return nil, err
	}

	return &tokens, nil
}

func (r *mongoTokenRepo) DeleteByUserID(ctx context.Context, userID string) error {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return domain.ErrInvalidUserID
	}

	result, err := r.coll.DeleteOne(ctx, bson.M{"user_id": objID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return domain.ErrTokenNotFound
	}

	return nil
}

func (r *mongoTokenRepo) FindByAccessToken (ctx context.Context, accessToken string) (string, error) {

	var tokens domain.Token
	err := r.coll.FindOne(ctx, bson.M{"access_token": accessToken}).Decode(&tokens)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", domain.ErrTokenNotFound
		}
		return "", err
	}

	userID := tokens.UserID
	fmt.Println("userID in access Token : ", userID)
	return userID, nil
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

	filter := bson.M{"email": token.Email}
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
		log.Println(err.Error())
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
		return nil, domain.ErrTokenNotFound
	} else if result.Err() != nil {
		return nil, domain.ErrInternalServer
	}

	var existingToken *domain.VToken
	err := result.Decode(&existingToken)
	if err != nil {
		return nil, domain.ErrInternalServer
	}

	return existingToken, nil
}








