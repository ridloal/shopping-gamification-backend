package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ValidateRequest(input interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := c.ShouldBindJSON(input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request format",
				"details": err.Error(),
			})

			c.Abort()
			return
		}

		validate := validator.New()
		if err := validate.Struct(input); err != nil {
			errors := err.(validator.ValidationErrors)
			errMessages := make(map[string]string)

			for _, e := range errors {
				errMessages[e.Field()] = fmt.Sprintf("Field validation for '%s' failed on the '%s' tag", e.Field(), e.Tag())
			}

			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation failed",
				"details": errMessages,
			})
			c.Abort()
			return
		}

		c.Set("validated_input", input)
		c.Next()
	}
}
