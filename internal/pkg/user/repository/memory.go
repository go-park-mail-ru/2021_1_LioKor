package repository

import (
	"liokor_mail/internal/pkg/common"
	"liokor_mail/internal/pkg/user"
	"strings"
	"sync"
)

type UserStruct struct {
	Users map[string]user.User
	Mutex sync.Mutex
}

type SessionStruct struct {
	Sessions map[string]user.Session
	Mutex    sync.Mutex
}

type UserRepository struct {
	UserDB    UserStruct
	SessionDB SessionStruct
}

func (ur *UserRepository) CreateSession(session user.Session) {
	ur.SessionDB.Mutex.Lock()
	defer ur.SessionDB.Mutex.Unlock()
	ur.SessionDB.Sessions[session.SessionToken] = session
}

func (ur *UserRepository) GetSessionBySessionToken(token string) (user.Session, error) {
	ur.SessionDB.Mutex.Lock()
	defer ur.SessionDB.Mutex.Unlock()
	if session, exists := ur.SessionDB.Sessions[token]; exists {
		return session, nil
	}
	return user.Session{}, user.InvalidSessionError{"session doesn't exist"}
}

func (ur *UserRepository) GetUserByUsername(username string) (user.User, error) {
	username = strings.ToLower(username)

	ur.UserDB.Mutex.Lock()
	defer ur.UserDB.Mutex.Unlock()
	if user, exists := ur.UserDB.Users[username]; exists {
		return user, nil
	}
	return user.User{}, user.InvalidUserError{"user doesn't exist"}
}

func (ur *UserRepository) CreateUser(newUser user.User) error {
	username := strings.ToLower(newUser.Username)

	ur.UserDB.Mutex.Lock()
	defer ur.UserDB.Mutex.Unlock()
	if _, exists := ur.UserDB.Users[username]; exists {
		return user.InvalidUserError{"username"}
	}
	ur.UserDB.Users[username] = newUser
	return nil
}

func (ur *UserRepository) UpdateUser(username string, newData user.User) (user.User, error) {
	username = strings.ToLower(username)

	ur.UserDB.Mutex.Lock()
	defer ur.UserDB.Mutex.Unlock()
	if _, exists := ur.UserDB.Users[username]; !exists {
		return user.User{}, user.InvalidUserError{"user doesn't exist"}
	}

	if strings.HasPrefix(newData.AvatarURL, "data:") {
		const avatarStoragePath = "media/avatars/"
		pathToAvatar, err := common.DataURLToFile(avatarStoragePath + username, newData.AvatarURL, 500)
		if err != nil {
			return ur.UserDB.Users[username], user.InvalidImageError{"invalid image"}
		}
		newData.AvatarURL = pathToAvatar
	}

	ur.UserDB.Users[username] = newData
	return newData, nil
}

func (ur *UserRepository) ChangePassword(username string, newPSWD string) error {
	username = strings.ToLower(username)

	ur.UserDB.Mutex.Lock()
	defer ur.UserDB.Mutex.Unlock()
	if _, exists := ur.UserDB.Users[username]; !exists {
		return user.InvalidUserError{"user doesn't exist"}
	}
	data := ur.UserDB.Users[username]
	data.HashPassword = newPSWD
	ur.UserDB.Users[username] = data
	return nil
}

func (ur *UserRepository) RemoveSession(token string) error {
	ur.SessionDB.Mutex.Lock()
	defer ur.SessionDB.Mutex.Unlock()
	if _, exists := ur.SessionDB.Sessions[token]; !exists {
		return user.InvalidSessionError{"session doesn't exist"}
	}
	delete(ur.SessionDB.Sessions, token)
	return nil
}
