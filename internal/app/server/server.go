package server

import (
	"context"
	"fmt"
	echoPrometheus "github.com/globocom/echo-prometheus"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
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

	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"

	"liokor_mail/internal/app/server/middlewareHelpers"
	session "liokor_mail/internal/pkg/common/protobuf_sessions"
)

func GetPrivateKey(path string) (*rsa.PrivateKey, error) {
	keyString, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode([]byte(keyString))
	if block == nil {
		return nil, err
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func StartServer(config common.Config, quit chan os.Signal) {
	dbInstance, err := common.NewPostgresDataBase(config)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbInstance.Close()

	if config.Debug {
		log.Println("WARN: RUNNING IN THE DEBUG MODE! DON'T USE IN PRODUCTION!")
	}

	grpcConn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", config.AuthHost, config.AuthPort),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("Unable to connect to grpc: %v\n", err)
	}

	defer grpcConn.Close()

	sessManager := session.NewIsAuthClient(grpcConn)

	userRep := &userRepository.PostgresUserRepository{dbInstance}
	userUc := &userUsecase.UserUseCase{userRep, sessManager, config}
	userHandler := userDelivery.UserHandler{userUc}

	privateKey, err := GetPrivateKey(config.DkimPrivateKeyPath)
	if err != nil {
		log.Printf("WARN: Unable to load private key: %v", err)
		privateKey = nil
	} else {
		log.Println("INFO: Private key for DKIM successfully loaded!")
	}
	mailRep := &mailRepository.PostgresMailRepository{dbInstance}
	mailUC := &mailUsecase.MailUseCase{mailRep, config, privateKey}
	mailHander := mailDelivery.MailHandler{mailUC}

	e := echo.New()

	var configMetrics = echoPrometheus.NewConfig()
	configMetrics.Buckets = []float64{
		0.001, // 1ms
		0.01,  // 10ms
		0.05,  // 50ms
		0.1,   // 100ms
		0.25,  // 250ms
		0.5,   // 500ms
		1,     // 1s
	}
	e.Use(echoPrometheus.MetricsMiddlewareWithConfig(configMetrics))
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	isAuth := middlewareHelpers.AuthMiddleware{userUc, sessManager}

	middlewareHelpers.SetupLogger(e, config.ApiLogPath)
	middlewareHelpers.SetupCSRFAndCORS(e, config.AllowedOrigin, config.Debug)

	//p := prometheus.NewPrometheus("echo", nil)
	//p.Use(e)

	if config.Debug {
		e.Static("/media", "media")
		e.Static("/swagger", "swagger")
	}

	e.POST("/user/auth", userHandler.Auth)
	e.DELETE("/user/session", userHandler.Logout, isAuth.IsAuth)
	e.GET("/user", userHandler.Profile, isAuth.IsAuth)
	e.POST("/user", userHandler.SignUp)
	e.PUT("/user/:username", userHandler.UpdateProfile, isAuth.IsAuth)
	e.PUT("/user/:username/avatar", userHandler.UpdateAvatar, isAuth.IsAuth)
	e.PUT("/user/:username/password", userHandler.ChangePassword, isAuth.IsAuth)
	// e.GET("/user/:username", userHandler.ProfileByUsername)

	e.GET("/email/dialogues", mailHander.GetDialogues, isAuth.IsAuth)
	e.POST("/email/dialogue", mailHander.CreateDialogue, isAuth.IsAuth)
	e.DELETE("/email/dialogue", mailHander.DeleteDialogue, isAuth.IsAuth)
	e.GET("/email/emails", mailHander.GetEmails, isAuth.IsAuth)
	e.DELETE("/email/emails", mailHander.DeleteEmails, isAuth.IsAuth)
	e.POST("/email", mailHander.SendEmail, isAuth.IsAuth)

	e.GET("/email/folders", mailHander.GetFolders, isAuth.IsAuth)
	e.POST("/email/folder", mailHander.CreateFolder, isAuth.IsAuth)
	e.PUT("/email/folder", mailHander.UpdateFolder, isAuth.IsAuth)
	e.DELETE("/email/folder", mailHander.DeleteFolder, isAuth.IsAuth)

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
