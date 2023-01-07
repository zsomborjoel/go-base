package verificationtokens

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"
	"github.com/zsomborjoel/workoutxz/internal/common"
	"github.com/zsomborjoel/workoutxz/internal/users"
)

type VerificationToken struct {
	Token     string `db:"token"`
	CreatedAt int64 `db:"created_at"`
	ExpiredAt int64 `db:"expired_at"`
}

func CreateOne(user users.User) error {
	log.Debug().Msg("verificationtokens.CreateOne called")

	uuid, err := uuid.NewV4()
	if err != nil {
		return fmt.Errorf("An error occured in verificationtokens.CreateOne.NewV4: %w", err)
	}

	now := time.Now()

	token := VerificationToken{
		Token:     uuid.String(),
		CreatedAt: now.Unix(),
		ExpiredAt: now.Add(time.Hour * 24).Unix(),
	}

	db := common.GetDB()
	tx := db.MustBegin()

	st := `INSERT INTO verification_token (token, created_at, expired_at) 
			VALUES (:token, :created_at, :expired_at)`
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
