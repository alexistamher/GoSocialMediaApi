package router

import (
	"github.com/alexistamher/social-api-go/internal/handler"
	"github.com/alexistamher/social-api-go/internal/service"
	"github.com/gin-gonic/gin"
)

func NewHandlers(authService service.AuthService) Handlers {
	return Handlers{Auth: handler.NewAuthHandler(authService)}
}

type Handlers struct {
	Auth *handler.AuthHandler
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
	}

	return r
}
