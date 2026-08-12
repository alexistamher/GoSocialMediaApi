package router

import (
	"github.com/alexistamher/social-api-go/internal/handler"
	"github.com/alexistamher/social-api-go/internal/service"
	"github.com/gin-gonic/gin"
)

func NewHandlers(authService service.AuthService, postService service.PostService, commentService service.CommentService,
	reactionService service.ReactionService) Handlers {
	return Handlers{
		Auth:     handler.NewAuthHandler(authService),
		Post:     handler.NewPostHandler(postService),
		Comment:  handler.NewCommentHandler(commentService),
		Reaction: handler.NewReactionHandler(reactionService),
	}
}

type Handlers struct {
	Auth     *handler.AuthHandler
	Post     *handler.PostHandler
	Comment  *handler.CommentHandler
	Reaction *handler.ReactionHandler
}

func New(h Handlers, authMiddleware gin.HandlerFunc) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	public := r.Group("/api/v1")
	public.POST("/auth/register", h.Auth.Register)
	public.POST("/auth/login", h.Auth.Login)

	private := r.Group("/api/v1")
	private.Use(authMiddleware)
	{
		private.GET("/auth/info", h.Auth.GetInfo)

		private.POST("/posts", h.Post.CreatePost)
		private.GET("/posts", h.Post.GetUserPosts)
		private.GET("/posts/:post_id", h.Post.GetPostByID)
		private.DELETE("/posts/:post_id", h.Post.DeletePost)

		private.POST("/comments", h.Comment.AddComment)
		private.DELETE("/comments/:comment_id", h.Comment.DeleteComment)
		private.GET("/comments/:comment_id", h.Comment.GetCommentById)
		private.GET("/comments/post/:post_id", h.Comment.GetPostComments)

		private.POST("/reactions", h.Reaction.AddReaction)
		private.PUT("/reactions", h.Reaction.UpdateReaction)
		private.GET("/reactions/:target_id", h.Reaction.GetTargetReactions)
		private.DELETE("/reactions/:target_id", h.Reaction.DeleteReaction)
	}

	return r
}
