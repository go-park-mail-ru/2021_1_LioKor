package repository

import (
	"github.com/stretchr/testify/assert"
	"liokor_mail/internal/pkg/common"
	"testing"
	"time"
)

var dbConfig = common.Config{
	DBHost:           "127.0.0.1",
	DBPort:           5432,
	DBUser:           "postgres",
	DBPassword:       "12",
	DBDatabase:       "liokor_mail",
	DBConnectTimeout: 10,
}

func TestCreateSessionSuccess(t *testing.T) {
	dbInstance, err := common.NewPostgresDataBase(dbConfig)
	if err != nil {
		t.Errorf("Database error: %v\n", err)
	}
	defer dbInstance.Close()

	sessionRep := PostgresSessionsRepository{
		dbInstance,
	}

	expire, _ := time.Parse("2006-01-02T15:04:05Z", "2021-05-11T15:04:05Z")

	session := common.Session{
		UserId:       1,
		SessionToken: "sessionToken",
		Expiration:   expire,
	}

	err = sessionRep.Create(session)
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

	sessionRep := PostgresSessionsRepository{
		dbInstance,
	}

	invalidSession := common.Session{
		UserId:       1,
		SessionToken: "sessionToken",
		Expiration:   time.Now(),
	}
	err = sessionRep.Create(invalidSession)
	switch err.(type) {
	case common.InvalidSessionError:
		break
	default:
		t.Errorf("Created session for not unique session token: %v\n", err)
	}

	invalidSession = common.Session{
		UserId:       0,
		SessionToken: "uniqueSessionToken",
		Expiration:   time.Now(),
	}
	err = sessionRep.Create(invalidSession)
	switch err.(type) {
	case common.InvalidUserError:
		break
	default:
		t.Errorf("Created session for non existing user: %v\n", err)
	}
}

func TestGetSession(t *testing.T) {
	dbInstance, err := common.NewPostgresDataBase(dbConfig)
	if err != nil {
		t.Errorf("Database error: %v\n", err)
	}
	defer dbInstance.Close()

	sessionRep := PostgresSessionsRepository{
		dbInstance,
	}

	s, err := sessionRep.Get("sessionToken")
	if err != nil {
		t.Errorf("Couldn't find session: %v\n", err)
	}
	assert.Equal(t, 1, s.UserId)

	_, err = sessionRep.Get("invalidSessionToken")
	switch err.(type) {
	case common.InvalidSessionError:
		break
	default:
		t.Errorf("Found session for non existing token: %v\n", err)
	}
}

func TestDeleteSessionSuccess(t *testing.T) {
	dbInstance, err := common.NewPostgresDataBase(dbConfig)
	if err != nil {
		t.Errorf("Database error: %v\n", err)
	}
	defer dbInstance.Close()

	sessionRep := PostgresSessionsRepository{
		dbInstance,
	}

	err = sessionRep.Delete("sessionToken")
	if err != nil {
		t.Errorf("Couldn't remove session: %v\n", err)
	}
}

func TestDeleteSessionFail(t *testing.T) {
	dbInstance, err := common.NewPostgresDataBase(dbConfig)
	if err != nil {
		t.Errorf("Database error: %v\n", err)
	}
	defer dbInstance.Close()

	sessionRep := PostgresSessionsRepository{
		dbInstance,
	}

	err = sessionRep.Delete("invalidSessionToken")
	switch err.(type) {
	case common.InvalidSessionError:
		break
	default:
		t.Errorf("Deleted session for non existing token: %v\n", err)
	}
}
