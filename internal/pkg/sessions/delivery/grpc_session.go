package delivery

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"liokor_mail/internal/pkg/common"
	session "liokor_mail/internal/pkg/common/protobuf_sessions"
	"liokor_mail/internal/pkg/sessions"
)

type SessionsDelivery struct {
	SessionsUseCase sessions.SessionsUseCase
}

func (sd *SessionsDelivery) Create(ctx context.Context, s *session.Session) (*session.Session, error) {
	newSession, err := sd.SessionsUseCase.Create(common.Session{
		int(s.UserId),
		s.SessionToken,
		s.Expiration.AsTime(),
	},
	)
	if err != nil {
		switch err.(type) {
		case common.InvalidUserError:
			return nil, status.Error(codes.NotFound, ("пользователя не существует"))
		case common.InvalidSessionError:
			return nil, status.Error(codes.NotFound, "сессия уже существует")
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	return &session.Session{
		UserId:       int32(newSession.UserId),
		SessionToken: newSession.SessionToken,
		Expiration:   timestamppb.New(newSession.Expiration),
	}, nil
}

func (sd *SessionsDelivery) Get(ctx context.Context, token *session.SessionToken) (*session.Session, error) {
	s, err := sd.SessionsUseCase.Get(token.SessionToken)
	if err != nil {
		switch err.(type) {
		case common.InvalidSessionError:
			return nil, status.Error(codes.NotFound, "такой сессии не существует")
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	return &session.Session{
		UserId:       int32(s.UserId),
		SessionToken: s.SessionToken,
		Expiration:   timestamppb.New(s.Expiration),
	}, nil
}

func (sd *SessionsDelivery) Delete(ctx context.Context, token *session.SessionToken) (*session.Empty, error) {
	err := sd.SessionsUseCase.Delete(token.SessionToken)
	if err != nil {
		switch err.(type) {
		case common.InvalidSessionError:
			return nil, status.Error(codes.NotFound, "такой сессии не существует")
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	return &session.Empty{}, nil
}
