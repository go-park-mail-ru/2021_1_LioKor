package usecase

import (
	"github.com/golang/mock/gomock"
	"liokor_mail/internal/pkg/common"
	"liokor_mail/internal/pkg/sessions/mocks"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	sessionRep := mocks.NewMockSessionRepository(mockCtrl)
	sessionUC := SessionUsecase{
		SessionRepository: sessionRep,
	}

	s := common.Session{
		UserId: 1,
		SessionToken: "sessionToken",
		Expiration: time.Now().Add(10 * 24 * time.Hour),
	}

	sessionRep.EXPECT().Create(s).Return(nil).Times(1)
	_, err := sessionUC.Create(s)
	if err != nil {
		t.Errorf("Didn't create valid session: %v\n", err.Error())
	}

	sessionRep.EXPECT().Create(s).Return(common.InvalidSessionError{"Invalid session"}).Times(1)
	_, err = sessionUC.Create(s)
	switch err.(type) {
	case common.InvalidSessionError:
		break
	default:
		t.Errorf("Didn't pass invalid session: %v\n", err)
	}
}

func TestGet(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	sessionRep := mocks.NewMockSessionRepository(mockCtrl)
	sessionUC := SessionUsecase{
		SessionRepository: sessionRep,
	}

	s := common.Session{
		UserId: 1,
		SessionToken: "sessionToken",
		Expiration: time.Now().Add(10 * 24 * time.Hour),
	}

	sessionRep.EXPECT().Get(s.SessionToken).Return(s, nil).Times(1)
	_, err := sessionUC.Get(s.SessionToken)
	if err != nil {
		t.Errorf("Didn't create valid session: %v\n", err.Error())
	}

	s.Expiration = time.Now().Add(-10 * 24 * time.Hour)
	gomock.InOrder(
		sessionRep.EXPECT().Get(s.SessionToken).Return(s, nil).Times(1),
		sessionRep.EXPECT().Delete(s.SessionToken).Return(nil).Times(1),
	)
	_, err = sessionUC.Get(s.SessionToken)
	switch err.(type) {
	case common.InvalidSessionError:
		break
	default:
		t.Errorf("Didn't pass expired token: %v\n", err)
	}
}

func TestDelete(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	sessionRep := mocks.NewMockSessionRepository(mockCtrl)
	sessionUC := SessionUsecase{
		SessionRepository: sessionRep,
	}

	sessionToken := "sessionToken"
	sessionRep.EXPECT().Delete(sessionToken).Return(nil).Times(1)
	err := sessionUC.Delete(sessionToken)
	if err != nil {
		t.Errorf("Didn't delete valid token: %v\n", err.Error())
	}
}