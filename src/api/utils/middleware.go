package utils

import (
	"api/internal/users/domain"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Middleware that restricts access to certain routes if not an admin
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleStr := c.GetHeader("role")
		role, err := strconv.ParseBool(roleStr)
		if err != nil || role != domain.ROLE_ADMIN {
			c.JSON(http.StatusForbidden, gin.H{"error": "only admins can access this route"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// Middleware that allows access to the route only for admins or the user themselves
func AdminAndUserItself() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleStr := c.GetHeader("role")
		userUUID := c.GetHeader("uuid")

		var requestBody struct {
			UUID string `json:"uuid"`
		}

		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "unable to read request body"})
			c.Abort()
			return
		}

		// Store the body for future use (since it deletes once it reads it)
		c.Set("requestBody", bodyBytes)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		if err := json.Unmarshal(bodyBytes, &requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON in request body"})
			c.Abort()
			return
		}
		role, err := strconv.ParseBool(roleStr)

		if err != nil || role == domain.ROLE_ADMIN || userUUID == requestBody.UUID {
			c.Next()
		} else {
			c.JSON(http.StatusForbidden, gin.H{"error": "access forbidden"})
			c.Abort()
		}
	}
}
