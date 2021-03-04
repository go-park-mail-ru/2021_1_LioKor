package server

import (
	"lioKor_mail/internal/pkg/user"
	"lioKor_mail/internal/pkg/user/delivery"
	"lioKor_mail/internal/pkg/user/repository"
	"lioKor_mail/internal/pkg/user/usecase"
	"github.com/labstack/echo/v4"
)


func StartServer() {
	rep := &repository.UserRepository{
		map[string]user.User{},
		map[string]user.Session{},
	}
	uc := &usecase.UserUseCase{rep}
	userHandler := delivery.UserHandler{uc}
	e := echo.New()
	e.POST("/user/auth", userHandler.Auth)
	e.POST("/user/logout", userHandler.Logout)
	e.GET("/user", userHandler.Profile)
	e.POST("/user", userHandler.SignUp)
	e.PUT("/user/:username", userHandler.UpdateProfile)
	e.PUT("/user/:username/password", userHandler.ChangePassword)
	e.GET("/user/:username", userHandler.ProfileByUsername)
	e.Logger.Fatal(e.Start(":8000"))
}