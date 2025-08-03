package domain

import "context"

type BlogRepository interface {
	GetAllBlogs(ctx context.Context, page int, limit int, sort string) ([]Blog, int, error)
	GetBlogByID(ctx context.Context, id string) (*Blog, error)
	IncrementBlogViews(ctx context.Context, id string) error
	CreateBlog(ctx context.Context, blog Blog) (*Blog, error)
	UpdateBlog(ctx context.Context, id string, userID string, updatedBlog BlogUpdateInput) error
	DeleteBlog(ctx context.Context, id string, userID string) error
}
