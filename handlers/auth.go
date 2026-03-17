package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/aprimr/blogs-api/models"
	"github.com/aprimr/blogs-api/repository"
	"github.com/aprimr/blogs-api/utils"
	"github.com/aprimr/blogs-api/validation"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUserHandler handles user registration requests.
//
// Expected JSON payload from client:
//
//	{
//		"name": "John Doe",
//		"email": "email@example.com",
//		"password": "YourPassword123!"
//	}
func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	registerBody := models.RegisterBody{}

	// Parse Register Body
	err := json.NewDecoder(r.Body).Decode(&registerBody)
	if err != nil {
		utils.LogError("Invalid JSON in RegisterUserHandler", err)
		utils.SendError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate data
	if validation.IsEmptyString(registerBody.Name) {
		utils.SendError(w, "Name is required", http.StatusBadRequest)
		return
	}
	if validation.IsEmptyString(registerBody.Email) {
		utils.SendError(w, "Email is required", http.StatusBadRequest)
		return
	}
	if validation.IsEmptyString(registerBody.Password) {
		utils.SendError(w, "Password is required", http.StatusBadRequest)
		return
	}
	if !validation.IsValidEmail(registerBody.Email) {
		utils.SendError(w, "Invalid email", http.StatusBadRequest)
		return
	}
	if !validation.IsValidPassword(registerBody.Password) {
		utils.SendError(w, "Password must be at least 8 characters and include a letter, number, and special character", http.StatusBadRequest)
		return
	}

	// Hash Password
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(registerBody.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.LogError("Hashing password in RegisterUserHandler", err)
		utils.SendError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	// Replace actual password by hashPassword in registerBody
	registerBody.Password = string(hashPassword)

	// Create User in db
	err = repository.RegisterUser(r.Context(), registerBody)
	if err != nil {
		if err.Error() == "email already exists" {
			utils.SendError(w, "Email already in use", http.StatusConflict)
			return
		}
		utils.LogError("Registeration in RegisterUserHandler", err)
		utils.SendError(w, "User registration failed", http.StatusInternalServerError)
		return
	}

	utils.SendSuccess(w, "User registration successful", nil, http.StatusCreated)
}

// LoginUserHandler handles user login requests.
//
// Expected JSON payload from client:
//
//	{
//		"email": "email@example.com",
//		"password": "YourPassword123!"
//	}
func LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	loginBody := models.LoginBody{}

	// Parse Login Body
	err := json.NewDecoder(r.Body).Decode(&loginBody)
	if err != nil {
		utils.LogError("Invalid JSON in LoginUserHandler", err)
		utils.SendError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate data
	if validation.IsEmptyString(loginBody.Email) {
		utils.SendError(w, "Email is required", http.StatusBadRequest)
		return
	}
	if validation.IsEmptyString(loginBody.Password) {
		utils.SendError(w, "Password is required", http.StatusBadRequest)
		return
	}
	if !validation.IsValidEmail(loginBody.Email) {
		utils.SendError(w, "Invalid email", http.StatusBadRequest)
		return
	}
	if !validation.IsValidPassword(loginBody.Password) {
		utils.SendError(w, "Password must be at least 8 characters and include a letter, number, and special character", http.StatusBadRequest)
		return
	}

	// Call GetUser
	user, err := repository.GetUser(r.Context(), loginBody)
	if err != nil {
		if err.Error() == "invalid email" {
			utils.LogError("Invalid Email in LoginUserHandler", err)
			utils.SendError(w, "Invalid Email", http.StatusBadRequest)
			return
		}
		utils.LogError("Unexpected in calling GetUser in LoginUserHandler", err)
		utils.SendError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	// Compare user password and hashPassword
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginBody.Password))
	if err != nil {
		utils.SendError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	// Create a jwt token
	jwtToken, err := utils.CreateToken(user.Uid, user.Name, user.Email)
	if err != nil {
		utils.LogError("Creating JWT Token in LoginUserHandler", err)
		utils.SendError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	// Update user lastLogin
	err = repository.UpdateUserLastLogin(r.Context(), user.Uid)
	if err != nil {
		utils.LogError("Updating user last login", err)
		utils.SendError(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	// Send Jwt Token to client
	utils.SendSuccess(w, "User login successful", jwtToken, http.StatusOK)
}
