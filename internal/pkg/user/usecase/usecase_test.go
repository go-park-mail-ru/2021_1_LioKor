package usecase

import (
	"github.com/golang/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"liokor_mail/internal/pkg/common"
	"liokor_mail/internal/pkg/user"
	"liokor_mail/internal/pkg/user/mocks"
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
	userUc := UserUseCase{
		Repository: mockRep,
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
		AvatarURL:    "/media/test",
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
	case user.InvalidUserError:
		break
	default:
		t.Errorf("Didn't pass invalid password: %v\n", err)
	}

	//Testing invalid username
	wrongUsernameCreds := user.Credentials{
		Username: "test",
		Password: "password",
	}
	mockRep.EXPECT().GetUserByUsername("test").Return(user.User{}, user.InvalidUserError{"user doesn't exist"}).Times(1)
	err = userUc.Login(wrongUsernameCreds)
	switch err.(type) {
	case user.InvalidUserError:
		break
	default:
		t.Errorf("Didn't pass invalid credentials: %v\n", err)
	}
}

func TestLogout(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRep := mocks.NewMockUserRepository(mockCtrl)
	userUc := UserUseCase{
		Repository: mockRep,
		Config:     config,
	}

	sessionToken := "sessionToken"

	mockRep.EXPECT().RemoveSession(sessionToken).Return(nil).Times(1)
	err := userUc.Logout(sessionToken)
	if err != nil {
		t.Errorf("Didn't pass valid session token: %v\n", err)
	}

	mockRep.EXPECT().RemoveSession(sessionToken).Return(user.InvalidSessionError{"session doesn't exist"}).Times(1)
	err = userUc.Logout(sessionToken)
	switch err.(type) {
	case user.InvalidSessionError:
		break
	default:
		t.Errorf("Didn't pass invalid session token: %v\n", err)
	}
}

func TestCreateSession(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRep := mocks.NewMockUserRepository(mockCtrl)
	userUc := UserUseCase{
		Repository: mockRep,
		Config:     config,
	}
	username := "test"

	mockRep.EXPECT().CreateSession(gomock.Any()).Return(nil).Times(1)
	sessionToken, err := userUc.CreateSession(username)
	if err != nil || sessionToken.Value == "" {
		t.Errorf("Didn't create session: %v\n", err)
	}

	mockRep.EXPECT().CreateSession(gomock.Any()).Return(user.InvalidUserError{"user doesn't exist"}).Times(1)
	_, err = userUc.CreateSession(username)
	switch err.(type) {
	case user.InvalidUserError:
		break
	default:
		t.Errorf("Didn't pass invalid user: %v\n", err)
	}
}

func TestGetUserBySessionToken(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRep := mocks.NewMockUserRepository(mockCtrl)
	userUc := UserUseCase{
		Repository: mockRep,
		Config:     config,
	}

	sessionToken := "sessionToken"
	retSession := user.Session{
		Username:     "test",
		SessionToken: sessionToken,
		Expiration:   time.Now().Add(10 * 24 * time.Hour),
	}
	retUser := user.User{
		Username:     "test",
		HashPassword: "hash",
		AvatarURL:    "/media/test",
		FullName:     "Test test",
		ReserveEmail: "test@test.test",
		RegisterDate: "",
		IsAdmin:      false,
	}

	gomock.InOrder(
		mockRep.EXPECT().GetSessionBySessionToken(sessionToken).Return(retSession, nil).Times(1),
		mockRep.EXPECT().GetUserByUsername(retSession.Username).Return(retUser, nil).Times(1),
	)

	u, err := userUc.GetUserBySessionToken(sessionToken)
	if err != nil || u != retUser {
		t.Errorf("Didn't pass valid session token: %v\n", err)
	}

	mockRep.EXPECT().GetSessionBySessionToken(sessionToken).Return(user.Session{}, user.InvalidSessionError{"session doesn't exist"}).Times(1)
	_, err = userUc.GetUserBySessionToken(sessionToken)
	switch err.(type) {
	case user.InvalidSessionError:
		break
	default:
		t.Errorf("Didn't pass invalid session token: %v\n", err)
	}

	expiredSession := user.Session{
		Username:     "test",
		SessionToken: sessionToken,
		Expiration:   time.Now().AddDate(0, 0, -1),
	}
	mockRep.EXPECT().GetSessionBySessionToken(sessionToken).Return(expiredSession, nil).Times(1)
	_, err = userUc.GetUserBySessionToken(sessionToken)
	switch err.(type) {
	case user.InvalidSessionError:
		break
	default:
		t.Errorf("Didn't pass expired token: %v\n", err)
	}

	gomock.InOrder(
		mockRep.EXPECT().GetSessionBySessionToken(sessionToken).Return(retSession, nil).Times(1),
		mockRep.EXPECT().GetUserByUsername(retSession.Username).Return(user.User{}, user.InvalidUserError{"user doesn't exist"}).Times(1),
	)
	_, err = userUc.GetUserBySessionToken(sessionToken)
	switch err.(type) {
	case user.InvalidUserError:
		break
	default:
		t.Errorf("Didn't pass invalid user: %v\n", err)
	}
}

