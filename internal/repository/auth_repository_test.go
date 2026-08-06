package repository_test

import (
	"context"
	"log"
	"testing"

	"github.com/alexistamher/social-api-go/internal/domain/models"
	"github.com/alexistamher/social-api-go/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthRepository_Register(t *testing.T) {
	db, ctr, err := GetTestDB()
	defer func() {
		err := ctr.Terminate(context.Background())
		require.NoError(t, err)
	}()
	if err != nil {
		t.Fatalf("error opening database: %s", err)
	}

	a := assert.New(t)
	r := repository.NewAuthRepository(db)

	user := &models.User{
		Username:    "JohnC",
		Password:    "passwordHash",
		Email:       "jconnor@mail.com",
		DisplayName: "John Connor",
		// Bio:          "Bio",
		// AvatarURL:    "avatarUrl",
	}

	ID, err := r.Register(user)
	log.Println(err)
	a.NoError(err)
	a.NotEmpty(ID)
}
