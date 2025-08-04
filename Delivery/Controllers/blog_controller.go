package controllers

import (
	"net/http"
	"strconv"

	domain "github.com/gedyzed/blog-starter-project/Domain"
	"github.com/gin-gonic/gin"
)

type BlogHandler struct {
	blogUsecase domain.BlogUsecase
}

func NewBlogHandler(blogUsecase domain.BlogUsecase) *BlogHandler {
	return &BlogHandler{blogUsecase: blogUsecase}
}

func (h *BlogHandler) UpdateBlog(c *gin.Context) {
	id := c.Param("id")

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user ID is not a string"})
		return
	}

	var input domain.BlogUpdateInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.blogUsecase.UpdateBlog(c.Request.Context(), id, userIDStr, input) // CHANGED: added context
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "blog updated successfully"})
}

func (h *BlogHandler) DeleteBlog(c *gin.Context) {
	id := c.Param("id")

	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, ok := userIDVal.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user ID is not a string"})
		return
	}

	roleVal, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "role not found"})
		return

	}

	role, ok := roleVal.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "role is not a string"})
		return
	}

	blog, err := h.blogUsecase.GetBlogByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "blog not found"})
		return
	}
	if blog.UserID.Hex() != userID && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "not authorized to delete this blog"})
		return
	}
	err = h.blogUsecase.DeleteBlog(c.Request.Context(), id, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "blog deleted successfully"})
}

func (h *BlogHandler) GetAllBlogs(c *gin.Context) {
	ctx := c.Request.Context()

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	sort := c.DefaultQuery("sort", "latest")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page number"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		return
	}

	result, err := h.blogUsecase.GetAllBlogs(ctx, page, limit, sort)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *BlogHandler) GetBlogById(c *gin.Context) {
	ctx := c.Request.Context()
	blogID := c.Param("id")

	blog, err := h.blogUsecase.GetBlogByID(ctx, blogID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, blog)
}

func (h *BlogHandler) CreateBlog(c *gin.Context) {
	ctx := c.Request.Context()
	var newBlog domain.Blog

	if err := c.ShouldBindJSON(&newBlog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON body"})
		return
	}

	createdBlog, err := h.blogUsecase.CreateBlog(ctx, newBlog)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdBlog)
}

func (h *BlogHandler) LikeBlog(c *gin.Context) {
	ctx := c.Request.Context()
	blogID := c.Param("id")

	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, ok := userIDVal.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user ID isnot string"})
		return
	}
	err := h.blogUsecase.LikeBlog(ctx, blogID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "blog liked successfully"})
}

func (h *BlogHandler) DislikeBlog(c *gin.Context) {
	ctx := c.Request.Context()
	blogID := c.Param("id")

	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, ok := userIDVal.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user ID isnot a string"})
		return
	}

	err := h.blogUsecase.DislikeBlog(ctx, blogID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "blog disliked successfully"})

}