func TestSignUp(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRep := mocks.NewMockUserRepository(mockCtrl)
	userUc := UserUseCase{
		Repository: mockRep,
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
	case user.InvalidUserError:
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
	case user.InvalidUserError:
		break
	default:
		t.Errorf("Didn't pass incorrect password: %v\n", err)
	}

	mockRep.EXPECT().CreateUser(gomock.Any()).Return(user.InvalidUserError{"username exists"}).Times(1)
	err = userUc.SignUp(u)
	switch err.(type) {
	case user.InvalidUserError:
		break
	default:
		t.Errorf("Didn't pass invalid user: %v\n", err)
	}
}

func TestUpdateUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRep := mocks.NewMockUserRepository(mockCtrl)
	userUc := UserUseCase{
		Repository: mockRep,
		Config:     config,
	}

	username := "test"
	newData := user.User{
		Username:     "",
		HashPassword: "",
		AvatarURL:    avatarBase64,
		FullName:     "New Fullname",
		ReserveEmail: "newtest@test.test",
		RegisterDate: "",
		IsAdmin:      false,
	}
	retUser := user.User{
		Username:     "test",
		HashPassword: "hash",
		AvatarURL:    "/media/test",
		FullName:     "Test test",
		ReserveEmail: "test@test.test",
		RegisterDate: "",
		IsAdmin:      false,
	}
	updUser := user.User{
		Username:     "test",
		HashPassword: "hash",
		AvatarURL:    "/media/someRandomString",
		FullName:     "New Fullname",
		ReserveEmail: "newtest@test.test",
		RegisterDate: "",
		IsAdmin:      false,
	}

	gomock.InOrder(
		mockRep.EXPECT().GetUserByUsername(username).Return(retUser, nil).Times(1),
		mockRep.EXPECT().UpdateUser(username, gomock.Any()).Return(updUser, nil).Times(1),
	)
	_, err := userUc.UpdateUser(username, newData)
	if err != nil {
		t.Errorf("Didn't pass valid user: %v\n", err)
	}

	invalidNewData := user.User{
		Username:     "",
		HashPassword: "",
		AvatarURL:    "invalidImage",
		FullName:     "New Fullname",
		ReserveEmail: "newtest@test.test",
		RegisterDate: "",
		IsAdmin:      false,
	}

	mockRep.EXPECT().GetUserByUsername(username).Return(retUser, nil).Times(1)
	_, err = userUc.UpdateUser(username, invalidNewData)
	switch err.(type) {
	case user.InvalidImageError:
		break
	default:
		t.Errorf("Didn't pass invalid image: %v\n", err)
	}

	mockRep.EXPECT().GetUserByUsername(username).Return(user.User{}, user.InvalidUserError{"user doesn't exist"}).Times(1)
	_, err = userUc.UpdateUser(username, newData)
	switch err.(type) {
	case user.InvalidUserError:
		break
	default:
		t.Errorf("Didn't pass invalid username: %v\n", err)
	}

	gomock.InOrder(
		mockRep.EXPECT().GetUserByUsername(username).Return(retUser, nil).Times(1),
		mockRep.EXPECT().UpdateUser(username, gomock.Any()).Return(user.User{}, user.InvalidUserError{"username"}).Times(1),
	)
	_, err = userUc.UpdateUser(username, newData)
	switch err.(type) {
	case user.InvalidUserError:
		break
	default:
		t.Errorf("Didn't pass invalid update: %v\n", err)
	}

}

func TestGetUserByUsername(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRep := mocks.NewMockUserRepository(mockCtrl)
	userUc := UserUseCase{
		Repository: mockRep,
		Config:     config,
	}
	username := "test"
	retUser := user.User{
		Username:     "test",
		HashPassword: "hash",
		AvatarURL:    "/media/test",
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

	mockRep.EXPECT().GetUserByUsername(username).Return(user.User{}, user.InvalidUserError{"user doesn't exist"}).Times(1)
	_, err = userUc.GetUserByUsername(username)
	switch err.(type) {
	case user.InvalidUserError:
		break
	default:
		t.Errorf("Didn't pass invalid user: %v\n", err)
	}
}

func TestChangePassword(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRep := mocks.NewMockUserRepository(mockCtrl)
	userUc := UserUseCase{
		Repository: mockRep,
		Config:     config,
	}

	hashPSWD, err := bcrypt.GenerateFromPassword([]byte("StrongPassword1"), bcrypt.DefaultCost)
	if err != nil {
		t.Errorf("bcrypt error: %v\n", err)
	}
	sessionUser := user.User{
		Username:     "test",
		HashPassword: string(hashPSWD),
		AvatarURL:    "/media/test",
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

	mockRep.EXPECT().ChangePassword(sessionUser.Username, gomock.Any()).Return(user.InvalidUserError{"user doesn't exist"}).Times(1)
	err = userUc.ChangePassword(sessionUser, newPSWD)
	switch err.(type) {
	case user.InvalidUserError:
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
	case user.InvalidUserError:
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
	case user.InvalidUserError:
		break
	default:
		t.Errorf("Didn't pass invalid old password: %v\n", err)
	}

}
