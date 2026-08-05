package handler

import "github.com/gin-gonic/gin"

type ContextKey = string

const ContextKeyUserId ContextKey = "user_id"

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken := c.Request.Header.Get("Authorization")

		//TODO: Agregar logica para verificar el token

		c.Set(ContextKeyUserId, authToken)
		c.Next()
	}
}
