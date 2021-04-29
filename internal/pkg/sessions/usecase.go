package sessions

import "liokor_mail/internal/pkg/common"

type SessionsUseCase interface {
	Create(session common.Session) (common.Session, error)
	Get(token string) (common.Session, error)
	Delete(token string) error
}