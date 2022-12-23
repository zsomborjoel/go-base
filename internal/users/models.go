package users

import "github.com/zsomborjoel/workoutxz/internal/common"

type User struct {
	Id        int64  `db:"user_id"`
	FirstName string `db:"firstname"`
	LastName  string `db:"lastname"`
	UserName  string `db:"username"`
	Password  string `db:"password"`
}

func FindUserByUserName(name string) (User, error) {
	db := common.GetDB()
	var u User
	err := db.Get(&u, "SELECT * FROM users WHERE username=$1", name)
	return u, err
}
