package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

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
func (bc *BlogHandler) DeleteBlog(c *gin.Context) {
	blogID := c.Param("id")

	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userRole, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	err := bc.blogUsecase.DeleteBlog(c.Request.Context(), blogID, userID.(string), userRole.(string))
	if err != nil {
		if err.Error() == "blog not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
		} else if err.Error() == "unauthorized access" {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to delete this blog"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Blog deleted successfully"})
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

	blog, err := h.blogUsecase.ViewBlog(ctx, blogID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, blog)
}

func (h *BlogHandler) CreateBlog(c *gin.Context) {
	ctx := c.Request.Context()

	// Get user ID from context (set by middleware)
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

	var newBlog domain.Blog
	if err := c.ShouldBindJSON(&newBlog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON body"})
		return
	}
	createdBlog, err := h.blogUsecase.CreateBlog(ctx, newBlog, userID)
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

func (h *BlogHandler) FilterBlogs(c *gin.Context) {
	ctx := c.Request.Context()
	rawTags := c.QueryArray("tags")
	fromDate := c.Query("fromDate")
	toDate := c.Query("toDate")
	sort := strings.TrimSpace(c.Query("sortBy"))

	const format = "2006-01-02"
	var startDate *time.Time
	var endDate *time.Time

	if fromDate != "" {
		t, err := time.Parse(format, fromDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid fromDate format, expected YYYY-MM-DD"})
			return
		}
		startDate = &t
	}
	if toDate != "" {
		t, err := time.Parse(format, toDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid fromDate format, expected YYYY-MM-DD"})
			return
		}
		endDate = &t
	}

	tags := make([]string, 0)
	for _, tag := range rawTags {
		clean := strings.TrimSpace(tag)
		if clean != "" {
			tags = append(tags, clean)
		}
	}

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
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

	result, err := h.blogUsecase.FilterBlogs(ctx, tags, startDate, endDate, sort, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *BlogHandler) SearchBlogs(c *gin.Context) {
	ctx := c.Request.Context()
	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter is required"})
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "page must be a positive integer"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		return
	}

	blogs, err := h.blogUsecase.SearchBlogs(ctx, query, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"query": query,
		"page":  page,
		"blogs": blogs,
	})
}
