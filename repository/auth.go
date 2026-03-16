package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aprimr/blogs-api/db"
	"github.com/aprimr/blogs-api/models"
	"github.com/jackc/pgx/v5"
)

func RegisterUser(ctx context.Context, registerBody models.RegisterBody) error {
	// Query for register
	query := "INSERT INTO users (name, email, password, lastLogin) VALUES($1, $2, $3, $4)"

	// Execute register query
	_, err := db.Pool.Exec(ctx, query, registerBody.Name, registerBody.Email, registerBody.Password, time.Now())

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return fmt.Errorf("email already exists")
		}
		return err
	}

	return nil
}

func GetUser(ctx context.Context, loginBody models.LoginBody) (*models.User, error) {
	// Query for getting user details
	query := "SELECT (uid, name, email, password, isVerfied, lastLogin, createdAt) FROM users WHERE email=$1"

	// User model
	user := models.User{}

	// Execute Query
	row := db.Pool.QueryRow(ctx, query, loginBody.Email)
	err := row.Scan(&user.Uid, &user.Name, &user.Email, &user.Password, &user.IsVerified, &user.LastLogin, &user.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("invalid email")
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}
