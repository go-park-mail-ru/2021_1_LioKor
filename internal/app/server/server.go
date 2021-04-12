package server

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"liokor_mail/internal/pkg/common"
	userDelivery "liokor_mail/internal/pkg/user/delivery"
	userRepository "liokor_mail/internal/pkg/user/repository"
	userUsecase "liokor_mail/internal/pkg/user/usecase"
	mailDelivery "liokor_mail/internal/pkg/mail/delivery"
	mailUsecase "liokor_mail/internal/pkg/mail/usecase"
	mailRepository "liokor_mail/internal/pkg/mail/repository"
	"log"
	"os"
	"strconv"
	"time"
)

func StartServer(config common.Config, quit chan os.Signal) {
	dbInstance, err := common.NewPostgresDataBase(config.DbString)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbInstance.Close()

	userRep := &userRepository.PostgresUserRepository{dbInstance}
	userUc := &userUsecase.UserUseCase{userRep, config}
	userHandler := userDelivery.UserHandler{userUc}

	mailRep := &mailRepository.PostgresMailRepository{dbInstance}
	mailUC := &mailUsecase.MailUseCase{mailRep}
	mailHander := mailDelivery.MailHandler{mailUC, userUc}


	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.CSRF())
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

	e.GET("/api/dialogues", mailHander.GetDialogues)
	e.GET("/api/emails", mailHander.GetEmails)
	e.POST("/api/emails", mailHander.SendEmail)

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
