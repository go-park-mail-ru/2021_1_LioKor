package usecase

import (
	"database/sql"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"liokor_mail/internal/pkg/common"
	session "liokor_mail/internal/pkg/common/protobuf_sessions"
	"liokor_mail/internal/pkg/user"
	mocks "liokor_mail/internal/pkg/user/mocks"
	sMocks "liokor_mail/internal/pkg/common/protobuf_sessions/mocks"
	"testing"
	"time"
)

var avatarBase64 = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAFIAAABaCAIAAACkHZahAAAAAXNSR0IArs4c6QAAAARnQU1BAACxjwv8YQUAAAAJcEhZcwAADsMAAA7DAcdvqGQAAAPzSURBVHhe7ZqvUzMxEEBPIpFIJBKJrEQiECgGieQ/QCJRDBKJRFZWIpFIZCUS2e9Ndyezcx2+9pLdMJ3kqV6uJXmX3fw6hlWTdO2W6Not0bVbomu3RNduia7dEl17CvP5fDabvby86HVFXKrO1D46OhqG4eDgQK9rgfPh4WF51ZnaVCzodRWSs6ClWWT+mIctdT88PGhRPBJiwu3trZZmkal9d3cn1VeL8+fnZ6kRCp0hU/vn50ebUCXOl8tliq+rqystLSC/0akdx8fH0UM63St1nZ6e8sS1tIB87RTnEBrq9/f3Ws0wvL29aWkZ+do89ZubG21OWKhb5/Pzcy0tprS5oUP6yNklvIVSbRvqX19fWuqBHbp9naFUm9YwzEjjfAe2NEu7O4NDTqZQJNW1qBgb3u7O4KC9WCykfcxkWlSGDW+XWXoTB216Iw1s9JKWFhAa3oKDNtiBrdxc/1BMeAs+2rSPntHGlpnbCNeiANz+9Mg8b1T//v5OER6U1YLnE7XmZPvn56fe2A362W6n4yIcnAOJ7mI8l3afnJzs3nQ7YwGPT2/E4J8/7+/vaWDfcSa3zhX2cxAybDw+PqrEMMzncy39BescN2ONCNGGi4sLMWGIIvK1dIM/cYYo7eVymcbk30KdTZt8AWo6Q5Q2vL6+qtMwkPBauob+Z37Se9WdIVAbUqizS9Oi1erj44NBXsphNptVdoZYbXv0xzhHCZNzKgHiv74zxGpDSmBsr6+v5TOwMiEL9EvVCdemM21IC2dnZ75HMVMJ14anpyfVXXN5efkngW0J1x4lMxDeXue+2cRqj87S9dOaks1pOVHazMzMTKq4nsBI5sVikdYwwKV+uzoh2mw57TDG7J2SmSktPQ6vNzsZ+GuP3kJvBvPmZF4fZ232jOIDuP02M6ctGt/hKWhpRTy1rTPbZhahemMDYju9VLApUA03bevMauQ/m00hna4DA8ForxKNj7Z13n07ld5aCzX/H8RBO89Z4Ld2/ONSbwRTqm2PtfO2zXZKY4Sbet6aR742hjZKS44Kss9bs8nUposYt6ShUH48knHeWkKO9miNSZ+79I89b41O8snatnH0D7mtNzxIh1DRST5N257v0uHuk221JJ+gbZ1J7KBFZZ0k31W75jm+zaP/rHBL2Em7prNALVId2a5FrmzXLl+QZEAna5Ux5zDbtdNcVc1ZILGlXnA3366tNQe/Z9+E6lKog6/5du00rup1ReLMt8vIa43y/1zPI8j8D/pwKiNzl235HmiDNSfpuNQbueyHNqAqr8Rdlm57oy14rYj3TNuLrt0SXbslunZLdO2W6Not0bVbomu3RNduiSa1V6t/FcorivktFkAAAAAASUVORK5CYII="

