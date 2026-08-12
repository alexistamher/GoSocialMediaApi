package repository_test

import (
	"context"
	"encoding/json"
	"log"
	"path/filepath"

	"github.com/testcontainers/testcontainers-go"
	cpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db, _ = GetTestDB()

func GetTestDB() (*gorm.DB, testcontainers.Container) {
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
		testcontainers.WithCmd("postgres", "-c", "max_connections=300"),
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

	return db, ctr
}

//nolint:unused
func printJson(data any) {
	res, _ := json.MarshalIndent(data, "\t", "\t")
	log.Printf("json: %s", string(res))
}
