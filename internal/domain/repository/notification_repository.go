package repository

import "github.com/gin-gonic/gin"

type NotificationRepository interface {
	RegisterConnection(ctx *gin.Context, connectionID string)
}