var config = common.Config{
	AvatarStoragePath: "../../../../media/avatars/",
}

func TestLogin(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRep := mocks.NewMockUserRepository(mockCtrl)
	mockSession := sMocks.NewMockIsAuthClient(mockCtrl)
	userUc := UserUseCase{
		Repository: mockRep,
		SessionManager: mockSession,
		Config:     config,
	}

	//Testing valid credentials
	creds := user.Credentials{
		Username: "test",
		Password: "StrongPassword1",
	}
	hashPSWD, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		t.Errorf("bcrypt error: %v\n", err)
	}
	retUser := user.User{
		Username:     "test",
		HashPassword: string(hashPSWD),
		AvatarURL:    common.NullString{sql.NullString{String: "/media/test", Valid: true}},
		FullName:     "Test test",
		ReserveEmail: "test@test.test",
		RegisterDate: "",
		IsAdmin:      false,
	}
	mockRep.EXPECT().GetUserByUsername("test").Return(retUser, nil).Times(1)
	err = userUc.Login(creds)
	if err != nil {
		t.Errorf("Didn't pass valid credentials: %v\n", err)
	}

	//Testing invalid password
	wrongPswdCreds := user.Credentials{
		Username: "test",
		Password: "password",
	}
	mockRep.EXPECT().GetUserByUsername("test").Return(retUser, nil).Times(1)
	err = userUc.Login(wrongPswdCreds)
	switch err.(type) {
	case common.InvalidUserError:
		break
	default:
		t.Errorf("Didn't pass invalid password: %v\n", err)
	}

	//Testing invalid username
	wrongUsernameCreds := user.Credentials{
		Username: "test",
		Password: "password",
	}
	mockRep.EXPECT().GetUserByUsername("test").Return(user.User{}, common.InvalidUserError{"user doesn't exist"}).Times(1)
	err = userUc.Login(wrongUsernameCreds)
	switch err.(type) {
	case common.InvalidUserError:
		break
	default:
		t.Errorf("Didn't pass invalid credentials: %v\n", err)
	}
}

func TestLogout(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRep := mocks.NewMockUserRepository(mockCtrl)
	mockSession := sMocks.NewMockIsAuthClient(mockCtrl)
	userUc := UserUseCase{
		Repository: mockRep,
		SessionManager: mockSession,
		Config:     config,
	}

	sessionToken := "sessionToken"

	mockSession.EXPECT().Delete(gomock.Any(), &session.SessionToken{ SessionToken: sessionToken }).Return(&session.Empty{}, nil).Times(1)
	err := userUc.Logout(sessionToken)
	if err != nil {
		t.Errorf("Didn't pass valid session token: %v\n", err)
	}

	mockSession.EXPECT().Delete(gomock.Any(), &session.SessionToken{ SessionToken: sessionToken }).Return(&session.Empty{}, status.Error(codes.NotFound, "Not found")).Times(1)
	err = userUc.Logout(sessionToken)
	switch err.(type) {
	case common.InvalidSessionError:
		break
	default:
		t.Errorf("Didn't pass invalid session token: %v\n", err)
	}
}

func TestCreateSession(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRep := mocks.NewMockUserRepository(mockCtrl)
	mockSession := sMocks.NewMockIsAuthClient(mockCtrl)
	userUc := UserUseCase{
		Repository: mockRep,
		SessionManager: mockSession,
		Config:     config,
	}
	sessionUser := user.User{
		Id : 1,
		Username: "test",
	}

	newSession:= session.Session{
		UserId: int32(sessionUser.Id),
		SessionToken: common.GenerateRandomString(),
		Expiration: timestamppb.New(time.Now().Add(10 * 24 * time.Hour)),
	}

	gomock.InOrder(
		mockRep.EXPECT().GetUserByUsername(sessionUser.Username).Return(sessionUser, nil).Times(1),
		mockSession.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&newSession, nil).Times(1),
	)
	s, err := userUc.CreateSession(sessionUser.Username)
	if err != nil {
		t.Errorf("Didn't create session: %v\n", err)
	}

	assert.Equal(t, s.UserId, int(newSession.UserId))
	assert.Equal(t, s.SessionToken, newSession.SessionToken)

	mockRep.EXPECT().GetUserByUsername(sessionUser.Username).Return(user.User{}, common.InvalidUserError{"user doesn't exist"}).Times(1)
	_, err = userUc.CreateSession(sessionUser.Username)
	switch err.(type) {
	case common.InvalidUserError:
		break
	default:
		t.Errorf("Didn't pass invalid user: %v\n", err)
	}
}


