package server

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"liokor_mail/internal/pkg/common"
	mailDelivery "liokor_mail/internal/pkg/mail/delivery"
	mailRepository "liokor_mail/internal/pkg/mail/repository"
	mailUsecase "liokor_mail/internal/pkg/mail/usecase"
	userDelivery "liokor_mail/internal/pkg/user/delivery"
	userRepository "liokor_mail/internal/pkg/user/repository"
	userUsecase "liokor_mail/internal/pkg/user/usecase"
	"log"
	"os"
	"time"

	"liokor_mail/internal/app/server/middlewareHelpers"
)

func StartServer(config common.Config, quit chan os.Signal) {
	dbInstance, err := common.NewPostgresDataBase(config.DbString)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbInstance.Close()

	if config.Debug {
		log.Println("WARN: RUNNING IN THE DEBUG MODE! DON'T USE IN PRODUCTION!")
	}

	userRep := &userRepository.PostgresUserRepository{dbInstance}
	userUc := &userUsecase.UserUseCase{userRep, config}
	userHandler := userDelivery.UserHandler{userUc}

	mailRep := &mailRepository.PostgresMailRepository{dbInstance}
	mailUc := &mailUsecase.MailUseCase{mailRep, config}
	mailHander := mailDelivery.MailHandler{mailUc, userUc}

	e := echo.New()

	middlewareHelpers.SetupLogger(e, config.ApiLogPath)
	middlewareHelpers.SetupCSRFAndCORS(e, config.AllowedOrigin, config.Debug)

	e.Static("/media", "media")
	e.Static("/swagger", "swagger")

	e.POST("/user/auth", userHandler.Auth)
	e.DELETE("/user/session", userHandler.Logout)
	e.GET("/user", userHandler.Profile)
	e.POST("/user", userHandler.SignUp)
	e.PUT("/user/:username", userHandler.UpdateProfile)
	e.PUT("/user/:username/password", userHandler.ChangePassword)
	// e.GET("/user/:username", userHandler.ProfileByUsername)

	e.GET("/email/dialogues", mailHander.GetDialogues)
	e.GET("/email/emails", mailHander.GetEmails)
	e.POST("/email", mailHander.SendEmail)

	go func() {
		addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
		err := e.Start(addr)
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
