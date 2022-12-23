package pings

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zsomborjoel/workoutxz/internal/common"
)

func PingRegister(router *gin.RouterGroup) {
	router.GET("", ping)
	router.GET("/db", pingDb)
}

func ping(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusOK)
}

func pingDb(c *gin.Context) {
	db := common.GetDB()

	var result int
	err := db.Get(&result, "SELECT 1")
	if (err != nil) {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	
	c.Writer.WriteHeader(http.StatusOK)
}
