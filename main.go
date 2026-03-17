package main

import (
	"net/http"
	"os"

	"github.com/aprimr/blogs-api/db"
	"github.com/aprimr/blogs-api/handlers"
	"github.com/aprimr/blogs-api/middlewares"
	"github.com/aprimr/blogs-api/utils"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		utils.LogFatal("Error loading .env", err)
		return
	}

	// Connect DB
	db.Connect()

	// Create router
	r := chi.NewRouter()

	// Routes
	r.Route("/api/v1", func(r chi.Router) {

		// Public routes
		r.Post("/register", handlers.RegisterUserHandler)
		r.Post("/login", handlers.LoginUserHandler)

		// Private routes (authentication required)
		r.Group(func(r chi.Router) {
			r.Use(middlewares.Authentication)

			r.Post("/blog", handlers.CreateBlogHandler)
		})
	})

	// Start server
	port := ":" + os.Getenv("PORT")
	utils.LogInfo("Server running on PORT" + port)
	err = http.ListenAndServe(port, r)
	if err != nil {
		utils.LogFatal("Error starting Server", err)
	}
}
