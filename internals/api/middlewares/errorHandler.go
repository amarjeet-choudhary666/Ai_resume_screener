package middlewares

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Internal server error",
					"message": "Something went wrong on our end",
				})
				c.Abort()
			}
		}()

		c.Next()

		// Handle errors after request processing
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				log.Printf("Request error: %v", err.Err)
			}

			// Return the last error as JSON
			lastError := c.Errors.Last()
			if lastError != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Request processing error",
					"message": lastError.Error(),
				})
			}
		}
	}
}

func ValidationErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check for validation errors
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				if err.Type == gin.ErrorTypeBind {
					c.JSON(http.StatusBadRequest, gin.H{
						"error":   "Validation error",
						"message": err.Err.Error(),
					})
					c.Abort()
					return
				}
			}
		}
	}
}
