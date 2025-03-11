package utils

import (
	auth_service "api/internal/auth/usecase"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	uuid "github.com/tentone/mssql-uuid"

	"github.com/gin-gonic/gin"
)

// Middleware that restricts access to certain routes if not an admin
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleStr := c.GetHeader("role")
		role, err := strconv.ParseBool(roleStr)
		if err != nil || !role {
			c.JSON(http.StatusForbidden, gin.H{"error": "only admins can access this route"})
			c.Abort()
			return
		}
		c.Set("role", role)
		c.Next()
	}
}

// Middleware that allows access to the route only for admins or the user themselves
func AdminAndUserItself() gin.HandlerFunc {
	return func(c *gin.Context) {
		var roleStr = c.GetHeader("role")
		var userUUIDStr = c.GetHeader("uuid")

		var requestBody struct {
			UUID uuid.UUID `json:"uuid"`
		}

		var bodyBytes, err = io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "unable to read request body"})
			c.Abort()
			return
		}

		// Store the body for future use (since it deletes once it reads it)
		c.Set("requestBody", bodyBytes)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// Parse JSON body
		if err := json.Unmarshal(bodyBytes, &requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON in request body"})
			c.Abort()
			return
		}

		var userUUID uuid.UUID

		// if userUUIDStr is empty
		if len(userUUIDStr) == 0 {
			// Returns a nil UUID
			userUUID = uuid.FromStringOrNil(userUUIDStr)
		} else {
			// changes from string to uuid
			userUUID, err = uuid.FromString(userUUIDStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID format in header"})
				c.Abort()
				return
			}
		}

		var role bool

		if len(roleStr) == 0 {
			// which means it will not have access (unless it's the own user)
			role = false
		} else {
			// changes string value to bool
			role, err = strconv.ParseBool(roleStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID format in header"})
				c.Abort()
				return
			}
		}

		// Allow access if the user is an admin or the UUID matches
		if !role || userUUID == requestBody.UUID {
			c.Set("role", role)
			c.Set("uuid", userUUID)
			c.Next()
		} else {
			c.JSON(http.StatusForbidden, gin.H{"error": "access forbidden"})
			c.Abort()
		}
	}
}

// Validates the JWT token to authorize access
func AuthMiddleware(authService auth_service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token from Authorization header
		var authHeader = c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization token required"})
			c.Abort()
			return
		}

		var tokenStr = authHeader

		// Validate the token using the service
		var isValid, err = authService.IsTokenValid(c.Request.Context(), tokenStr)
		if err != nil || !isValid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("token", tokenStr)
		c.Next()
	}
}
