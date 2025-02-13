package version

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const apiVersion = "0.0.0"

// Returns the API version and timestamp in ISO 8601 format
func GetVersion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version":   apiVersion,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}
