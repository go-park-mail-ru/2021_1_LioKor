package user

type UserRepository interface {
	CreateSession(session Session) error
	GetSessionBySessionToken(token string) (Session, error)
	GetUserByUsername(username string) (User, error)
	CreateUser(user User) error
	UpdateUser(username string, newData User) (User, error)
	ChangePassword(username string, newPSWD string) error
	RemoveSession(token string) error
	Close()
}
