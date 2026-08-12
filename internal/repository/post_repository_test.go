package repository_test

import (
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
		Author:     models.Author{ID: *UserID},
	}

	rpost, _ := postRepo.AddPost(post)
	post = &models.Post{
		Content:    "this is a post test with parent",
		Visibility: models.Friends,
		Author:     models.Author{ID: *UserID},
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
			Author:     models.Author{ID: *UserID},
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

func TestPostRepository_AddPostReactions_Success(t *testing.T) {
	tx := db.Begin()
	t.Cleanup(func() {
		tx.Rollback()
	})
	a := assert.New(t)
	postRepo := repository.NewPostRepository(tx)
	authRepo := repository.NewAuthRepository(tx)
	reacRepo := repository.NewReactionRepository(tx)
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
		Author:     models.Author{ID: *userIds[0]},
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
		_, _ = reacRepo.AddReaction(rpost.ID, *userIds[userID], string(reactions[reactionIdx]), string(models.PostType))
	}
	postReactions, _ := reacRepo.GetTargetReactions(rpost.ID)
	a.Len(postReactions, 9)

	previewReactions, _ := reacRepo.GetTargetPreviewReactions(rpost.ID)
	prevCounter := 0
	for _, count := range previewReactions {
		prevCounter += count
	}
	a.Equal(prevCounter, len(postReactions))
}

func TestPostRepository_RemovePostSuccess(t *testing.T) {
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

	userId, _ := authRepo.Register(user)
	post := &models.Post{
		Content:    "This is a test post",
		Visibility: models.Public,
		Author:     models.Author{ID: *userId},
	}

	rpost, _ := postRepo.AddPost(post)
	posts, _, _ := postRepo.GetPostsByUserID(*userId, nil, uint(1))
	a.Len(posts, 1)
	err := postRepo.DeletePost(rpost.ID)
	a.NoError(err)
	posts, _, _ = postRepo.GetPostsByUserID(*userId, nil, uint(1))
	a.Len(posts, 0)
}

func TestPostRepository_RemoveReactionFromTarget(t *testing.T) {
	tx := db.Begin()
	t.Cleanup(func() {
		tx.Rollback()
	})
	a := assert.New(t)
	postRepo := repository.NewPostRepository(tx)
	authRepo := repository.NewAuthRepository(tx)
	reacRepo := repository.NewReactionRepository(tx)
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

	rpostReaction, _ := reacRepo.AddReaction(rpost.ID, *userId, string(models.LikeType), string(models.PostType))

	reactions, _ := reacRepo.GetTargetReactions(rpost.ID)
	a.Len(reactions, 1)
	err := reacRepo.DeleteReaction(rpostReaction.ID)
	a.NoError(err)
	reactions, _ = reacRepo.GetTargetReactions(rpost.ID)
	a.Len(reactions, 0)
}

func TestPostRepository_RemovePostAlongWithCommentsAndReactions_Success(t *testing.T) {
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
	comments, _ := postRepo.GetCommentsByPostID(rpost.ID)
	a.Len(comments, 0)

	comment := &models.Comment{
		Content:         "this is a comment",
		Author:          models.Author{ID: *userIds[1]},
		PostID:          rpost.ID,
		ParentCommentID: nil,
	}
	_, err := cmntRepo.AddComment(comment)
	a.NoError(err)

	reacRepo.AddReaction(rpost.ID, *userIds[0], "like", "post")
	reacRepo.AddReaction(rpost.ID, *userIds[1], "love", "post")
	reacRepo.AddReaction(rpost.ID, *userIds[2], "haha", "post")

	comments, err = postRepo.GetCommentsByPostID(rpost.ID)
	a.NoError(err)
	a.Len(comments, 1)

	reactions, err := reacRepo.GetTargetReactions(rpost.ID)
	a.NoError(err)
	a.Len(reactions, 3)

	err = postRepo.DeletePost(rpost.ID)
	a.NoError(err)

	reactions, err = reacRepo.GetTargetReactions(rpost.ID)
	a.NoError(err)
	a.Len(reactions, 0)

	comments, _ = postRepo.GetCommentsByPostID(rpost.ID)
	a.Len(comments, 0)
}

func TestPostRepository_GetDetailedPosts(t *testing.T) {
	tx := db.Begin()
	t.Cleanup(func() {
		tx.Rollback()
	})
	a := assert.New(t)
	postRepo := repository.NewPostRepository(tx)
	cmntRepo := repository.NewCommentRepository(tx)
	authRepo := repository.NewAuthRepository(tx)
	reacRepo := repository.NewReactionRepository(tx)

	usedIDs := make([]string, 10)

	for i := range 10 {
		user := &models.User{
			Username:    "JohnC" + strconv.Itoa(i),
			Password:    "passwordHash",
			Email:       "jconnor" + strconv.Itoa(i) + "@mail.com",
			DisplayName: "John Connor" + strconv.Itoa(i),
		}
		userId, _ := authRepo.Register(user)
		usedIDs[i] = *userId
	}

	post := &models.Post{
		Content:    "This is a test post",
		Visibility: models.Public,
		Author:     models.Author{ID: usedIDs[0]},
	}
	rpost, _ := postRepo.AddPost(post)

	for range 2 {
		comment := &models.Comment{
			Content:         "this is a comment",
			Author:          models.Author{ID: usedIDs[0]},
			PostID:          rpost.ID,
			ParentCommentID: nil,
		}
		cmntRepo.AddComment(comment)
	}

	reactions := []models.ReactionType{
		models.LikeType,
		models.LoveType,
		models.HahaType,
		models.WowType,
		models.SadType,
		models.AngryType,
	}
	for _, userID := range usedIDs {
		reaction := reactions[rand.Intn(len(reactions))]
		_, err := reacRepo.AddReaction(rpost.ID, userID, string(reaction), string(models.PostType))
		a.NoError(err)
	}

	offset := int(0)
	limit := int(20)
	posts, _, _ := postRepo.GetAllPosts(usedIDs[0], &offset, &limit)

	post = posts[0]
	a.Equal(post.Author.ID, usedIDs[0])

}

//TODO: TestPostRepository_GetPublicNotesFeed: this means return all public friend's posts
