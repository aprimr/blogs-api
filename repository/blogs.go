package repository

import (
	"context"

	"github.com/aprimr/blogs-api/db"
	"github.com/aprimr/blogs-api/models"
)

func CreateBlog(ctx context.Context, uid string, blogBody models.BlogBody) (*models.Blog, error) {
	// query for creating new row
	query := "INSERT INTO blogs (uid, title, description, content, is_private) VALUES($1, $2, $3, $4, $5) RETURNING blogid, uid, title, description, content, is_deleted, is_private, updated_at, created_at"

	// execute query and scan returned row
	blog := models.Blog{}
	row := db.Pool.QueryRow(ctx, query, uid, blogBody.Title, blogBody.Description, blogBody.Content, blogBody.IsPrivate)
	err := row.Scan(
		&blog.BlogId,
		&blog.Uid,
		&blog.Title,
		&blog.Description,
		&blog.Content,
		&blog.IsDeleted,
		&blog.IsPrivate,
		&blog.UpdatedAt,
		&blog.CreatedAt,
	)

	if err != nil {
		return &models.Blog{}, err
	}

	return &blog, nil
}

func DeleteBlog(ctx context.Context, uid string, blogid string) error {
	query := "UPDATE blogs SET is_deleted=true WHERE uid=$1 AND blogid=$2"

	_, err := db.Pool.Exec(ctx, query, uid, blogid)
	if err != nil {
		return err
	}

	return nil
}
