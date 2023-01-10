package email

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func EmailRegister(r *gin.RouterGroup) {
	r.POST("/send", Send)
}

func Send(c *gin.Context) {
	log.Debug().Msg("Send called")

	p := c.Query("target")

	err := SendEmail(p)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}
