package users

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/zsomborjoel/workoutxz/internal/common"
)

type User struct {
	Id       string `db:"id"`
	UserName string `db:"username"`
	Email    string `db:"email"`
	Password string `db:"password"`
}

func FindUserByUserName(name string) (User, error) {
	log.Debug().Msg("users.FindUserByUserName called")

	db := common.GetDB()
	var u User
	err := db.Get(&u, "SELECT * FROM users WHERE username=$1", name)
	if err != nil {
		return u, fmt.Errorf("An error occured in users.FindUserByUserName.Get: %w", err)
	}
	return u, nil
}

func ExistByUserName(name string) (bool, error) {
	log.Debug().Msg("users.ExistByUserName called")

	db := common.GetDB()
	var i int
	err := db.Get(&i, "SELECT 1 FROM users WHERE username=$1", name)
	if err != nil {
		return false, fmt.Errorf("An error occured in users.ExistByUserName.Get: %w", err)
	}

	if i == 1 {
		return true, nil
	}

	return false, nil
}

func CreateOne(user User) error {
	log.Debug().Msg("users.CreateOne called")

	db := common.GetDB()
	tx := db.MustBegin()

	st := `INSERT INTO users (user_id, username, email, password) 
			VALUES (:user_id, :username, :email, :password)`
	_, err := tx.NamedExec(st, &user)
	if err != nil {
		return fmt.Errorf("An error occured in users.CreateOne.NamedExec: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("An error occured in users.CreateOne.Commit: %w", err)
	}

	return nil
}
