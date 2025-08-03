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
	GetByEmail(ctx context.Context, email string)(*User, error)
	GetByUsername(ctx context.Context, username string)(*User, error)
}


// User represents a user in the system
type User struct {
	ID        string 			 `json:"id" bson:"_id"`
    Firstname string             `json:"firstname" bson:"firstname"`
    Lastname  string             `json:"lastname" bson:"lastname"`
    Username  string             `json:"username" bson:"username"`
    Email     string             `json:"email" bson:"email"`
    VCode     string             `json:"vcode" bson:"-"`
    Role      string             `json:"role" bson:"role"`
    Password  string             `json:"password" bson:"password"`
    CreatedAt time.Time          `json:"created_at" bson:"created_at"`
    UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	
    // ← embed Profile here
    Profile Profile `json:"profile" bson:"profile"`
}

type ContactInformation struct {
	PhoneNumber string `json:"phone_number"`
	Location    string `json:"location"`
}

type Profile struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id"`
	Bio                string             `json:"bio" bson:"bio"`
	ContactInformation ContactInformation  `json:"contact_information" bson:"contact_information"`
	ProfilePicture     string             `json:"profile_picture" bson:"profile_picture"`
	CreatedAt          time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt          time.Time          `json:"updated_at" bson:"updated_at"`
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


type VToken struct {
	UserID string `json:"user_id" bson:"user_id"`
	TokenType string `json:"token_type" bson:"token_type"`
	Token string `json:"-" bson:"token"`
	ExpiresAt time.Time `json:"expires_at" bson:"expires_at"`
}

type EmailRequest struct {
	Email string `json:"email"`
}

type ITokenRepo interface{
	CreateVCode(ctx context.Context, token *VToken) error
	DeleteVCode(ctx context.Context, id string) error
	GetVCode(ctx context.Context, id string)(*VToken, error)
	Save(ctx context.Context, tokens Token) error
	FindByUserID(ctx context.Context, userID string) (*Token, error)
	DeleteByUserID(ctx context.Context, userID string) error
}

type IPasswordService interface {
	Hash(string) (string, error)
	Verify(password, hashedPassword string) error
}

type ITokenService interface {
	SendEmail(to []string, subject string, body string) error
}

type IJWTService interface {
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




