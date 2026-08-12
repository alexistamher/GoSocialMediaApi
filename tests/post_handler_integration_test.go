package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexistamher/social-api-go/internal/handler"
	"github.com/alexistamher/social-api-go/internal/handler/dto"
	"github.com/alexistamher/social-api-go/internal/repository"
	"github.com/alexistamher/social-api-go/internal/router"
	"github.com/alexistamher/social-api-go/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestPostHandler_Integration_Flow(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	db, cleanup := setupAuthIntegrationDB(t)
	defer cleanup()

	authRepo := repository.NewAuthRepository(db)
	authSvc := service.NewAuthService(authRepo)
	postRepo := repository.NewPostRepository(db)
	postSvc := service.NewPostService(postRepo)
	cmntRepo := repository.NewCommentRepository(db)
	cmntSvc := service.NewCommentService(cmntRepo)
	reacRepo := repository.NewReactionRepository(db)
	reacSvc := service.NewReactionService(reacRepo)

	handlers := router.Handlers{
		Auth:     handler.NewAuthHandler(authSvc),
		Post:     handler.NewPostHandler(postSvc),
		Comment:  handler.NewCommentHandler(cmntSvc),
		Reaction: handler.NewReactionHandler(reacSvc),
	}

	r := router.New(handlers, handler.AuthMiddleware())

	var creds dto.AuthResponse
	var postID string
	var commentID string

	t.Run("1. Register user", func(t *testing.T) {
		regReq := dto.RegisterRequest{
			Username:    "integrationUser",
			Email:       "integration@example.com",
			Password:    "password123",
			DisplayName: "Integration User",
		}

		body, _ := json.Marshal(regReq)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		err := json.Unmarshal(w.Body.Bytes(), &creds)
		assert.NoError(t, err)
	})

	t.Run("2. Create post and get posts", func(t *testing.T) {
		reqNewPost := dto.CreatePostRequest{
			Content:    "some test content",
			Visibility: "public",
		}
		body, _ := json.Marshal(reqNewPost)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/posts", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+creds.AccessToken)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		var res dto.PostResponse
		err := json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)
		postID = res.ID

		req = httptest.NewRequest(http.MethodGet, "/api/v1/posts", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+creds.AccessToken)
		w = httptest.NewRecorder()

		r.ServeHTTP(w, req)
		var postRes map[string]any
		err = json.Unmarshal(w.Body.Bytes(), &postRes)
		nextCursor, _ := postRes["next_cursor"].(float64)
		posts, _ := postRes["posts"].([]any)

		assert.NoError(t, err)
		assert.NotZero(t, nextCursor)
		assert.NotEmpty(t, posts)
	})

	t.Run("3. Insert and get post comments", func(t *testing.T) {
		reqNewComment := dto.AddCommentRequest{
			Content: "some test content",
			PostID:  postID,
		}
		body, _ := json.Marshal(reqNewComment)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/comments", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+creds.AccessToken)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		var res dto.CommentResponse
		err := json.Unmarshal(w.Body.Bytes(), &res)

		assert.NoError(t, err)
		assert.NotEmpty(t, res.Author.ID)
		assert.NotEmpty(t, res.Author.Username)
		assert.NotEmpty(t, res.Author.DisplayName)

		commentID = res.ID

		req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/comments/post/%s", postID), nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+creds.AccessToken)
		w = httptest.NewRecorder()

		r.ServeHTTP(w, req)

		var commentsRes []dto.CommentResponse
		err = json.Unmarshal(w.Body.Bytes(), &commentsRes)

		assert.NoError(t, err)
		assert.Len(t, commentsRes, 1)
	})

	t.Run("4. Insert comment into another comment and get comments inside the comment", func(t *testing.T) {
		reqNewComment := dto.AddCommentRequest{
			Content:         "some test content",
			PostID:          postID,
			ParentCommentID: &commentID,
		}
		body, _ := json.Marshal(reqNewComment)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/comments", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+creds.AccessToken)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		var res dto.CommentResponse
		err := json.Unmarshal(w.Body.Bytes(), &res)

		assert.NoError(t, err)
		assert.NotEmpty(t, res.Author.ID)
		assert.NotEmpty(t, res.Author.Username)
		assert.NotEmpty(t, res.Author.DisplayName)
		assert.NotNil(t, res.ParentCommentID)

		req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/comments/%s", commentID), nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+creds.AccessToken)
		w = httptest.NewRecorder()

		r.ServeHTTP(w, req)

		var commentsRes dto.CommentDetailsResponse
		err = json.Unmarshal(w.Body.Bytes(), &commentsRes)

		assert.NoError(t, err)
		assert.Len(t, commentsRes.Comments, 1)
	})

	t.Run("5. Insert reaction at post and then retrieve it", func(t *testing.T) {
		postReaction := dto.AddReactionRequest{
			TargetID:           postID,
			ReactionType:       "haha",
			ReactionTargetType: "post",
		}
		body, _ := json.Marshal(postReaction)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/reactions", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+creds.AccessToken)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		var postReacRes dto.CreatedReactionResponse
		err := json.Unmarshal(w.Body.Bytes(), &postReacRes)

		assert.NoError(t, err)
		assert.NotEmpty(t, postReacRes.ID)
		assert.NotEmpty(t, postReacRes.CreatedAt)

		req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/posts/%s", postID), nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+creds.AccessToken)
		w = httptest.NewRecorder()

		r.ServeHTTP(w, req)

		var postRes *dto.PostWithDetailsResponse
		err = json.Unmarshal(w.Body.Bytes(), &postRes)

		assert.NoError(t, err)
		assert.NotEmpty(t, postRes.Reactions)

		reaction := postRes.Reactions[0]

		assert.NotNil(t, reaction)
		assert.NotEmpty(t, reaction.ID)
		assert.NotEmpty(t, reaction.ReactionType)
		assert.NotEmpty(t, reaction.TargetID)
		assert.NotEmpty(t, reaction.CreatedAt)
	})

	t.Run("6. Insert reaction at comment and then retrieve it", func(t *testing.T) {
		cmntReaction := dto.AddReactionRequest{
			TargetID:           commentID,
			ReactionType:       "like",
			ReactionTargetType: "comment",
		}
		body, _ := json.Marshal(cmntReaction)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/reactions", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+creds.AccessToken)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		var postReacRes dto.CreatedReactionResponse
		err := json.Unmarshal(w.Body.Bytes(), &postReacRes)

		assert.NoError(t, err)
		assert.NotEmpty(t, postReacRes.ID)
		assert.NotEmpty(t, postReacRes.CreatedAt)

		req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/comments/%s", commentID), nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+creds.AccessToken)
		w = httptest.NewRecorder()

		r.ServeHTTP(w, req)

		var postRes *dto.CommentDetailsResponse
		err = json.Unmarshal(w.Body.Bytes(), &postRes)

		assert.NoError(t, err)
		assert.NotEmpty(t, postRes.Reactions)

		reaction := postRes.Reactions[0]

		assert.NotNil(t, reaction)
		assert.NotEmpty(t, reaction.ID)
		assert.NotEmpty(t, reaction.ReactionType)
		assert.NotEmpty(t, reaction.CreatedAt)
	})
}
