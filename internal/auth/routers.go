package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/zsomborjoel/workoutxz/internal/common"
	"github.com/zsomborjoel/workoutxz/internal/email"
	"github.com/zsomborjoel/workoutxz/internal/users"
	"github.com/zsomborjoel/workoutxz/internal/verificationtokens"
)

func AuthRegister(r *gin.RouterGroup) {
	r.POST("/registration", Registration)
	r.GET(common.ConfirmRegistrationEndpoint, ConfirmRegistration)
}

func Registration(c *gin.Context) {
	log.Debug().Msg("Registration called")

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Errorf("Invalid body: %w", err))
		return
	}

	var rr RegistrationRequest
	err = json.Unmarshal(body, &rr)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Errorf("Unmarshal error occured: %w", err))
		return
	}

	var u users.User
	s := RegistrationRequestSerializer{c, rr}
	u, err = s.Model()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if err := users.CreateOne(u); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if err := verificationtokens.CreateOne(u); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if err := email.SendEmail(u.Email); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}

func ConfirmRegistration(c *gin.Context) {
	log.Debug().Msg("ConfirmRegistration called")

	t := c.Param("token")
	vt, err := verificationtokens.IsValid(t)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	err = users.ActivateOne(users.User{Id: vt.UserId})
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}
