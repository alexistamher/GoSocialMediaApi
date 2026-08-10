package repository_test

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/alexistamher/social-api-go/internal/domain/models"
	"github.com/alexistamher/social-api-go/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestCommentRepository_AddCommentToPost_Success(t *testing.T) {
	tx := db.Begin()
	t.Cleanup(func() {
		tx.Rollback()
	})
	a := assert.New(t)
	postRepo := repository.NewPostRepository(tx)
	cmntRepo := repository.NewCommentRepository(tx)
	authRepo := repository.NewAuthRepository(tx)
	var userIds []*string

	for i := range 3 {
		user := &models.User{
			Username:    "JohnC" + strconv.Itoa(i),
			Password:    "passwordHash",
			Email:       "jconnor" + strconv.Itoa(i) + "@mail.com",
			DisplayName: "John Connor" + strconv.Itoa(i),
		}
		UserID, _ := authRepo.Register(user)
		userIds = append(userIds, UserID)
	}

	post := &models.Post{
		Content:    "This is a test post",
		Visibility: models.Public,
		Author:     models.User{ID: *userIds[0]},
	}
	rpost, _ := postRepo.AddPost(post)

	var commentID *string
	for i := range 5 {
		userIdx := rand.Intn(2)
		comment := &models.Comment{
			Content:         "this is a comment" + strconv.Itoa(i),
			Author:          models.User{ID: *userIds[userIdx]},
			PostID:          rpost.ID,
			ParentCommentID: commentID,
		}
		rcomment, _ := cmntRepo.AddComment(comment)
		if rand.Intn(2) == 0 {
			commentID = &rcomment.ID
		} else {
			commentID = nil
		}

	}

	comments, _ := cmntRepo.GetCommentsByPostID(rpost.ID)
	a.Len(comments, 5)
}

func TestCommentRepository_RemoveCommentSuccess(t *testing.T) {
	tx := db.Begin()
	t.Cleanup(func() {
		tx.Rollback()
	})
	a := assert.New(t)
	postRepo := repository.NewPostRepository(tx)
	cmntRepo := repository.NewCommentRepository(tx)
	authRepo := repository.NewAuthRepository(tx)
	user := &models.User{
		Username:    "JohnC",
		Password:    "passwordHash",
		Email:       "jconnor@mail.com",
		DisplayName: "John Connor",
	}

	userId, _ := authRepo.Register(user)
	post := &models.Post{
		Content:    "This is a test post",
		Visibility: models.Public,
		Author:     models.User{ID: *userId},
	}
	rpost, _ := postRepo.AddPost(post)

	comment := &models.Comment{
		Content:         "this is a comment",
		Author:          models.User{ID: *userId},
		PostID:          rpost.ID,
		ParentCommentID: nil,
	}
	rcomment, _ := cmntRepo.AddComment(comment)
	_ = rcomment
	comments, _ := cmntRepo.GetCommentsByPostID(rpost.ID)
	a.Len(comments, 1)
	err := cmntRepo.DeleteComment(rcomment.ID)
	a.NoError(err)
	comments, _ = cmntRepo.GetCommentsByPostID(rpost.ID)
	a.Len(comments, 0)
}

func TestCommentRepository_RemoveNestedCommentSuccess(t *testing.T) {
	tx := db.Begin()
	t.Cleanup(func() {
		tx.Rollback()
	})
	a := assert.New(t)
	postRepo := repository.NewPostRepository(tx)
	cmntRepo := repository.NewCommentRepository(tx)
	authRepo := repository.NewAuthRepository(tx)
	user := &models.User{
		Username:    "JohnC",
		Password:    "passwordHash",
		Email:       "jconnor@mail.com",
		DisplayName: "John Connor",
	}

	userId, _ := authRepo.Register(user)
	post := &models.Post{
		Content:    "This is a test post",
		Visibility: models.Public,
		Author:     models.User{ID: *userId},
	}
	rpost, _ := postRepo.AddPost(post)

	comment := &models.Comment{
		Content:         "this is a comment",
		Author:          models.User{ID: *userId},
		PostID:          rpost.ID,
		ParentCommentID: nil,
	}
	rcomment, err := cmntRepo.AddComment(comment)
	a.NoError(err)

	comment.ParentCommentID = &rcomment.ID
	_, _ = cmntRepo.AddComment(comment)

	comments, _ := cmntRepo.GetCommentsByPostID(rpost.ID)
	a.Len(comments, 2)
	err = cmntRepo.DeleteComment(rcomment.ID)
	a.NoError(err)
	comments, _ = cmntRepo.GetCommentsByPostID(rpost.ID)
	a.Len(comments, 0)

}

func TestCommentRepository_RemoveCommentWithAllReactions_Success(t *testing.T) {
	tx := db.Begin()
	t.Cleanup(func() {
		tx.Rollback()
	})
	a := assert.New(t)
	postRepo := repository.NewPostRepository(tx)
	cmntRepo := repository.NewCommentRepository(tx)
	authRepo := repository.NewAuthRepository(tx)

	userIds := make([]*string, 3)
	for i := range 3 {
		user := &models.User{
			Username:    "JohnC" + strconv.Itoa(i),
			Password:    "passwordHash",
			Email:       "jconnor" + strconv.Itoa(i) + "@mail.com",
			DisplayName: "John Connor" + strconv.Itoa(i),
		}
		userId, _ := authRepo.Register(user)
		userIds[i] = userId
	}

	post := &models.Post{
		Content:    "This is a test post",
		Visibility: models.Public,
		Author:     models.User{ID: *userIds[0]},
	}
	rpost, _ := postRepo.AddPost(post)

	comment := &models.Comment{
		Content:         "this is a comment",
		Author:          models.User{ID: *userIds[0]},
		PostID:          rpost.ID,
		ParentCommentID: nil,
	}
	rcomment, err := cmntRepo.AddComment(comment)
	a.NoError(err)

	postRepo.AddReaction(rcomment.ID, *userIds[0], "like", "comment")
	postRepo.AddReaction(rcomment.ID, *userIds[1], "love", "comment")
	postRepo.AddReaction(rcomment.ID, *userIds[2], "haha", "comment")

	reactions, err := postRepo.GetTargetReactions(rcomment.ID)
	a.NoError(err)
	a.Len(reactions, 3)

	err = cmntRepo.DeleteComment(rcomment.ID)
	a.NoError(err)

	reactions, err = postRepo.GetTargetReactions(rcomment.ID)
	a.NoError(err)
	a.Len(reactions, 0)

	comments, _ := cmntRepo.GetCommentsByPostID(rpost.ID)
	a.Len(comments, 0)
}
