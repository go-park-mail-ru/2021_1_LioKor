package server

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"liokor_mail/internal/pkg/user/delivery"
	"liokor_mail/internal/pkg/user/repository"
	"liokor_mail/internal/pkg/user/usecase"
	"log"
	"os"
	"strconv"
	"time"
)

func StartServer(host string, port int, allowedOrigins []string, quit chan os.Signal, dbConfig string) {
	userRep, err := repository.NewPostgresUserRepository(dbConfig)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer userRep.Close()

	userUc := &usecase.UserUseCase{userRep}
	userHandler := delivery.UserHandler{userUc}

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: allowedOrigins,
		AllowCredentials: true,
	}))

	e.POST("/user/auth", userHandler.Auth)
	e.DELETE("/user/session", userHandler.Logout)
	e.GET("/user", userHandler.Profile)
	e.POST("/user", userHandler.SignUp)
	e.PUT("/user/:username", userHandler.UpdateProfile)
	e.PUT("/user/:username/password", userHandler.ChangePassword)
	e.GET("/user/:username", userHandler.ProfileByUsername)

	go func() {
		if err := e.Start(host + ":" + strconv.Itoa(port)); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
