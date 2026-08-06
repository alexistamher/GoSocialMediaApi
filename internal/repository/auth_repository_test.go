package repository_test

import (
	"log"
	"testing"

	"github.com/alexistamher/social-api-go/internal/domain/models"
	"github.com/alexistamher/social-api-go/internal/repository"
	"github.com/stretchr/testify/assert"
)

var db = GetTestDB()

func TestAuthRepository_Register_Success(t *testing.T) {
	tx := db.Begin()
	defer func() {
		tx.Rollback()
	}()

	a := assert.New(t)
	r := repository.NewAuthRepository(tx)

	user := &models.User{
		Username:    "JohnC",
		Password:    "passwordHash",
		Email:       "jconnor@mail.com",
		DisplayName: "John Connor",
	}

	ID, err := r.Register(user)
	a.NoError(err)
	a.NotEmpty(ID)
}

func TestAuthRepository_Login_Success(t *testing.T) {
	tx := db.Begin()
	defer func() {
		tx.Rollback()
	}()

	a := assert.New(t)
	r := repository.NewAuthRepository(tx)

	user := &models.User{
		Username:    "JohnC",
		Password:    "passwordHash",
		Email:       "jconnor@mail.com",
		DisplayName: "John Connor",
	}

	ID, _ := r.Register(user)
	log.Printf("pass: %s", user.Password)

	ruser, err := r.Login(user.Email, user.Password)
	if err != nil {
		t.Fatalf("error logging in: %s", err)
	}

	a.Equal(ID, ruser)
}
func TestAuthRepository_GetUserInfo_Success(t *testing.T) {
	tx := db.Begin()
	defer func() {
		tx.Rollback()
	}()

	a := assert.New(t)
	r := repository.NewAuthRepository(tx)

	user := &models.User{
		Username:    "JohnC",
		Password:    "passwordHash",
		Email:       "jconnor@mail.com",
		DisplayName: "John Connor",
	}

	ID, _ := r.Register(user)
	log.Printf("pass: %s", user.Password)

	ruser, err := r.GetUserInfo(*ID)
	if err != nil {
		t.Fatalf("error logging in: %s", err)
	}

	a.Equal(*ID, ruser.ID)
	a.Equal(user.Username, ruser.Username)
	a.Equal(user.Email, ruser.Email)
	a.Equal(user.DisplayName, ruser.DisplayName)
}
