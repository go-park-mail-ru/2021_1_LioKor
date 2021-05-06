package sessions

import "liokor_mail/internal/pkg/common"

type SessionRepository interface {
	Create(session common.Session) error
	Get(token string) (common.Session, error)
	Delete(token string) error
}
