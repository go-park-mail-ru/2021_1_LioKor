package repository

import (
	"github.com/stretchr/testify/assert"
	"liokor_mail/internal/pkg/common"
	"liokor_mail/internal/pkg/user"
	"testing"
	"time"
)


var dbConfig = "host=localhost user=postgres password=12 dbname=liokor_mail_test sslmode=disable"

func TestCreateUserSuccess(t *testing.T) {
	dbInstance, err := common.NewPostgresDataBase(dbConfig)
	if err != nil {
		t.Errorf("Database error: %v\n", err)
	}
	defer dbInstance.Close()

	userRep := PostgresUserRepository{
		dbInstance,
	}

	newUser := user.User{
		Username: "newTestUser",
		HashPassword: "hashPassword",
		AvatarURL: "/media/test",
		FullName: "New Test User",
		ReserveEmail: "newtest@test.test",
		RegisterDate: "",
		IsAdmin: false,
	}

	err = userRep.CreateUser(newUser)
	if err != nil {
		t.Errorf("Couldn't create user: %v\n", err)
	}
}

func TestCreateUserFail(t *testing.T) {
	dbInstance, err := common.NewPostgresDataBase(dbConfig)
	if err != nil {
		t.Errorf("Database error: %v\n", err)
	}
	defer dbInstance.Close()

	userRep := PostgresUserRepository{
		dbInstance,
	}

	newUser := user.User{
		Username: "newTestUser",
		HashPassword: "hashPassword",
		AvatarURL: "/media/",
		FullName: "New Test User",
		ReserveEmail: "newtest@test.test",
		RegisterDate: "",
		IsAdmin: false,
	}

	err = userRep.CreateUser(newUser)
	switch err.(type) {
	case user.InvalidUserError:
		break
	default:
		t.Errorf("Created user with existing username: %v\n", err)
	}
}

func TestGetUserByUsername(t *testing.T) {
	dbInstance, err := common.NewPostgresDataBase(dbConfig)
	if err != nil {
		t.Errorf("Database error: %v\n", err)
	}
	defer dbInstance.Close()

	userRep := PostgresUserRepository{
		dbInstance,
	}

	retUser := user.User{
		Username: "newTestUser",
		HashPassword: "hashPassword",
		AvatarURL: "/media/test",
		FullName: "New Test User",
		ReserveEmail: "newtest@test.test",
		RegisterDate: "",
		IsAdmin: false,
	}

	u, err := userRep.GetUserByUsername("newTestUser")
	if err != nil {
		t.Errorf("Didn't find user: %v\n", err)
	}
	assert.Equal(t, retUser.Username, u.Username)

	_, err = userRep.GetUserByUsername("invalidUser")
	switch err.(type) {
	case user.InvalidUserError:
		break
	default:
		t.Errorf("Found non existing user: %v\n", err)
	}
}

func TestUpdateUser(t *testing.T) {
	dbInstance, err := common.NewPostgresDataBase(dbConfig)
	if err != nil {
		t.Errorf("Database error: %v\n", err)
	}
	defer dbInstance.Close()

	userRep := PostgresUserRepository{
		dbInstance,
	}

	updUser := user.User{
		Username: "newTestUser",
		HashPassword: "hashPassword",
		AvatarURL: "/media/test",
		FullName: "New Name",
		ReserveEmail: "newtest@test.test",
		RegisterDate: "",
		IsAdmin: false,
	}

	u, err := userRep.UpdateUser(updUser.Username, updUser)
	if err != nil {
		t.Errorf("Couldn't update user: %v\n", err)
	}
	assert.Equal(t, updUser.Username, u.Username)

	_, err = userRep.UpdateUser("invalidUser", updUser)
	switch err.(type) {
	case user.InvalidUserError:
		break
	default:
		t.Errorf("Found non existing user: %v\n", err)
	}
}

