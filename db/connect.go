package db

import (
	"context"
	"os"

	"github.com/aprimr/blogs-api/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func Connect() {
	var err error

	// Get Database URL
	databaseURL := os.Getenv("DATABASE_URL")

	// Inti connection
	Pool, err = pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		utils.LogFatal("Unable to connecting to database: ", err)
	}

	// Ping Database
	err = Pool.Ping(context.Background())
	if err != nil {
		utils.LogFatal("Unable to ping database: ", err)
	}

	utils.LogInfo("Connected to databse")
}
