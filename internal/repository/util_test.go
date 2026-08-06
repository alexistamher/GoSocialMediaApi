package repository_test

import (
	"context"
	"log"
	"path/filepath"

	cpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetTestDB() *gorm.DB {
	ctx := context.Background()

	dbName := "postgres"
	dbUser := "postgres"
	dbPassword := "postgres"

	ctr, err := cpostgres.Run(
		ctx,
		"postgres:16-alpine",
		cpostgres.WithInitScripts(filepath.Join("../../social-media-db.sql")),
		cpostgres.WithDatabase(dbName),
		cpostgres.WithUsername(dbUser),
		cpostgres.WithPassword(dbPassword),
		cpostgres.BasicWaitStrategies(),
	)
	if err != nil {
		panic("failed to start container: " + err.Error())
	}

	connStr, err := ctr.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		panic("failed to get connection string: " + err.Error())
	}

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database: " + err.Error())
	}

	log.Println("success database connection")

	return db
}
