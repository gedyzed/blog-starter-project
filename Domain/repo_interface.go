package domain

import "context"

type BlogRepository interface {
	GetAllBlogs(ctx context.Context, page int, limit int, sort string) ([]Blog, int, error)
	GetBlogByID(ctx context.Context, id string) (*Blog, error)
	IncrementBlogViews(ctx context.Context, id string) error
	CreateBlog(ctx context.Context, blog Blog) (*Blog, error)
	UpdateBlog(ctx context.Context, id string, userID string, updatedBlog BlogUpdateInput) error
	DeleteBlog(ctx context.Context, id string) error
	LikeBlog(ctx context.Context, blogID string, userID string) error
	DislikeBlog(ctx context.Context, blogID string, userID string) error
}

type CommentRepository interface {
	CreateComment(ctx context.Context, blogID string, userID string, comment Comment) (*Comment, error)
	GetAllComments(ctx context.Context, blogID string, page int, limit int, sort string) ([]Comment, int, error)
	GetCommentByID(ctx context.Context, blogID string, id string) (*Comment, error)
	EditComment(ctx context.Context, blogID string, id string, userID string, message string) error
	DeleteComment(ctx context.Context, blogID string, id string, userID string) error
}


type IUserRepository interface {

	Add(ctx context.Context, user *User) (string, error)
	Update(ctx context.Context, filterField, filterValue string, user *User) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string)(*User, error)
	GetByUsername(ctx context.Context, username string)(*User, error)
}

type ITokenRepo interface{
	Save(ctx context.Context, tokens *Token) error
	FindByUserID(ctx context.Context, userID string) (*Token, error)
	DeleteByUserID(ctx context.Context, userID string) error
}

type IVTokenRepo interface {
	CreateVCode(ctx context.Context, token *VToken) error
	DeleteVCode(ctx context.Context, id string) error
	GetVCode(ctx context.Context, id string)(*VToken, error)

}

type IPasswordService interface {
	Hash(string) (string, error)
	Verify(password, hashedPassword string) error
}

type IVTokenService interface {
	SendEmail(to []string, subject string, body string) error
}

type ITokenService interface {
	GenerateTokens(ctx context.Context, userID string) (*Token, error)
	VerifyAccessToken(string) (string, error)
	RefreshTokens(ctx context.Context, refreshToken string) (*Token, error)
}

