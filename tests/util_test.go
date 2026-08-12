package integration_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	cpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupAuthIntegrationDB(t *testing.T) (*gorm.DB, func()) {
	t.Helper()
	ctx := context.Background()

	dbName := "social_media_test"
	dbUser := "postgres"
	dbPassword := "postgres"

	scriptPath, err := filepath.Abs("../social-media-db.sql")
	require.NoError(t, err, "failed to get absolute path of sql script")

	ctr, err := cpostgres.Run(
		ctx,
		"postgres:16-alpine",
		cpostgres.WithInitScripts(scriptPath),
		cpostgres.WithDatabase(dbName),
		cpostgres.WithUsername(dbUser),
		cpostgres.WithPassword(dbPassword),
		cpostgres.BasicWaitStrategies(),
	)
	require.NoError(t, err, "failed to start testcontainer")

	connStr, err := ctr.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err, "failed to get container connection string")

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	require.NoError(t, err, "failed to connect to database")

	cleanup := func() {
		_ = ctr.Terminate(ctx)
	}

	return db, cleanup
}