func TestSignUp(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRep := mocks.NewMockUserRepository(mockCtrl)
	mockSession := sMocks.NewMockIsAuthClient(mockCtrl)
	userUc := UserUseCase{
		Repository: mockRep,
		SessionManager: mockSession,
		Config:     config,
	}

	u := user.UserSignUp{
		Username:     "test",
		Password:     "StrongPassword1",
		AvatarURL:    "/media/test",
		FullName:     "Test test",
		ReserveEmail: "test@test.test",
	}

	mockRep.EXPECT().CreateUser(gomock.Any()).Return(nil).Times(1)
	err := userUc.SignUp(u)
	if err != nil {
		t.Errorf("Didn't pass valid sign up: %v\n", err)
	}

	incorrectUsername := user.UserSignUp{
		Username:     "тест",
		Password:     "StrongPassword1",
		AvatarURL:    "/media/test",
		FullName:     "Test test",
		ReserveEmail: "test@test.test",
	}
	err = userUc.SignUp(incorrectUsername)
	switch err.(type) {
	case common.InvalidUserError:
		break
	default:
		t.Errorf("Didn't pass incorrect username: %v\n", err)
	}

	incorrectPassword := user.UserSignUp{
		Username:     "test",
		Password:     "password",
		AvatarURL:    "/media/test",
		FullName:     "Test test",
		ReserveEmail: "test@test.test",
	}
	err = userUc.SignUp(incorrectPassword)
	switch err.(type) {
	case common.InvalidUserError:
		break
	default:
		t.Errorf("Didn't pass incorrect password: %v\n", err)
	}

	mockRep.EXPECT().CreateUser(gomock.Any()).Return(common.InvalidUserError{"username exists"}).Times(1)
	err = userUc.SignUp(u)
	switch err.(type) {
	case common.InvalidUserError:
		break
	default:
		t.Errorf("Didn't pass invalid user: %v\n", err)
	}
}

func TestUpdateUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRep := mocks.NewMockUserRepository(mockCtrl)
	mockSession := sMocks.NewMockIsAuthClient(mockCtrl)
	userUc := UserUseCase{
		Repository: mockRep,
		SessionManager: mockSession,
		Config:     config,
	}

	username := "test"
	newData := user.User{
		Username:     "",
		HashPassword: "",
		AvatarURL:    common.NullString{sql.NullString{String: "", Valid: false}},
		FullName:     "New Fullname",
		ReserveEmail: "newtest@test.test",
		RegisterDate: "",
		IsAdmin:      false,
	}
	retUser := user.User{
		Username:     "test",
		HashPassword: "hash",
		AvatarURL:    common.NullString{sql.NullString{String: "/media/test", Valid: true}},
		FullName:     "Test test",
		ReserveEmail: "test@test.test",
		RegisterDate: "",
		IsAdmin:      false,
	}
	updUser := user.User{
		Username:     "test",
		HashPassword: "hash",
		AvatarURL:    common.NullString{sql.NullString{String: "/media/test", Valid: true}},
		FullName:     "New Fullname",
		ReserveEmail: "newtest@test.test",
		RegisterDate: "",
		IsAdmin:      false,
	}

	gomock.InOrder(
		mockRep.EXPECT().GetUserByUsername(username).Return(retUser, nil).Times(1),
		mockRep.EXPECT().UpdateUser(username, updUser).Return(updUser, nil).Times(1),
	)
	_, err := userUc.UpdateUser(username, newData)
	if err != nil {
		t.Errorf("Didn't pass valid user: %v\n", err)
	}

	mockRep.EXPECT().GetUserByUsername(username).Return(user.User{}, common.InvalidUserError{"user doesn't exist"}).Times(1)
	_, err = userUc.UpdateUser(username, newData)
	switch err.(type) {
	case common.InvalidUserError:
		break
	default:
		t.Errorf("Didn't pass invalid username: %v\n", err)
	}

	gomock.InOrder(
		mockRep.EXPECT().GetUserByUsername(username).Return(retUser, nil).Times(1),
		mockRep.EXPECT().UpdateUser(username, updUser).Return(user.User{}, common.InvalidUserError{"username"}).Times(1),
	)
	_, err = userUc.UpdateUser(username, newData)
	switch err.(type) {
	case common.InvalidUserError:
		break
	default:
		t.Errorf("Didn't pass invalid update: %v\n", err)
	}

}

