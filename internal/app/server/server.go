package server

import (
	"liokor_mail/internal/pkg/user"
	"liokor_mail/internal/pkg/user/delivery"
	"liokor_mail/internal/pkg/user/repository"
	"liokor_mail/internal/pkg/user/usecase"
	"github.com/labstack/echo/v4"
)

func CORSHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Access-Control-Allow-Origin", "https://mail.liokor.ru")
		return next(c)
	}
}

func StartServer(host string, port string) {
	userRep := &repository.UserRepository{
		map[string]user.User{},
		map[string]user.Session{},
	}
	userUc := &usecase.UserUseCase{userRep}
	userHandler := delivery.UserHandler{userUc}

	e := echo.New()
	e.Use(CORSHeader)

	e.POST("/user/auth", userHandler.Auth)
	e.DELETE("/user/session", userHandler.Logout)
	e.GET("/user", userHandler.Profile)
	e.POST("/user", userHandler.SignUp)
	e.PUT("/user/:username", userHandler.UpdateProfile)
	e.PUT("/user/:username/password", userHandler.ChangePassword)
	e.GET("/user/:username", userHandler.ProfileByUsername)

	e.Logger.Fatal(e.Start(host + ":" + port))
}
