package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/zsomborjoel/workoutxz/internal/email"
	"github.com/zsomborjoel/workoutxz/internal/users"
	"github.com/zsomborjoel/workoutxz/internal/verificationtokens"
)

func AuthRegister(r *gin.RouterGroup) {
	r.POST("/signup", Registration)
}

func Registration(c *gin.Context) {
	log.Debug().Msg("Registration called")

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Errorf("Invalid body: %w", err).Error())
		return
	}

	var rr RegistrationRequest
	err = json.Unmarshal(body, &rr)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Errorf("Unmarshal error occured: %w", err).Error())
		return
	}

	var u users.User
	s := RegistrationRequestSerializer{c, rr}
	u, err = s.Model()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if err := users.CreateOne(u); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if err := verificationtokens.CreateOne(u); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	email.SendEmail("zsomborjoel@gmai.com")

	c.Writer.WriteHeader(http.StatusOK)
}
