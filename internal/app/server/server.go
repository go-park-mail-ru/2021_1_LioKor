package server

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"liokor_mail/internal/pkg/common"
	mailDelivery "liokor_mail/internal/pkg/mail/delivery"
	mailRepository "liokor_mail/internal/pkg/mail/repository"
	mailUsecase "liokor_mail/internal/pkg/mail/usecase"
	userDelivery "liokor_mail/internal/pkg/user/delivery"
	userRepository "liokor_mail/internal/pkg/user/repository"
	userUsecase "liokor_mail/internal/pkg/user/usecase"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
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
	mailUC := &mailUsecase.MailUseCase{mailRep, config}
	mailHander := mailDelivery.MailHandler{mailUC, userUc}

	e := echo.New()

	if len(config.ApiLogPath) > 0 {
		logFile, err := os.Create(config.ApiLogPath)
		if err == nil {
			defer logFile.Close()
			e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
				Output: logFile,
			}))
			log.Printf("INFO: Logging API calls to %s\n", config.ApiLogPath)
		} else {
			log.Println("WARN: Unable to create log file!")
		}
	} else {
		log.Println("WARN: Logging disabled")
	}

	if len(config.AllowedOrigin) > 0 {
		url, err := url.Parse(config.AllowedOrigin)
		if err != nil {
			log.Fatal(err)
			return
		}
		csrfCookieDomain := url.Hostname()
		if len(csrfCookieDomain) == 0 {
			log.Fatal("Invalid domain specified in allowedOrigin")
			return
		}

		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins:     []string{config.AllowedOrigin},
			AllowCredentials: true,
		}))

		e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
			Skipper: func(c echo.Context) bool {
				host := c.Request().Host
				if strings.HasPrefix(host, "localhost:") || host == "localhost" {
					return true
				}
				return false
			},
			CookieSameSite: http.SameSiteStrictMode,
			CookieDomain:   csrfCookieDomain,
			CookiePath:     "/",
		}))
		log.Printf("INFO: %s added to CORS and CSRF protection enabled for %s\n", config.AllowedOrigin, csrfCookieDomain)
	} else {
		log.Println("WARN: allowedOrigin not set in config => CORS and CSRF middlewares are not enabled!")
	}

	e.Static("/media", "media")
	e.Static("/swagger", "swagger")

	e.POST("/user/auth", userHandler.Auth)
	e.DELETE("/user/session", userHandler.Logout)
	e.GET("/user", userHandler.Profile)
	e.POST("/user", userHandler.SignUp)
	e.PUT("/user/:username", userHandler.UpdateProfile)
	e.PUT("/user/:username/password", userHandler.ChangePassword)
	e.GET("/user/:username", userHandler.ProfileByUsername)

	e.GET("/email/dialogues", mailHander.GetDialogues)
	e.GET("/email/emails", mailHander.GetEmails)
	e.POST("/email", mailHander.SendEmail)

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
