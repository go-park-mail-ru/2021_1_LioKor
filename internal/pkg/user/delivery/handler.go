package delivery

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2021_1_LioKor/internal/pkg/user"
	"net/http"
)

type UserHandler struct {
	UserUsecase user.UseCase
}

func (h *UserHandler) Authenticate(w http.ResponseWriter, r *http.Request) {
	var creds user.Credentials
	err :=json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.UserUsecase.Login(creds)
	if err != nil {
		switch err.(type) {
		case user.InvalidUserError:
			w.WriteHeader(http.StatusUnauthorized)
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	session, err := h.UserUsecase.CreateSession(creds.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name: "session_token",
		Value: session.Value,
		Expires: session.Expiration,
		HttpOnly: true,
	})
	http.Redirect(w, r, "/user", http.StatusFound)
}

func (h *UserHandler) UserPage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		sessionToken, err := r.Cookie("session_token")
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		sessionUser, err := h.UserUsecase.GetUserBySessionToken(sessionToken.Value)
		if err != nil {
			switch err.(type) {
			case user.InvalidSessionError, user.InvalidUserError:
				w.WriteHeader(http.StatusUnauthorized)
				return
			default:
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		jsonUser, err := json.Marshal(sessionUser)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(jsonUser)

	case http.MethodPost:
		//sign up
		var newUser user.UserSignUp
		err := json.NewDecoder(r.Body).Decode(&newUser)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = h.UserUsecase.SignUp(newUser)
		if err != nil {
			switch err.(type) {
			case user.InvalidUserError:
				w.WriteHeader(http.StatusConflict)
				return
			default:
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		w.WriteHeader(http.StatusOK)

	case http.MethodPut:
		_, err := r.Cookie("session_token")
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		sessionUser := user.User{}
		err =json.NewDecoder(r.Body).Decode(&sessionUser)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		sessionUser, err = h.UserUsecase.UpdateUser(sessionUser.Username, sessionUser)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		jsonUser, err := json.Marshal(sessionUser)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(jsonUser)
	default:
		//do smth

	}
}

