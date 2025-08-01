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

type blogRepository struct {
	collection *mongo.Collection
}

func NewBlogRepository(coll *mongo.Collection) domain.IBlogRepository {
	return &blogRepository{collection: coll}

}

func (r *blogRepository) DeleteBlog(id string, userID string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	userpreID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID, "user_id": userpreID}
	ctx := context.Background()
	res, err := r.collection.DeleteOne(ctx, filter)

	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("no blog found")
	}

	return nil
}

func (r *blogRepository) UpdateBlog(id string, userID string, updatedBlog domain.BlogUpdateInput) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	userpreID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	ctx := context.Background()
	filter := bson.M{"_id": objID, "user_id": userpreID}
	update := bson.M{
		"title":   updatedBlog.Title,
		"content": updatedBlog.Content,
		"tags":    updatedBlog.Tags,
		"updated": time.Now(),
	}

	opts := options.Update().SetUpsert(false)
	res, err := r.collection.UpdateOne(ctx, filter, bson.M{"$set": update}, opts)

	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("no blog found")
	}
	return nil
}
