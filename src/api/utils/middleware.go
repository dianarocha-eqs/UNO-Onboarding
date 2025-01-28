package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Middleware that restricts access to certain routes if not an admin
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetHeader("role")
		if role != "true" {
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
		role := c.GetHeader("role")
		userUUID := c.GetHeader("uuid")

		// Read the body, but do not consume it yet
		var requestBody struct {
			UUID string `json:"uuid"`
		}

		// Read the request body into a buffer so it can be reused later (if not, body is cleaned after the middleware reads the uuid)
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "unable to read request body"})
			c.Abort()
			return
		}

		// Store the body bytes in the context (to be used later)
		c.Set("requestBody", bodyBytes)

		// Reset the request body so that it can be read again by the handler
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// Bind the request body to the struct for easy validation
		if err := json.Unmarshal(bodyBytes, &requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON in request body"})
			c.Abort()
			return
		}

		// Allow access if the user is an admin or the UUID matches
		if role != "true" || userUUID == requestBody.UUID {
			c.Next()
		} else {
			c.JSON(http.StatusForbidden, gin.H{"error": "access forbidden"})
			c.Abort()
		}
	}
}
