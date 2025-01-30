package utils

import (
	"api/internal/users/domain"
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
		userUUIDStr := c.GetHeader("uuid")

		var requestBody struct {
			UUID uuid.UUID `json:"uuid"`
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

		// Parse JSON body
		if err := json.Unmarshal(bodyBytes, &requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON in request body"})
			c.Abort()
			return
		}

		var userUUID uuid.UUID
		// convert string to uuid
		userUUID, err = uuid.FromString(userUUIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID format in header"})
			c.Abort()
			return
		}

		var role bool
		// Convert string to boolean
		role, err = strconv.ParseBool(roleStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role format in header"})
			c.Abort()
			return
		}

		// Allow access if the user is an admin or the UUID matches
		if role == domain.ROLE_ADMIN || userUUID == requestBody.UUID {
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
		var roleStr string
		roleStr = c.GetHeader("role")

		var role bool
		var err error
		role, err = strconv.ParseBool(roleStr)

		if err != nil || role != domain.ROLE_ADMIN {
			c.JSON(http.StatusForbidden, gin.H{"error": "only admins can access this route"})
			c.Abort()
			return
		}

		// Get optional query parameters (from headers)
		search := c.GetHeader("search")
		sortDirectionStr := c.GetHeader("sortDirection")

		// Set default direction if nothing provided
		// Gives an error if the provided int is wrong (# 1, -1 )
		var sortDirection int
		if sortDirectionStr != "" {
			// Convert sortDirectionStr to int
			sortDirection, err = strconv.Atoi(sortDirectionStr)
			if err != nil || (sortDirection != 1 && sortDirection != -1) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid value for sorting direction"})
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
