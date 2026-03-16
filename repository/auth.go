package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aprimr/blogs-api/db"
	"github.com/aprimr/blogs-api/models"
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
