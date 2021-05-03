package sessions

import (
	"context"
	session "liokor_mail/internal/pkg/common/protobuf_sessions"
)

type SessionsDelivery interface {
	Create(ctx context.Context, s *session.Session) (*session.Session, error)
	Get(ctx context.Context, token *session.SessionToken) (*session.Session, error)
	Delete(ctx context.Context, token *session.SessionToken) (*session.Empty, error)
}
