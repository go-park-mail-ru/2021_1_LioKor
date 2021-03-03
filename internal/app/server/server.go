package server

import (
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2021_1_LioKor/internal/pkg/user"
	"github.com/go-park-mail-ru/2021_1_LioKor/internal/pkg/user/delivery"
	"github.com/go-park-mail-ru/2021_1_LioKor/internal/pkg/user/repository"
	"github.com/go-park-mail-ru/2021_1_LioKor/internal/pkg/user/usecase"
)


func StartServer() {
	userHandler := delivery.UserHandler{
		usecase.UserUseCase{
			repository.UserRepository{
				map[string]user.User{},
				map[string]user.Session{},
			},
		},
	}
	http.HandleFunc("/user/auth", userHandler.Authenticate)
	http.HandleFunc("/user", userHandler.UserPage)

	log.Fatal(http.ListenAndServe(":8000", nil))
}