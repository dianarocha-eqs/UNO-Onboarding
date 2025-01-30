package utils

import (
	"api/internal/users/domain"
	"bytes"
	"encoding/json"
	"fmt"
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

		// if role has no value
		if len(userUUIDStr) == 0 {
			userUUID, err = uuid.FromString(userUUID.String())
			fmt.Println(userUUID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID format in header"})
				c.Abort()
				return
			}
		} else {
			userUUID, err = uuid.FromString(userUUIDStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID format in header"})
				c.Abort()
				return
			}
		}

		var role bool

		if roleStr == "" {
			role = domain.ROLE_USER
		} else {
			role, err = strconv.ParseBool(roleStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID format in header"})
				c.Abort()
				return
			}
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
