package repository_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/alexistamher/social-api-go/internal/domain/models"
	"github.com/alexistamher/social-api-go/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestPostRespository_AddPost_Success(t *testing.T) {
	tx := db.Begin()
	t.Cleanup(func() {
		tx.Rollback()
	})
	a := assert.New(t)
	authRepo := repository.NewAuthRepository(tx)
	postRepo := repository.NewPostRepository(tx)

	user := &models.User{
		Username:    "JohnC",
		Password:    "passwordHash",
		Email:       "jconnor@mail.com",
		DisplayName: "John Connor",
	}

	UserID, _ := authRepo.Register(user)

	post := &models.Post{
		Content:    "this is a post test",
		Visibility: "public",
		Author:     models.User{ID: *UserID},
	}

	rpost, _ := postRepo.AddPost(post)
	post = &models.Post{
		Content:    "this is a post test with parent",
		Visibility: models.Friends,
		Author:     models.User{ID: *UserID},
		PostParent: rpost,
	}
	_, _ = postRepo.AddPost(post)

	posts, _, _ := postRepo.GetPostsByUserID(*UserID, nil, uint(2))

	a.Len(posts, 2)
	a.Equal(posts[1].Content, post.Content)
	a.Equal(models.Public, posts[0].Visibility)
	a.Equal(models.Friends, posts[1].Visibility)
	a.Equal(posts[0].Author.ID, post.Author.ID)
	a.Equal(posts[1].Author.ID, post.Author.ID)
}

func TestPostRepository_GetFeedNotes_Success(t *testing.T) {
	tx := db.Begin()
	t.Cleanup(func() {
		tx.Rollback()
	})
	a := assert.New(t)
	postRepo := repository.NewPostRepository(tx)
	authRepo := repository.NewAuthRepository(tx)
	user := &models.User{
		Username:    "JohnC",
		Password:    "passwordHash",
		Email:       "jconnor@mail.com",
		DisplayName: "John Connor",
	}

	UserID, _ := authRepo.Register(user)

	for i := range 100 {
		time.Sleep(time.Duration(1 * int(time.Millisecond)))
		post := &models.Post{
			Content:    "post #" + strconv.Itoa(i),
			Visibility: models.Public,
			Author:     models.User{ID: *UserID},
		}
		_, _ = postRepo.AddPost(post)
	}

	limit := uint(40)
	var offset *int
	posts, offset, _ := postRepo.GetPostsByUserID(*UserID, nil, limit)
	postIds := make([]string, len(posts))
	for i, post := range posts {
		postIds[i] = post.Content
	}
	a.Len(posts, 40)
	posts, offset, _ = postRepo.GetPostsByUserID(*UserID, offset, limit)
	for i, post := range posts {
		postIds[i] = post.Content
	}
	a.Len(posts, 40)
	limit = 20
	posts, _, _ = postRepo.GetPostsByUserID(*UserID, offset, limit)
	postIds = make([]string, len(posts))
	for i, post := range posts {
		postIds[i] = post.Content
	}
	a.Len(posts, 20)
}
