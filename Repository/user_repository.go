package repository

import (
	"context"
	"errors"
	"time"

	"github.com/gedyzed/blog-starter-project/Domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrUserAlreadyExist = errors.New("user already exists")
)

type mongoUserRepo struct {
	coll *mongo.Collection
}

func NewUserMongoRepo(db *mongo.Database) domain.IUserRepository {
	return &mongoUserRepo{
		coll: db.Collection("users"),
	}
}

func (r *mongoUserRepo) Add(ctx context.Context, user *domain.User) error {
	_, err := r.coll.InsertOne(ctx, user)
	if err != nil {
		if we, ok := err.(mongo.WriteException); ok {
			for _, e := range we.WriteErrors {
				if e.Code == 1100 {
					return ErrUserAlreadyExist
				}

			}
		}

		return err
	}

	return nil
}

func (r *mongoUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User

	filter := bson.M{"email": email}
	err := r.coll.FindOne(ctx, filter).Decode(&user)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (r *mongoUserRepo) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	var user domain.User
	filter := bson.M{"username": username}
	err := r.coll.FindOne(ctx, filter).Decode(&user)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (r *mongoUserRepo) Update(ctx context.Context, id string, user *domain.User) error {
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{{Key: "email", Value: user.Email}, {Key: "role", Value: user.Role}, {Key: "updated_at", Value: time.Now()}, {Key: "password", Value: user.Password}}

	result, err := r.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		if we, ok := err.(mongo.WriteException); ok {
			for _, e := range we.WriteErrors {
				if e.Code == 1100 {
					return ErrUserAlreadyExist
				}

			}
		}

		return err
	}

	if result.MatchedCount == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *mongoUserRepo) Delete(ctx context.Context, id string) error {
	result, err := r.coll.DeleteOne(ctx, bson.D{{Key: "_id", Value: id}})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *mongoUserRepo) Get(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User

	err := r.coll.FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}
