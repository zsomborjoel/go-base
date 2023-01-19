package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/zsomborjoel/workoutxz/internal/common"
	"github.com/zsomborjoel/workoutxz/internal/email"
	"github.com/zsomborjoel/workoutxz/internal/user"
	"github.com/zsomborjoel/workoutxz/internal/verificationtoken"
)

func AuthRegister(r *gin.RouterGroup) {
	r.POST("/registration", Registration)
	r.GET(common.ConfirmRegistrationEndpoint, ConfirmRegistration)
	r.PUT("/resend-verification", ResendVerification)
}

func Registration(c *gin.Context) {
	log.Debug().Msg("Registration called")

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("Invalid body: %w", err))
		return
	}

	var rr RegistrationRequest
	err = json.Unmarshal(body, &rr)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("Unmarshal error occured: %w", err))
		return
	}

	var u user.User
	s := RegistrationRequestSerializer{c, rr}
	u, err = s.Model()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	err = user.ExistByUserName(u.UserName)
	if err != nil  {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if err := user.CreateOne(u); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	t, err := verificationtoken.CreateOne(u)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if err := email.SendEmail(u.Email, t); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}

func ConfirmRegistration(c *gin.Context) {
	log.Debug().Msg("ConfirmRegistration called")

	t := c.Query("token")

	if t == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("Invalid token"))
		return
	}

	vt, err := verificationtoken.IsValid(t)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	err = user.ActivateOne(user.User{Id: vt.UserId})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	verificationtoken.DeleteOne(t)

	c.Writer.WriteHeader(http.StatusOK)
}

func ResendVerification(c *gin.Context) {
	log.Debug().Msg("ConfirmRegistration called")

	t := c.Query("token")

	if t == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("Invalid token"))
		return
	}

	_, err := verificationtoken.IsValid(t)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	new, err := verificationtoken.UpdateToken(t)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, new)
}
