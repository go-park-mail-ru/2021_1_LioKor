package repository

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"liokor_mail/internal/pkg/common"
	"liokor_mail/internal/pkg/user"
	"testing"
)

var dbConfig = common.Config{
	DBHost:           "127.0.0.1",
	DBPort:           5432,
	DBUser:           "postgres",
	DBPassword:       "12",
	DBDatabase:       "liokor_mail",
	DBConnectTimeout: 10,
}

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
		Username:     "newTestUser",
		HashPassword: "hashPassword",
		AvatarURL:    common.NullString{sql.NullString{String: "/media/test", Valid: true}},
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
	dbInstance, err := common.NewPostgresDataBase(dbConfig)
	if err != nil {
		t.Errorf("Database error: %v\n", err)
	}
	defer dbInstance.Close()

	userRep := PostgresUserRepository{
		dbInstance,
	}

	newUser := user.User{
		Username:     "newTestUser",
		HashPassword: "hashPassword",
		AvatarURL:    common.NullString{sql.NullString{String: "/media", Valid: true}},
		FullName:     "New Test User",
		ReserveEmail: "newtest@test.test",
		RegisterDate: "",
		IsAdmin:      false,
	}

	err = userRep.CreateUser(newUser)
	switch err.(type) {
	case common.InvalidUserError:
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
		Username:     "newTestUser",
		HashPassword: "hashPassword",
		AvatarURL:    common.NullString{sql.NullString{String: "/media/test", Valid: true}},
		FullName:     "New Test User",
		ReserveEmail: "newtest@test.test",
		RegisterDate: "",
		IsAdmin:      false,
	}

	u, err := userRep.GetUserByUsername("newTestUser")
	if err != nil {
		t.Errorf("Didn't find user: %v\n", err)
	}
	assert.Equal(t, retUser.Username, u.Username)

	_, err = userRep.GetUserByUsername("invalidUser")
	switch err.(type) {
	case common.InvalidUserError:
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
		Username:     "newTestUser",
		HashPassword: "hashPassword",
		AvatarURL:    common.NullString{sql.NullString{String: "/media/test", Valid: true}},
		FullName:     "New Name",
		ReserveEmail: "newtest@test.test",
		RegisterDate: "",
		IsAdmin:      false,
	}

	u, err := userRep.UpdateUser(updUser.Username, updUser)
	if err != nil {
		t.Errorf("Couldn't update user: %v\n", err)
	}
	assert.Equal(t, updUser.FullName, u.FullName)

	_, err = userRep.UpdateUser("invalidUser", updUser)
	switch err.(type) {
	case common.InvalidUserError:
		break
	default:
		t.Errorf("Found non existing user: %v\n", err)
	}
}

func TestUpgradeAvatar(t *testing.T) {
	dbInstance, err := common.NewPostgresDataBase(dbConfig)
	if err != nil {
		t.Errorf("Database error: %v\n", err)
	}
	defer dbInstance.Close()

	userRep := PostgresUserRepository{
		dbInstance,
	}

	updUser := user.User{
		Username:     "newTestUser",
		HashPassword: "hashPassword",
		AvatarURL:    common.NullString{sql.NullString{String: "/media/newTest", Valid: true}},
		FullName:     "New Name",
		ReserveEmail: "newtest@test.test",
		RegisterDate: "",
		IsAdmin:      false,
	}

	u, err := userRep.UpdateAvatar(updUser.Username, common.NullString{sql.NullString{String: "/media/newTest", Valid: true}})
	if err != nil {
		t.Errorf("Couldn't update avatar: %v\n", err)
	}
	assert.Equal(t, updUser.AvatarURL, u.AvatarURL)
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
	case common.InvalidUserError:
		break
	default:
		t.Errorf("Changed non existing user: %v\n", err)
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