func TestUpdateAvatar(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRep := mocks.NewMockUserRepository(mockCtrl)
	mockSession := sMocks.NewMockIsAuthClient(mockCtrl)
	userUc := UserUseCase{
		Repository: mockRep,
		SessionManager: mockSession,
		Config:     config,
	}

	username := "test"
	retUser := user.User{
		Username:     "test",
		HashPassword: "hash",
		AvatarURL:    common.NullString{sql.NullString{String: "/media/test", Valid: true}},
		FullName:     "Test test",
		ReserveEmail: "test@test.test",
		RegisterDate: "",
		IsAdmin:      false,
	}
	updUser := user.User{
		Username:     "test",
		HashPassword: "hash",
		AvatarURL:    common.NullString{sql.NullString{String: "/media/randomString", Valid: true}},
		FullName:     "Test test",
		ReserveEmail: "test@test.test",
		RegisterDate: "",
		IsAdmin:      false,
	}

	gomock.InOrder(
		mockRep.EXPECT().GetUserByUsername(username).Return(retUser, nil).Times(1),
		mockRep.EXPECT().UpdateAvatar(username, gomock.Any()).Return(updUser, nil).Times(1),
	)
	_, err := userUc.UpdateAvatar(username, avatarBase64)
	if err != nil {
		t.Errorf("Didn't pass valid user: %v\n", err)
	}

	invalidAvatarURL := "invalidavatarurl"
	mockRep.EXPECT().GetUserByUsername(username).Return(retUser, nil).Times(1)
	_, err = userUc.UpdateAvatar(username, invalidAvatarURL)
	switch err.(type) {
	case common.InvalidImageError:
		break
	default:
		t.Errorf("Didn't pass invalid username: %v\n", err)
	}

}

func TestGetUserByUsername(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRep := mocks.NewMockUserRepository(mockCtrl)
	mockSession := sMocks.NewMockIsAuthClient(mockCtrl)
	userUc := UserUseCase{
		Repository: mockRep,
		SessionManager: mockSession,
		Config:     config,
	}

	username := "test"
	retUser := user.User{
		Username:     "test",
		HashPassword: "hash",
		AvatarURL:    common.NullString{sql.NullString{String: "/media/test", Valid: true}},
		FullName:     "Test test",
		ReserveEmail: "test@test.test",
		RegisterDate: "",
		IsAdmin:      false,
	}

	mockRep.EXPECT().GetUserByUsername(username).Return(retUser, nil).Times(1)
	u, err := userUc.GetUserByUsername(username)
	if err != nil || u != retUser {
		t.Errorf("Didn't pass valid user: %v\n", err)
	}

	mockRep.EXPECT().GetUserByUsername(username).Return(user.User{}, common.InvalidUserError{"user doesn't exist"}).Times(1)
	_, err = userUc.GetUserByUsername(username)
	switch err.(type) {
	case common.InvalidUserError:
		break
	default:
		t.Errorf("Didn't pass invalid user: %v\n", err)
	}
}

