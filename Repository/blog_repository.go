package repository

import (
	"context"
	"errors"
	"fmt"
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

func NewBlogRepository(coll *mongo.Collection) domain.BlogRepository {
	return &blogRepository{collection: coll}
}

func (r *blogRepository) GetAllBlogs(ctx context.Context, page int, limit int, sort string) ([]domain.Blog, int, error) {
	var blogs []domain.Blog
	skip := int64((page - 1) * limit)

	findOptions := options.Find().SetSkip(skip).SetLimit(int64(limit))
	switch sort {
	case "popular":
		findOptions.SetSort(bson.D{{Key: "view_count", Value: -1}})
	case "oldest":
		findOptions.SetSort(bson.D{{Key: "created", Value: 1}})
	default:
		findOptions.SetSort(bson.D{{Key: "created", Value: -1}})
	}

	cursor, err := r.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch blogs: %w", err)
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &blogs); err != nil {
		return nil, 0, fmt.Errorf("failed to decode blogs: %w", err)
	}

	totalCount, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count blogs: %w", err)
	}

	return blogs, int(totalCount), nil
}

func (r *blogRepository) GetBlogByID(ctx context.Context, id string) (*domain.Blog, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid blog id: %w", err)
	}

	var blog domain.Blog
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&blog)
	if err != nil {
		return nil, fmt.Errorf("blog not found: %w", err)
	}

	return &blog, nil
}

func (r *blogRepository) IncrementBlogViews(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid blog id: %w", err)
	}
	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$inc": bson.M{"view_count": 1}})
	return err
}

func (r *blogRepository) CreateBlog(ctx context.Context, blog domain.Blog) (*domain.Blog, error) {
	blog.ID = primitive.NewObjectID()
	blog.Created = time.Now()
	blog.Updated = blog.Created
	blog.ViewCount = 0

	_, err := r.collection.InsertOne(ctx, blog)
	if err != nil {
		return nil, fmt.Errorf("failed to insert blog: %w", err)
	}
	return &blog, nil
}

func (r *blogRepository) UpdateBlog(ctx context.Context, id string, userID string, input domain.BlogUpdateInput) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID, "user_id": userObjID}
	update := bson.M{"$set": bson.M{
		"title":   input.Title,
		"content": input.Content,
		"tags":    input.Tags,
		"updated": time.Now(),
	}}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("no blog found")
	}
	return nil
}

func (r *blogRepository) DeleteBlog(ctx context.Context, id string, userID string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	res, err := r.collection.DeleteOne(ctx, bson.M{"_id": objID, "user_id": userObjID})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("no blog found")
	}
	return nil
}
