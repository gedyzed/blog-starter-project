package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IUserRepository interface {
	Add(ctx context.Context, user *User) error
	Update(ctx context.Context, id string, user *User) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
}

// User represents a user in the system
type User struct {
	ID        string    `json:"id" bson:"user_id"`
	Firstname string    `json:"firstname" bson:"firstname"`
	LastName  string    `json:"lastname" bson:"lastname"`
	Username  string    `json:"username" bson:"username"`
	Email     string    `json:"email" bson:"email"`
	Role      string    `json:"role" bson:"role"`
	Password  string    `json:"-" bson:"password"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// Profile represents a user's profile information
type Profile struct {
	ID             string    `json:"id" bson:"profile_id"`
	UserID         string    `json:"user_id" bson:"user_id"`
	Bio            string    `json:"bio" bson:"bio"`
	ContactInfo    string    `json:"contact_info" bson:"contact_info"`
	PhoneNumber    string    `json:"phone_number" bson:"phone_number"`
	Location       string    `json:"location" bson:"location"`
	ProfilePicture string    `json:"profile_picture" bson:"profile_picture"`
	CreatedAt      time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" bson:"updated_at"`
}

// Blog represents a blog post
type Blog struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"` // uses MongoDB's native ObjectID
	UserID    primitive.ObjectID `json:"user_id" bson:"user_id"`
	Title     string             `json:"title" bson:"title"`
	Content   string             `json:"content" bson:"content"`
	Created   time.Time          `json:"created" bson:"created"`
	Updated   time.Time          `json:"updated" bson:"updated"`
	ViewCount int                `json:"view_count" bson:"view_count"`
	Tags      []string           `json:"tags" bson:"tags"`
}

// Comment represents a comment on a blog post
type Comment struct {
	ID      string    `json:"id" bson:"comment_id"`
	BlogID  string    `json:"blog_id" bson:"blog_id"`
	UserID  string    `json:"user_id" bson:"user_id"` // Commentor's ID
	Message string    `json:"message" bson:"message"`
	Created time.Time `json:"created" bson:"created"`
}

type ITokenRepository interface {
	Save(context.Context, Token) error
	FindByUserID(context.Context, string) (*Token, error)
	DeleteByUserID(context.Context, string) error
}

// Token represents authentication tokens
type Token struct {
	UserID        string    `json:"user_id" bson:"user_id"`
	AccessToken   string    `json:"access_token" bson:"access_token"`
	RefreshToken  string    `json:"refresh_token" bson:"refresh_token"`
	AccessExpiry  time.Time `json:"access_expiry" bson:"access_expiry"`
	RefreshExpiry time.Time `json:"refresh_expiry" bson:"refresh_expiry"`
	CreatedAt     time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time `json:"expires_at" bson:"expires_at"`
}

// Like represents a like on a blog post
type Like struct {
	ID        string    `json:"id" bson:"like_id"`
	BlogID    string    `json:"blog_id" bson:"blog_id"`
	UserID    string    `json:"user_id" bson:"user_id"`
	Liked     bool      `json:"liked" bson:"liked"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

// Dislike represents a dislike on a blog post
type Dislike struct {
	ID        string    `json:"id" bson:"dislike_id"`
	BlogID    string    `json:"blog_id" bson:"blog_id"`
	UserID    string    `json:"user_id" bson:"user_id"`
	Disliked  bool      `json:"disliked" bson:"disliked"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

// AISuggestion represents AI-generated suggestions
type AISuggestion struct {
	ID         string    `json:"id" bson:"suggestion_id"`
	UserID     string    `json:"user_id" bson:"user_id"`
	BlogID     string    `json:"blog_id" bson:"blog_id"`
	Prompt     string    `json:"prompt" bson:"prompt"`
	Suggestion string    `json:"suggestion" bson:"suggestion"`
	CreatedAt  time.Time `json:"created_at" bson:"created_at"`
}

type IPasswordService interface {
	Hash(string) (string, error)
	Verify(string, string) error
}

type ITokenService interface {
	GenerateTokens(ctx context.Context, userID string) (*Token, error)
	VerifyAccessToken(string) (string, error)
	RefreshTokens(ctx context.Context, refreshToken string) (*Token, error)
}

// BlogUpdateInput for updating a blog
type BlogUpdateInput struct {
	UserID  string   `json:"user_id"`
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

type PaginatedBlogResponse struct {
	Blogs       []Blog `json:"blogs"`
	TotalCount  int    `json:"total_count"`
	TotalPages  int    `json:"total_pages"`
	CurrentPage int    `json:"current_page"`
}
