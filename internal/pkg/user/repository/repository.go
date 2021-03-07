package repository

import (
	"liokor_mail/internal/pkg/user"
)

type UserRepository struct {
	UserDB    map[string]user.User
	SessionDB map[string]user.Session
}

func (ur *UserRepository) CreateSession(session user.Session) error {
	if _, exists := ur.SessionDB[session.SessionToken]; exists {
		return user.InvalidSessionError{"session token exists"}
	}
	ur.SessionDB[session.SessionToken] = session
	return nil
}

func (ur *UserRepository) GetSessionBySessionToken(token string) (user.Session, error) {
	if session, exists := ur.SessionDB[token]; exists {
		return session, nil
	}
	return user.Session{}, user.InvalidSessionError{"session doesn't exist"}
}

func (ur *UserRepository) GetUserByUsername(username string) (user.User, error) {
	if user, exists := ur.UserDB[username]; exists {
		return user, nil
	}
	return user.User{}, user.InvalidUserError{"user doesn't exist"}
}

func (ur *UserRepository) CreateUser(newUser user.User) error {
	if _, exists := ur.UserDB[newUser.Username]; exists {
		return user.InvalidUserError{"username taken"}
	}
	ur.UserDB[newUser.Username] = newUser
	return nil
}

func (ur *UserRepository) UpdateUser(username string, newData user.User) (user.User, error) {
	if _, exists := ur.UserDB[username]; !exists {
		return user.User{}, user.InvalidUserError{"user doesn't exist"}
	}
	ur.UserDB[username] = newData
	return newData, nil
}

func (ur *UserRepository) ChangePassword(username string, newPSWD string) error {
	if _, exists := ur.UserDB[username]; !exists {
		return user.InvalidUserError{"user doesn't exist"}
	}
	data := ur.UserDB[username]
	data.HashPassword = newPSWD
	ur.UserDB[username] = data
	return nil
}

func (ur *UserRepository) RemoveSession(token string) error {
	if _, exists := ur.SessionDB[token]; !exists {
		return user.InvalidSessionError{"session doesn't exist"}
	}
	delete(ur.SessionDB, token)
	return nil
}
