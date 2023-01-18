package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		var finalErr error
		for _, err := range c.Errors {
			finalErr = err
			log.Error().Err(err).Msg("http error")
		}

		// status -1 doesn't overwrite existing status code
		c.JSON(-1, finalErr)
	}
}
