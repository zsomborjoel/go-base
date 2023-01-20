package auth

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zsomborjoel/workoutxz/internal/user"
)

type UserClaim struct {
	jwt.RegisteredClaims
	user.User
}

func CreateJWTToken(user user.User) (string, error) {
	key := os.Getenv("JWT_KEY")


	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim{
		RegisteredClaims: jwt.RegisteredClaims{},
		User: user,
	})

	signedString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", fmt.Errorf("Error creating signed string: %v", err)
	}

	return signedString, nil
}
