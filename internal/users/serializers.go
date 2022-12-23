package users

import "github.com/gin-gonic/gin"

type UserResponse struct {
	Id        int64
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	UserName  string `json:"userName"`
}

type UserSerializer struct {
	C *gin.Context
	User
}

func (s *UserSerializer) Response() (UserResponse, error) {
	r := UserResponse{
		Id:        s.User.Id,
		FirstName: s.User.FirstName,
		LastName:  s.User.LastName,
		UserName:  s.User.UserName,
	}

	return r, nil
}
