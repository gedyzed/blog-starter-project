package repository

import (
	"context"
	"errors"

	domain "github.com/gedyzed/blog-starter-project/Domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

	cursor, err := r.coll.Find(ctx, bson.M{})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	if !cursor.Next(ctx) {
		user.Role = "admin"
	} else {
		user.Role = "regular"
	}

	// Insert the user
	_, err = r.coll.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (r *mongoUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {

	// Check for duplicate username
	filter := bson.M{"email": email}
	result := r.coll.FindOne(ctx, filter)

	if result.Err() != nil && errors.Is(result.Err(), mongo.ErrNoDocuments) {
		return nil, errors.New("user not found")
	}

	if result.Err() != nil {
		return nil, errors.New("error while decoding data")
	}

	var user domain.User
	err := result.Decode(&user)
	if err != nil {
		return nil, errors.New("error while decoding data")
	}

	return &user, nil
}

func (r *mongoUserRepo) GetByUsername(ctx context.Context, username string) (*domain.User, error) {

		// Check for duplicate username
	filter := bson.M{"username": username}
	result := r.coll.FindOne(ctx, filter)

	if result.Err() != nil && errors.Is(result.Err(), mongo.ErrNoDocuments) {
		return nil, errors.New("user not found")
	}

	if result.Err() != nil {
		return nil, errors.New("error while decoding data")
	}

	var user domain.User
	err := result.Decode(&user)
	if err != nil {
		return nil, errors.New("error while decoding data")
	}

	return &user, nil
}

func (r *mongoUserRepo) Update(ctx context.Context, id string, user *domain.User) error{
	return nil
}
func (r *mongoUserRepo) Delete(ctx context.Context, id string)error {
	return nil
}
func (r *mongoUserRepo) Get(ctx context.Context, id string) (*domain.User, error){

	return nil, nil
}

