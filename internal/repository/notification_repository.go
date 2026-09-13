package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	drepository "github.com/alexistamher/social-api-go/internal/domain/repository"
	"github.com/gin-gonic/gin"
)

type notificationRepository struct {
}

func NewNotificationRepository() drepository.NotificationRepository {
	return &notificationRepository{}
}

func (n *notificationRepository) RegisterConnection(c *gin.Context, connectionID string) {
	url := os.Getenv("NOTIFICATION_URL")
	reqBody := map[string]string{
		"action_type":   "register",
		"connection_id": connectionID,
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/connections", url), bytes.NewBuffer(bodyBytes))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(resp.StatusCode, gin.H{"data": string(body)})
}
