package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func UsersRegister(r *gin.RouterGroup) {
	r.GET("/:username", UserRetrieveByUserName)
}

func UserRetrieveByUserName(c *gin.Context) {
	un := c.Param("username")
	u, err := FindUserByUserName(un)
	if err != nil {
		c.JSON(http.StatusNotFound, err)
	}

	s := UserSerializer{c, u}
	r, err := s.Response()
	if err != nil {
		c.JSON(http.StatusNotFound, err)
	}

	c.JSON(http.StatusOK, r)
}