func TestChangePassword(t *testing.T) {
	dbInstance, err := common.NewPostgresDataBase(dbConfig)
	if err != nil {
		t.Errorf("Database error: %v\n", err)
	}
	defer dbInstance.Close()

	userRep := PostgresUserRepository{
		dbInstance,
	}

	err = userRep.ChangePassword("newTestUser", "newHashPassword")
	if err != nil {
		t.Errorf("Couldn't change password: %v\n", err)
	}

	err = userRep.ChangePassword("invalidUser", "newHashPassword")
	switch err.(type) {
	case user.InvalidUserError:
		break
	default:
		t.Errorf("Changed non existing user: %v\n", err)
	}
}

func TestCreateSessionSuccess(t *testing.T) {
	dbInstance, err := common.NewPostgresDataBase(dbConfig)
	if err != nil {
		t.Errorf("Database error: %v\n", err)
	}
	defer dbInstance.Close()

	userRep := PostgresUserRepository{
		dbInstance,
	}

	expire, _ := time.Parse("2006-01-02T15:04:05Z", "2021-04-11T15:04:05Z")

	session := user.Session{
		Username: "newTestUser",
		SessionToken: "sessionToken",
		Expiration: expire,
	}

	err = userRep.CreateSession(session)
	if err != nil {
		t.Errorf("Couldn't create session: %v\n", err)
	}
}

func TestCreateSessionFail(t *testing.T) {
	dbInstance, err := common.NewPostgresDataBase(dbConfig)
	if err != nil {
		t.Errorf("Database error: %v\n", err)
	}
	defer dbInstance.Close()

	userRep := PostgresUserRepository{
		dbInstance,
	}

	invalidSession := user.Session{
		Username: "invalidUser",
		SessionToken: "sessionToken",
		Expiration: time.Now(),
	}
	err = userRep.CreateSession(invalidSession)
	switch err.(type) {
	case user.InvalidUserError:
		break
	default:
		t.Errorf("Created session for non existing user: %v\n", err)
	}
}

func TestGetSessionBySessionToken(t *testing.T) {
	dbInstance, err := common.NewPostgresDataBase(dbConfig)
	if err != nil {
		t.Errorf("Database error: %v\n", err)
	}
	defer dbInstance.Close()

	userRep := PostgresUserRepository{
		dbInstance,
	}

	s, err := userRep.GetSessionBySessionToken("sessionToken")
	if err != nil {
		t.Errorf("Couldn't find session: %v\n", err)
	}
	assert.Equal(t, "newTestUser", s.Username)

	_, err = userRep.GetSessionBySessionToken("invalidSessionToken")
	switch err.(type) {
	case user.InvalidSessionError:
		break
	default:
		t.Errorf("Found session for non existing token: %v\n", err)
	}
}

func TestRemoveSessionSuccess(t *testing.T) {
	dbInstance, err := common.NewPostgresDataBase(dbConfig)
	if err != nil {
		t.Errorf("Database error: %v\n", err)
	}
	defer dbInstance.Close()

	userRep := PostgresUserRepository{
		dbInstance,
	}

	err = userRep.RemoveSession("sessionToken")
	if err != nil {
		t.Errorf("Couldn't remove session: %v\n", err)
	}
}

func TestRemoveSessionFail(t *testing.T) {
	dbInstance, err := common.NewPostgresDataBase(dbConfig)
	if err != nil {
		t.Errorf("Database error: %v\n", err)
	}
	defer dbInstance.Close()

	userRep := PostgresUserRepository{
		dbInstance,
	}

	err = userRep.RemoveSession("invalidSessionToken")
	switch err.(type) {
	case user.InvalidSessionError:
		break
	default:
		t.Errorf("Deleted session for non existing token: %v\n", err)
	}
}

func TestRemoveUserSuccess(t *testing.T) {
	dbInstance, err := common.NewPostgresDataBase(dbConfig)
	if err != nil {
		t.Errorf("Database error: %v\n", err)
	}
	defer dbInstance.Close()

	userRep := PostgresUserRepository{
		dbInstance,
	}

	err = userRep.RemoveUser("newTestUser")
	if err != nil {
		t.Errorf("Couldn't remove user: %v\n", err)
	}
}