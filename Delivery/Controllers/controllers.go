package controllers

import (
	"net/http"

	domain "github.com/gedyzed/blog-starter-project/Domain"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlogHandler struct {
	blogUsecase domain.BlogUsecase
}

func NewBlogHandler(blogUsecase domain.BlogUsecase) *BlogHandler {
	return &BlogHandler{blogUsecase: blogUsecase}

}

func (h *BlogHandler) UpdateBlog(c *gin.Context) {
	blogIDparam := c.Param("id")
	blogID, err := primitive.ObjectIDFromHex(blogIDparam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid blog id"})
		return
	}

	var input domain.BlogUpdateInput
	err = c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(primitive.ObjectID)
	err = h.blogUsecase.UpdateBlog(blogID, userID, input)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "blog updated successfully"})

}

func (h *BlogHandler) DeleteBlog(c *gin.Context) {
	blogIDparam := c.Param("id")
	blogID, err := primitive.ObjectIDFromHex(blogIDparam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid blog id"})
		return
	}

	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID := userIDVal.(primitive.ObjectID)
	err = h.blogUsecase.DeleteBlog(blogID, userID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "blog deleted successfully"})
}
