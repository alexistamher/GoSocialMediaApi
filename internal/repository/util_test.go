package repository_test

import (
	"context"
	"log"
	"path/filepath"

	cpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetTestDB() (*gorm.DB, *cpostgres.PostgresContainer, error) {
	ctx := context.Background()

	dbName := "postgres"
	dbUser := "postgres"
	dbPassword := "postgres"

	ctr, err := cpostgres.Run(ctx,
		"postgres:16-alpine",
		cpostgres.WithInitScripts(filepath.Join("../../social-media-db.sql")),
		cpostgres.WithDatabase(dbName),
		cpostgres.WithUsername(dbUser),
		cpostgres.WithPassword(dbPassword),
		cpostgres.BasicWaitStrategies(),
	)
	if err != nil {
		log.Printf("failed to start container: %s", err)
		return nil, nil, err
	}

	connStr, err := ctr.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Printf("failed to get connection string: %s", err)
		return nil, nil, err
	}

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Printf("failed to connect to database: %s", err)
		return nil, nil, err
	}

	log.Println("success database connection")

	return db, ctr, nil
}