func TestGetUserById(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRep := mocks.NewMockUserRepository(mockCtrl)
	mockSession := sMocks.NewMockIsAuthClient(mockCtrl)
	userUc := UserUseCase{
		Repository: mockRep,
		SessionManager: mockSession,
		Config:     config,
	}

	userId := 1
	retUser := user.User{
		Id : 1,
		Username:     "test",
		HashPassword: "hash",
		AvatarURL:    common.NullString{sql.NullString{String: "/media/test", Valid: true}},
		FullName:     "Test test",
		ReserveEmail: "test@test.test",
		RegisterDate: "",
		IsAdmin:      false,
	}

	mockRep.EXPECT().GetUserById(userId).Return(retUser, nil).Times(1)
	u, err := userUc.GetUserById(userId)
	if err != nil || u != retUser {
		t.Errorf("Didn't pass valid user: %v\n", err)
	}

	mockRep.EXPECT().GetUserById(userId).Return(user.User{}, common.InvalidUserError{"user doesn't exist"}).Times(1)
	_, err = userUc.GetUserById(userId)
	switch err.(type) {
	case common.InvalidUserError:
		break
	default:
		t.Errorf("Didn't pass invalid user: %v\n", err)
	}
}
func TestChangePassword(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRep := mocks.NewMockUserRepository(mockCtrl)
	mockSession := sMocks.NewMockIsAuthClient(mockCtrl)
	userUc := UserUseCase{
		Repository: mockRep,
		SessionManager: mockSession,
		Config:     config,
	}

	hashPSWD, err := bcrypt.GenerateFromPassword([]byte("StrongPassword1"), bcrypt.DefaultCost)
	if err != nil {
		t.Errorf("bcrypt error: %v\n", err)
	}
	sessionUser := user.User{
		Username:     "test",
		HashPassword: string(hashPSWD),
		AvatarURL:    common.NullString{sql.NullString{String: "/media/test", Valid: true}},
		FullName:     "Test test",
		ReserveEmail: "test@test.test",
		RegisterDate: "",
		IsAdmin:      false,
	}
	newPSWD := user.ChangePassword{
		OldPassword: "StrongPassword1",
		NewPassword: "NewStrongPassword2",
	}

	mockRep.EXPECT().ChangePassword(sessionUser.Username, gomock.Any()).Return(nil).Times(1)
	err = userUc.ChangePassword(sessionUser, newPSWD)
	if err != nil {
		t.Errorf("Didn't pass valid change: %v\n", err)
	}

	mockRep.EXPECT().ChangePassword(sessionUser.Username, gomock.Any()).Return(common.InvalidUserError{"user doesn't exist"}).Times(1)
	err = userUc.ChangePassword(sessionUser, newPSWD)
	switch err.(type) {
	case common.InvalidUserError:
		break
	default:
		t.Errorf("Didn't pass invalid user: %v\n", err)
	}

	newInvalidPSWD := user.ChangePassword{
		OldPassword: "StrongPassword1",
		NewPassword: "New",
	}
	err = userUc.ChangePassword(sessionUser, newInvalidPSWD)
	switch err.(type) {
	case common.InvalidUserError:
		break
	default:
		t.Errorf("Didn't pass invalid new password: %v\n", err)
	}

	newPSWDInvalidOld := user.ChangePassword{
		OldPassword: "StrongPasswo1",
		NewPassword: "NewStrongPassword2",
	}
	err = userUc.ChangePassword(sessionUser, newPSWDInvalidOld)
	switch err.(type) {
	case common.InvalidUserError:
		break
	default:
		t.Errorf("Didn't pass invalid old password: %v\n", err)
	}

}
