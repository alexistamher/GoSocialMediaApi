package repository_test

import (
	"encoding/json"
	"math/rand"
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

func TestPostRepository_AddCommentToPost_Success(t *testing.T) {
	tx := db.Begin()
	t.Cleanup(func() {
		tx.Rollback()
	})
	a := assert.New(t)
	postRepo := repository.NewPostRepository(tx)
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
		rcomment, _ := postRepo.AddComment(comment)
		if rand.Intn(2) == 0 {
			commentID = &rcomment.ID
		} else {
			commentID = nil
		}

	}

	comments, _ := postRepo.GetCommentsByPostID(rpost.ID)

	jcomments, _ := json.Marshal(comments)
	print(string(jcomments))
	a.Len(comments, 5)
}

func TestPostRepository_AddPostReactions_Success(t *testing.T) {
	tx := db.Begin()
	t.Cleanup(func() {
		tx.Rollback()
	})
	a := assert.New(t)
	postRepo := repository.NewPostRepository(tx)
	authRepo := repository.NewAuthRepository(tx)
	var userIds []*string

	for i := range 9 {
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

	var reactions = []models.ReactionType{
		models.LikeType,
		models.LoveType,
		models.HahaType,
		models.WowType,
		models.SadType,
		models.AngryType,
	}
	for userID := range userIds {
		reactionIdx := rand.Intn(6)
		_, _ = postRepo.AddPostReaction(rpost.ID, *userIds[userID], string(reactions[reactionIdx]))
	}
	postReactions, _ := postRepo.GetTargetReactions(rpost.ID)
	a.Len(postReactions, 9)

	previewReactions, _ := postRepo.GetTargetPreviewReactions(rpost.ID)
	prevCounter := 0
	for _, count := range previewReactions {
		prevCounter += count
	}
	a.Equal(prevCounter, len(postReactions))
}

//TODO: TestPostRepository_RemoveReaction from target (post/comment)

//TODO: TestPostRepository_GetPublicNotesFeed: this means return all public friend's posts
