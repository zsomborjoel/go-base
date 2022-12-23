package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/zsomborjoel/workoutxz/internal/common"
	"github.com/zsomborjoel/workoutxz/internal/pings"
	"github.com/zsomborjoel/workoutxz/internal/users"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	common.LoadEnvVariables()

	common.Init()

	r := gin.Default()
	r.Use(common.CORSMiddleware())

	v1 := r.Group("/api")
	pings.PingRegister(v1.Group("/ping"))
	users.UsersRegister(v1.Group("/users"))

	r.Run()
}
