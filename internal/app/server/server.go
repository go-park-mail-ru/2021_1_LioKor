package server

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"liokor_mail/internal/pkg/common"
	"liokor_mail/internal/pkg/user"
	"liokor_mail/internal/pkg/user/delivery"
	"liokor_mail/internal/pkg/user/repository"
	"liokor_mail/internal/pkg/user/usecase"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

func StartServer(config common.Config, quit chan os.Signal) {
	userRep := &repository.UserRepository{
		repository.UserStruct{map[string]user.User{}, sync.Mutex{}},
		repository.SessionStruct{map[string]user.Session{}, sync.Mutex{}},
	}
	userUc := &usecase.UserUseCase{userRep, config}
	userHandler := delivery.UserHandler{userUc}

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     config.AllowedOrigins,
		AllowCredentials: true,
	}))
	e.Static("/media", "media")

	e.POST("/user/auth", userHandler.Auth)
	e.DELETE("/user/session", userHandler.Logout)
	e.GET("/user", userHandler.Profile)
	e.POST("/user", userHandler.SignUp)
	e.PUT("/user/:username", userHandler.UpdateProfile)
	e.PUT("/user/:username/password", userHandler.ChangePassword)
	e.GET("/user/:username", userHandler.ProfileByUsername)

	go func() {
		err := e.Start(config.Host + ":" + strconv.Itoa(config.Port))
		if err != nil {
			log.Println("Server was shut down with no errors!")
		} else {
			log.Fatal("Error occured while trying to shut down server: " + err.Error())
		}
	}()
	<-quit

	log.Println("Interrupt signal received. Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Fatal("Server shut down timeout with an error: " + err.Error())
	}
}
