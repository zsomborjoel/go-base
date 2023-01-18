package verificationtoken

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"
	"github.com/zsomborjoel/workoutxz/internal/common"
	"github.com/zsomborjoel/workoutxz/internal/user"
)

type VerificationToken struct {
	Token     string `db:"token"`
	CreatedAt int64  `db:"created_at"`
	ExpiredAt int64  `db:"expired_at"`
	UserId    string `db:"user_id"`
}

const ExpirationTime = time.Hour * 24

func CreateOne(user user.User) error {
	log.Debug().Msg("verificationtokens.CreateOne called")

	uuid, err := uuid.NewV4()
	if err != nil {
		return fmt.Errorf("An error occured in verificationtokens.CreateOne.NewV4: %w", err)
	}

	now := time.Now()

	token := VerificationToken{
		Token:     uuid.String(),
		CreatedAt: now.Unix(),
		ExpiredAt: now.Add(ExpirationTime).Unix(),
		UserId:    user.Id,
	}

	db := common.GetDB()
	tx := db.MustBegin()

	st := `INSERT INTO verification_tokens (token, created_at, expired_at, user_id) 
			VALUES (:token, :created_at, :expired_at, :user_id)`
	_, err = tx.NamedExec(st, &token)
	if err != nil {
		return fmt.Errorf("An error occured in verificationtokens.CreateOne.NamedExec: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("An error occured in verificationtokens.CreateOne.Commit: %w", err)
	}

	return nil
}

func IsValid(token string) (VerificationToken, error) {
	log.Debug().Msg("verificationtokens.FindByToken called")

	db := common.GetDB()
	var vt VerificationToken
	err := db.Get(&vt, "SELECT * FROM verification_tokens WHERE token=$1", token)
	if err != nil {
		return vt, fmt.Errorf("An error occured in users.FindByToken.Get: %w", err)
	}

	if vt.ExpiredAt < time.Now().Unix() {
		return vt, fmt.Errorf("Verification token is expired")
	}

	return vt, nil
}
