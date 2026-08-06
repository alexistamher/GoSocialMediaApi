package handler

import (
	"errors"
	"net/http"

	derrors "github.com/alexistamher/social-api-go/internal/domain"
	"github.com/gin-gonic/gin"
)

func respondError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, derrors.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, derrors.ErrAlreadyExists):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case errors.Is(err, derrors.ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	case errors.Is(err, derrors.ErrUnauthorized):
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	case errors.Is(err, derrors.ErrForbidden):
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	case errors.Is(err, derrors.ErrMissingUserID):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
