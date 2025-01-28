package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Middleware that restricts access to certain routes if not an admin
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
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
		userUUID, uuidExists := c.Get("uuid") // Assuming UUID is set in context via token middleware
		var requestBody struct {
			ID string `json:"uuid"`
		}

		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			c.Abort()
			return
		}

		if !uuidExists || requestBody.ID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID missing"})
			c.Abort()
			return
		}

		// Allow access if admin or if UUID matches
		if role == "admin" || userUUID == requestBody.ID {
			c.Next()
		} else {
			c.JSON(http.StatusForbidden, gin.H{"error": "access forbidden"})
			c.Abort()
		}
	}
}
