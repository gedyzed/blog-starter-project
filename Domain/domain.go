package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user in the system
type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Firstname string             `json:"firstname" bson:"firstname"`
	Lastname  string             `json:"lastname" bson:"lastname"`
	Username  string             `json:"username,omitempty" bson:"username,omitempty"` // optional for OAuth
	Email     string             `json:"email" bson:"email"`
	VCode     string             `json:"vcode,omitempty" bson:"-"` // used only in logic, not saved
	Role      string             `json:"role" bson:"role"`
	Password  string             `json:"password,omitempty" bson:"password,omitempty"` // only for local users
	Provider  string             `json:"provider" bson:"provider"`                     // "local", "google", etc.
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`

	// embedded user profile
	Profile Profile `json:"profile" bson:"profile"` 
}


type ContactInformation struct {
	PhoneNumber string `json:"phone_number"`
	Location    string `json:"location"`
}

type Profile struct {
	Bio                string             `json:"bio" bson:"bio"`
	ContactInfo 	   ContactInformation `json:"contact_info" bson:"contact_information"`
	ProfilePic		   string             `json:"profile_picture" bson:"profile_picture"`
	CreatedAt          time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt          time.Time          `json:"updated_at" bson:"updated_at"`
}

// Blog represents a blog post
type Blog struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"` // uses MongoDB's native ObjectID
	UserID          primitive.ObjectID `json:"user_id" bson:"user_id"`
	AuthorName      string             `json:"author_name" bson:"author_name"`
	Title           string             `json:"title" bson:"title"`
	Content         string             `json:"content" bson:"content"`
	Created         time.Time          `json:"created" bson:"created"`
	Updated         time.Time          `json:"updated" bson:"updated"`
	ViewCount       int                `json:"view_count" bson:"view_count"`
	Tags            []string           `json:"tags" bson:"tags"`
	Likes           int                `json:"likes" bson:"likes"`
	Dislikes        int                `json:"dislikes" bson:"dislikes"`
	LikedUsers      []string           `json:"liked_users" bson:"liked_users"`
	DislikedUsers   []string           `json:"disliked_users" bson:"disliked_users"`
	CommentsCount   int                `json:"comments_count" bson:"comments_count"`
	PopularityScore float64            `json:"popularity_score" bson:"popularity_score"`
}

// Comment represents a comment on a blog post
type Comment struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	BlogID     primitive.ObjectID `json:"blog_id" bson:"blog_id"`
	UserID     primitive.ObjectID `json:"user_id" bson:"user_id"`
	AuthorName string             `json:"author_name" bson:"author_name"`
	Message    string             `json:"message" bson:"message"`
	Created    time.Time          `json:"created" bson:"created"`
	Updated    time.Time          `json:"updated_at" bson:"updated_at"`
}

// Token represents authentication tokens
type Token struct {
	UserID        string 			 `json:"user_id" bson:"user_id"`
	AccessToken   string             `json:"access_token" bson:"access_token"`
	RefreshToken  string             `json:"refresh_token" bson:"refresh_token"`
	AccessExpiry  time.Time          `json:"access_expiry" bson:"access_expiry"`
	RefreshExpiry time.Time          `json:"refresh_expiry" bson:"refresh_expiry"`
	CreatedAt     time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time          `json:"expires_at" bson:"expires_at"`
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
type AIPrompt struct {
	Prompt     string    `json:"prompt"`
}

type VToken struct {
	Email    string    `json:"email" bson:"email"`
	TokenType string    `json:"token_type" bson:"token_type"`
	Token     string    `json:"-" bson:"token"`
	ExpiresAt time.Time `json:"expires_at" bson:"expires_at"`
}

type EmailRequest struct {
	Email string `json:"email"`
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

type PromoteDemoteStruct struct {
	UserID string `json:"user_id" binding:"required"`
}

type ProfileUpdateInput struct {
	UserID 	    string 	    `json:"user_id" binding:"required"`
	Firstname   string       `json:"firstname"`
    Lastname    string       `json:"lastname"`
	Bio         string       `json:"bio"`
	ProfilePic  string       `json:"profile_picture"`
	Location    string       `json:"location"`
	PhoneNumber string      `json:"phone_number"`

}

// struct for google oauth response
type UserInfo struct {
    Sub           string `json:"sub"`
    Name          string `json:"name"`
    GivenName     string `json:"given_name"`
    FamilyName    string `json:"family_name"`
    Picture       string `json:"picture"`
    Email         string `json:"email"`
    EmailVerified bool   `json:"email_verified"`
    Locale        string `json:"locale"`
}





