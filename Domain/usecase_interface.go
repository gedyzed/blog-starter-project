package domain

import "context"

type BlogUsecase interface {
	GetAllBlogs(ctx context.Context, page int, limit int, sort string) (*PaginatedBlogResponse, error)
	GetBlogByID(ctx context.Context, id string) (*Blog, error)
	CreateBlog(ctx context.Context, blog Blog) (*Blog, error)
	UpdateBlog(ctx context.Context, id string, userID string, updatedBlog BlogUpdateInput) error
	DeleteBlog(ctx context.Context, id string, userID string) error
	LikeBlog(ctx context.Context, blogID string, userID string) error
	DislikeBlog(ctx context.Context, blogID string, userID string) error
}
