package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/aprimr/blogs-api/models"
	"github.com/aprimr/blogs-api/repository"
	"github.com/aprimr/blogs-api/utils"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
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

func GetBlogByBlogidHandler(w http.ResponseWriter, r *http.Request) {
	//
	// Extract blog id from URL params
	// .../blog/:blogid
	//
	blogid := chi.URLParam(r, "blogid")
	if strings.TrimSpace(blogid) == "" {
		utils.SendError(w, "Invalid blogid", http.StatusBadRequest)
		return
	}

	// Call GetBlogByBlogid
	blog, err := repository.GetBlogByBlogid(r.Context(), blogid)
	if err != nil {
		utils.LogError("GetBlogById", err)
		if err == pgx.ErrNoRows {
			utils.SendError(w, "No result found", http.StatusNotFound)
			return
		}
		utils.SendError(w, "Error fetching blog", http.StatusInternalServerError)
		return
	}

	utils.SendSuccess(w, "Blog fetched successfully", blog, http.StatusOK)
}

func DeleteBlogHandler(w http.ResponseWriter, r *http.Request) {
	//
	// Extract blog id from URL params
	// .../blog/:blogid
	//
	blogid := chi.URLParam(r, "blogid")
	if strings.TrimSpace(blogid) == "" {
		utils.SendError(w, "Invalid blogid", http.StatusBadRequest)
		return
	}

	// Get uid from context
	uid, ok := r.Context().Value("uid").(string)
	if !ok || uid == "" {
		utils.SendError(w, "User unauthorized", http.StatusUnauthorized)
		return
	}

	// Call delete
	err := repository.DeleteBlog(r.Context(), uid, blogid)
	if err != nil {
		utils.LogError("DeleteBlog", err)
		if err.Error() == "blog doesnt exists" {
			utils.SendError(w, "Blog not found", http.StatusNotFound)
			return
		}
		utils.SendError(w, "Failed to delete blog", http.StatusInternalServerError)
		return
	}

	utils.SendSuccess(w, "Blog deleted successfully", nil, http.StatusOK)
}

// Expected JSON payload from client:
//
//	{
//		"title": "My first blog",
//		"description": "Short description",
//		"content": "Content body",
//		"is_private": false
//	}
func UpdateBlogHandler(w http.ResponseWriter, r *http.Request) {
	//
	// Extract blog id from URL params
	// .../blog/:blogid
	//
	blogid := chi.URLParam(r, "blogid")
	if strings.TrimSpace(blogid) == "" {
		utils.SendError(w, "Invalid blogid", http.StatusBadRequest)
		return
	}

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
		utils.LogError("Invalid JSON in UpdateBlogHandler", err)
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

	// Call UpdateBlog
	blog, err := repository.UpdateBlog(r.Context(), uid, blogid, blogBody)
	if err != nil {
		utils.LogError("Updating Blog", err)
		if pgx.ErrNoRows == err {
			utils.SendError(w, "Blog not found", http.StatusNotFound)
			return
		}
		utils.SendError(w, "Error updating blog", http.StatusInternalServerError)
		return
	}

	utils.SendSuccess(w, "Blog updated successfully", blog, http.StatusOK)
}
