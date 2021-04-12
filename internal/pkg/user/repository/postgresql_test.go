package repository

import (
	"github.com/stretchr/testify/assert"
	"liokor_mail/internal/pkg/user"
	"testing"
	"time"
)

//TODO: mock tests

var dbConfig = "host=localhost user=postgres password=12 dbname=liokor_mail_test sslmode=disable"

func TestCreateUserSuccess(t *testing.T) {
	userRep, err := NewPostgresUserRepository(dbConfig)
	if err != nil {
		t.Errorf("Database error: %v\n", err)
	}
	defer userRep.Close()

	newUser := user.User{
		Username:     "newTestUser",
		HashPassword: "hashPassword",
		AvatarURL:    "/media/",
		FullName:     "New Test User",
		ReserveEmail: "newtest@test.test",
		RegisterDate: "",
		IsAdmin:      false,
	}

	err = userRep.CreateUser(newUser)
	if err != nil {
		t.Errorf("Couldn't create user: %v\n", err)
	}
}

func TestCreateUserFail(t *testing.T) {
	userRep, err := NewPostgresUserRepository(dbConfig)
	if err != nil {
		t.Errorf("Database error: %v\n", err)
	}
	defer userRep.Close()

	newUser := user.User{
		Username:     "newTestUser",
		HashPassword: "hashPassword",
		AvatarURL:    "/media/",
		FullName:     "New Test User",
		ReserveEmail: "newtest@test.test",
		RegisterDate: "",
		IsAdmin:      false,
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
	userRep, err := NewPostgresUserRepository(dbConfig)
	if err != nil {
		t.Errorf("Database error: %v\n", err)
	}
	defer userRep.Close()

	retUser := user.User{
		Username:     "newTestUser",
		HashPassword: "hashPassword",
		AvatarURL:    "/media/",
		FullName:     "New Test User",
		ReserveEmail: "newtest@test.test",
		RegisterDate: "",
		IsAdmin:      false,
	}

	u, err := userRep.GetUserByUsername("newTestUser")
	if err != nil {
		t.Errorf("Didn't find user: %v\n", err)
	}
	assert.Equal(t, retUser, u)

	_, err = userRep.GetUserByUsername("invalidUser")
	switch err.(type) {
	case user.InvalidUserError:
		break
	default:
		t.Errorf("Found non existing user: %v\n", err)
	}
}

func TestUpdateUser(t *testing.T) {
	userRep, err := NewPostgresUserRepository(dbConfig)
	if err != nil {
		t.Errorf("Database error: %v\n", err)
	}
	defer userRep.Close()

	updUser := user.User{
		Username:     "newTestUser",
		HashPassword: "hashPassword",
		AvatarURL:    "/media/test",
		FullName:     "New Name",
		ReserveEmail: "newtest@test.test",
		RegisterDate: "",
		IsAdmin:      false,
	}

	u, err := userRep.UpdateUser(updUser.Username, updUser)
	if err != nil {
		t.Errorf("Couldn't update user: %v\n", err)
	}
	assert.Equal(t, updUser, u)

	_, err = userRep.UpdateUser("invalidUser", updUser)
	switch err.(type) {
	case user.InvalidUserError:
		break
	default:
		t.Errorf("Found non existing user: %v\n", err)
	}
}

func TestChangePassword(t *testing.T) {
	userRep, err := NewPostgresUserRepository(dbConfig)
	if err != nil {
		t.Errorf("Database error: %v\n", err)
	}
	defer userRep.Close()

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
	userRep, err := NewPostgresUserRepository(dbConfig)
	if err != nil {
		t.Errorf("Database error: %v\n", err)
	}
	defer userRep.Close()

	expire, _ := time.Parse("2006-01-02T15:04:05Z", "2021-04-11T15:04:05Z")

	session := user.Session{
		Username:     "newTestUser",
		SessionToken: "sessionToken",
		Expiration:   expire,
	}

	err = userRep.CreateSession(session)
	if err != nil {
		t.Errorf("Couldn't create session: %v\n", err)
	}
}

func TestCreateSessionFail(t *testing.T) {
	userRep, err := NewPostgresUserRepository(dbConfig)
	if err != nil {
		t.Errorf("Database error: %v\n", err)
	}
	defer userRep.Close()

	invalidSession := user.Session{
		Username:     "invalidUser",
		SessionToken: "sessionToken",
		Expiration:   time.Now(),
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
	userRep, err := NewPostgresUserRepository(dbConfig)
	if err != nil {
		t.Errorf("Database error: %v\n", err)
	}
	defer userRep.Close()

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
	userRep, err := NewPostgresUserRepository(dbConfig)
	if err != nil {
		t.Errorf("Database error: %v\n", err)
	}
	defer userRep.Close()

	err = userRep.RemoveSession("sessionToken")
	if err != nil {
		t.Errorf("Couldn't remove session: %v\n", err)
	}
}

func TestRemoveSessionFail(t *testing.T) {
	userRep, err := NewPostgresUserRepository(dbConfig)
	if err != nil {
		t.Errorf("Database error: %v\n", err)
	}
	defer userRep.Close()

	err = userRep.RemoveSession("invalidSessionToken")
	switch err.(type) {
	case user.InvalidSessionError:
		break
	default:
		t.Errorf("Deleted session for non existing token: %v\n", err)
	}
}
