package usecase

import (
	"liokor_mail/internal/pkg/common"
	"liokor_mail/internal/pkg/sessions"
	"time"
	"log"
)

type SessionUsecase struct {
	SessionRepository sessions.SessionRepository
}

func (uc *SessionUsecase) Create(session common.Session) (common.Session, error) {
	err := uc.SessionRepository.Create(session)

	if err != nil {
		return common.Session{}, err
	}

	return session, nil

}
func (uc *SessionUsecase) Get(token string) (common.Session, error) {
	session, err := uc.SessionRepository.Get(token)
	if err != nil {
		return common.Session{}, err
	}

	if session.Expiration.Before(time.Now()) {
		err := uc.SessionRepository.Delete(token)
		if err != nil {
			log.Println("WARN: Unable to delete session token")
		}
		return common.Session{}, common.InvalidSessionError{"session token expired"}
	}
	return session, nil
}

func (uc *SessionUsecase) Delete(token string) error {
	return uc.SessionRepository.Delete(token)
}
