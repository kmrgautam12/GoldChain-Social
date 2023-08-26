package errorservice

import "github.com/gin-gonic/gin"

func ErrorResponse(c *gin.Context, code int, message string) {
	response := map[string]interface{}{
		"error": map[string]interface{}{
			"code":    code,
			"message": message,
			"status":  "FAILED",
		},
	}
	c.JSON(code, response)
}
