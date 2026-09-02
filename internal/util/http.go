package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func CallExternalService(c *gin.Context, connectionID string) {
	url := os.Getenv("NOTIFICATION_URL")
	reqBody := map[string]string{
		"action_type":   "register",
		"connection_id": connectionID,
	}
	bodyBytes, _ := json.Marshal(reqBody)
	// 1. Create the request
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/connections", url), bytes.NewBuffer(bodyBytes))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 2. Set headers if needed
	req.Header.Set("Content-Type", "application/json")

	// 3. Execute the request using the default client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	// 4. Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 5. Return the result to the client
	c.JSON(resp.StatusCode, gin.H{"data": string(body)})
}
