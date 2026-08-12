package repository_test

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/alexistamher/social-api-go/internal/domain/models"
	"github.com/alexistamher/social-api-go/internal/repository"
	rmodels "github.com/alexistamher/social-api-go/internal/repository/models"
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
		Author:     models.Author{ID: *userIds[0]},
	}
	rpost, _ := postRepo.AddPost(post)

	var commentID *string
	for i := range 5 {
		userIdx := rand.Intn(2)
		comment := &models.Comment{
			Content:         "this is a comment" + strconv.Itoa(i),
			Author:          models.Author{ID: *userIds[userIdx]},
			PostID:          rpost.ID,
			ParentCommentID: commentID,
		}
		rcomment, _ := cmntRepo.AddComment(comment)
		commentID = &rcomment.ID
	}

	var commentsCount int64
	_ = tx.Model(&rmodels.Comments{}).Where("post_id = ?", rpost.ID).Count(&commentsCount)
	a.Equal(commentsCount, int64(5))

	comments, _ := postRepo.GetCommentsByPostID(rpost.ID)
	a.Len(comments, 1)
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
		Author:     models.Author{ID: *userId},
	}
	rpost, _ := postRepo.AddPost(post)

	comment := &models.Comment{
		Content:         "this is a comment",
		Author:          models.Author{ID: *userId},
		PostID:          rpost.ID,
		ParentCommentID: nil,
	}
	rcomment, _ := cmntRepo.AddComment(comment)

	comment.ID = rcomment.ID
	comment.ParentCommentID = &comment.ID
	rcomment, _ = cmntRepo.AddComment(comment)
	commentDetails, _ := cmntRepo.GetCommentByID(*rcomment.ParentCommentID)
	a.Len(commentDetails.Comments, 1)
	err := cmntRepo.DeleteComment(rcomment.ID)
	a.NoError(err)
	commentDetails, _ = cmntRepo.GetCommentByID(rcomment.ID)
	a.Len(commentDetails.Comments, 0)
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
		Author:     models.Author{ID: *userId},
	}
	rpost, _ := postRepo.AddPost(post)

	comment := &models.Comment{
		Content:         "this is a comment",
		Author:          models.Author{ID: *userId},
		PostID:          rpost.ID,
		ParentCommentID: nil,
	}
	rcomment, err := cmntRepo.AddComment(comment)
	a.NoError(err)

	comment.ParentCommentID = &rcomment.ID
	rcomment, _ = cmntRepo.AddComment(comment)
	var getAllPostCommentsCount = func(postId string) (int64, error) {
		var count int64
		if err := tx.Table("comments").
			Where("post_id = ? AND deleted_at IS NULL", postId).
			Count(&count).Error; err != nil {
			return 0, err
		}
		return count, nil
	}

	comments, _ := getAllPostCommentsCount(rpost.ID)
	a.Equal(int64(2), comments)
	err = cmntRepo.DeleteComment(*rcomment.ParentCommentID)
	a.NoError(err)
	comments, _ = getAllPostCommentsCount(rpost.ID)
	a.Equal(int64(0), comments)

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
	reacRepo := repository.NewReactionRepository(tx)

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
		Author:     models.Author{ID: *userIds[0]},
	}
	rpost, _ := postRepo.AddPost(post)

	comment := &models.Comment{
		Content:         "this is a comment",
		Author:          models.Author{ID: *userIds[0]},
		PostID:          rpost.ID,
		ParentCommentID: nil,
	}
	rcomment, err := cmntRepo.AddComment(comment)
	a.NoError(err)

	reacRepo.AddReaction(rcomment.ID, *userIds[0], "like", "comment")
	reacRepo.AddReaction(rcomment.ID, *userIds[1], "love", "comment")
	reacRepo.AddReaction(rcomment.ID, *userIds[2], "haha", "comment")

	reactions, err := reacRepo.GetTargetReactions(rcomment.ID)
	a.NoError(err)
	a.Len(reactions, 3)

	err = cmntRepo.DeleteComment(rcomment.ID)
	a.NoError(err)

	reactions, err = reacRepo.GetTargetReactions(rcomment.ID)
	a.NoError(err)
	a.Len(reactions, 0)

	commentDetails, _ := cmntRepo.GetCommentByID(rpost.ID)
	a.Len(commentDetails.Comments, 0)
}

func TestCommentRepository_GetDetailedComments_Success(t *testing.T) {
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
		Author:     models.Author{ID: *userIds[0]},
	}
	rpost, _ := postRepo.AddPost(post)

	var commentID *string
	for i := range 5 {
		userIdx := rand.Intn(2)
		comment := &models.Comment{
			Content:         "this is a comment" + strconv.Itoa(i),
			Author:          models.Author{ID: *userIds[userIdx]},
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

	var count int64
	err := tx.Table("comments").
		Where("post_id = ? AND deleted_at IS NULL", rpost.ID).
		Count(&count).Error
	a.NoError(err)
	a.Equal(5, int(count))
}
