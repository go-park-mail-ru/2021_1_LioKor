package authServer

import (
	"google.golang.org/grpc"
	"liokor_mail/internal/pkg/common"
	session "liokor_mail/internal/pkg/common/protobuf_sessions"
	"os"
	"time"

	sessionDelivery "liokor_mail/internal/pkg/sessions/delivery"
	sessionRepository "liokor_mail/internal/pkg/sessions/repository"
	sessionUsecase "liokor_mail/internal/pkg/sessions/usecase"

	"context"
	"fmt"
	"log"
	"net"
)

func StartAuthServer(config common.Config, quit chan os.Signal) {
	dbInstance, err := common.NewPostgresDataBase(config)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbInstance.Close()

	if config.Debug {
		log.Println("WARN: RUNNING IN THE DEBUG MODE! DON'T USE IN PRODUCTION!")
	}

	sessionRep := &sessionRepository.PostgresSessionsRepository{dbInstance}
	sessionUC := &sessionUsecase.SessionUsecase{sessionRep}
	sessionDel := &sessionDelivery.SessionsDelivery{sessionUC}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.AuthPort))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()

	session.RegisterIsAuthServer(server, sessionDel)

	go func() {
		log.Printf("starting server at %v\n", config.AuthPort)
		err = server.Serve(lis)
		if err != nil {
			log.Fatal(err)
		}
	}()

	<-quit
	log.Println("Interrupt signal received. Shutting down server...")
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	server.GracefulStop()
}
