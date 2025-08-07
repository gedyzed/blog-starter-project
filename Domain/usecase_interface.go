package domain

import (
	"context"
	"time"
)

type BlogUsecase interface {
	GetAllBlogs(ctx context.Context, page int, limit int, sort string) (*PaginatedBlogResponse, error)
	ViewBlog(ctx context.Context, id string) (*Blog, error)
	CreateBlog(ctx context.Context, blog Blog, userID string) (*Blog, error)
	UpdateBlog(ctx context.Context, id string, userID string, updatedBlog BlogUpdateInput) error
	DeleteBlog(ctx context.Context, id string, userID string, role string) error
	LikeBlog(ctx context.Context, blogID string, userID string) error
	DislikeBlog(ctx context.Context, blogID string, userID string) error
	RefreshPopularity(ctx context.Context, blogID string) error
	FilterBlogs(ctx context.Context, tags []string, startDate, endDate *time.Time, sortBy string, page int, limit int) (*PaginatedBlogResponse, error)
	SearchBlogs(ctx context.Context, keyword string, page, limit int) (*PaginatedBlogResponse, error)
}

type CommentUsecase interface {
	CreateComment(ctx context.Context, blogID string, userID string, message string) (*Comment, error)
	GetAllComments(ctx context.Context, blogID string, page int, limit int, sort string) ([]Comment, int, error)
	GetCommentByID(ctx context.Context, blogID string, commentID string) (*Comment, error)
	EditComment(ctx context.Context, blogID string, commentID string, userID string, message string) error
	DeleteComment(ctx context.Context, blogID string, commentID string, userID string) error
}

type BlogRefreshDispatcher interface {
	Enqueue(blogID string)
}
