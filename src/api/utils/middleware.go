package utils

import (
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
		userUUIDStr := c.GetHeader("uuid")

		// Read the request body into a buffer
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
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // Reset body

		// Parse JSON body
		if err := json.Unmarshal(bodyBytes, &requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON in request body"})
			c.Abort()
			return
		}

		var userUUID uuid.UUID
		if len(userUUIDStr) == 0 {
			userUUID, err = uuid.FromString(userUUIDStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID format in header"})
				c.Abort()
				return
			}
		}

		requestUUID, err := uuid.FromString(requestBody.UUID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID format in request body"})
			c.Abort()
			return
		}

		// Allow access if the user is an admin or the UUID matches
		if role == "true" || userUUID == requestUUID {
			c.Next()
		} else {
			c.JSON(http.StatusForbidden, gin.H{"error": "access forbidden"})
			c.Abort()
		}
	}
}

// Middleware that parses sortDirection and search, and checks if the user is admin
func AdminAndSortMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if user is admin
		role := c.GetHeader("role")
		if role != "true" {
			c.JSON(http.StatusForbidden, gin.H{"error": "only admins can access this route"})
			c.Abort()
			return
		}

		// Get optional query parameters (from headers)
		search := c.GetHeader("search")
		sortDirectionStr := c.GetHeader("sortDirection")
		var sortDirection int
		var err error

		if sortDirectionStr != "" {
			// Convert sortDirectionStr to int
			sortDirection, err = strconv.Atoi(sortDirectionStr)
			if err != nil || (sortDirection != 1 && sortDirection != -1) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sortDirection"})
				c.Abort()
				return
			}
		} else {
			sortDirection = 0 // Default to no sorting
		}

		c.Set("search", search)
		c.Set("sortDirection", sortDirection)
		c.Next()
	}
}
