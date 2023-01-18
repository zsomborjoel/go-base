package auth

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"
	"github.com/zsomborjoel/workoutxz/internal/user"
	"golang.org/x/crypto/bcrypt"
)

type RegistrationRequest struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegistrationRequestSerializer struct {
	C *gin.Context
	RegistrationRequest
}

func (s *RegistrationRequestSerializer) Model() (user.User, error) {
	log.Debug().Msg("auth.Model called")

	uuid, err := uuid.NewV4()
	if err != nil {
		return user.User{}, fmt.Errorf("An error occured in auth.Model.NewV4: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(s.RegistrationRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		return user.User{}, fmt.Errorf("An error occured in auth.Model.GenerateFromPassword: %w", err)
	}

	return user.User{
		Id:       uuid.String(),
		UserName: s.RegistrationRequest.UserName,
		Email:    s.RegistrationRequest.Email,
		Password: string(hash),
		Active : false,
	}, nil
}
