package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/aprimr/blogs-api/models"
	"github.com/aprimr/blogs-api/repository"
	"github.com/aprimr/blogs-api/utils"
)

// Expected JSON payload from client:
//
//	{
//		"title": "My first blog",
//		"description": "Short description",
//		"content": "Content body",
//		"is_private": false"
//	}
func CreateBlogHandler(w http.ResponseWriter, r *http.Request) {
	// Get uid from context
	uid, ok := r.Context().Value("uid").(string)
	if !ok || uid == "" {
		utils.SendError(w, "User unauthorized", http.StatusUnauthorized)
		return
	}

	// For storing JSON payload
	blogBody := models.BlogBody{}

	// Parse JSON
	err := json.NewDecoder(r.Body).Decode(&blogBody)
	if err != nil {
		utils.LogError("Invalid JSON in CreateBlogHandler", err)
		utils.SendError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate data
	if len(strings.TrimSpace(blogBody.Title)) < 12 {
		utils.SendError(w, "Title must be atleast 12 characters", http.StatusBadRequest)
		return
	}
	if len(strings.TrimSpace(blogBody.Description)) < 30 {
		utils.SendError(w, "Description must be atleast 30 characters", http.StatusBadRequest)
		return
	}
	if len(strings.TrimSpace(blogBody.Content)) < 60 {
		utils.SendError(w, "Content must be atleast 60 characters", http.StatusBadRequest)
		return
	}

	// Call CreateBlog func
	blog, err := repository.CreateBlog(r.Context(), uid, blogBody)
	if err != nil {
		utils.LogError("CreateBlog", err)
		utils.SendError(w, "Failed to create blog", http.StatusInternalServerError)
		return
	}

	utils.SendSuccess(w, "Blog created successfully", blog, http.StatusCreated)
}
