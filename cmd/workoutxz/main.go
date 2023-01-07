package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/zsomborjoel/workoutxz/internal/auth"
	"github.com/zsomborjoel/workoutxz/internal/common"
	"github.com/zsomborjoel/workoutxz/internal/email"
	"github.com/zsomborjoel/workoutxz/internal/pings"
	"github.com/zsomborjoel/workoutxz/internal/users"
)

func main() {
	common.LoadEnvVariables()

	level := os.Getenv("LOG_LEVEL")
	zerolog.SetGlobalLevel(common.LogLevel(level))

	common.Init()

	r := gin.Default()
	r.Use(common.CORSMiddleware())

	v1 := r.Group("/api")
	pings.PingRegister(v1.Group("/ping"))
	users.UsersRegister(v1.Group("/users"))
	auth.AuthRegister(v1.Group("/auth"))
	email.EmailRegister(v1.Group("/email"))

	r.Run()
}
