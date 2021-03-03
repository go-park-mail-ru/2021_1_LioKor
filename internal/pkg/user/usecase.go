package user

type UseCase interface {
	Login(credentials Credentials) error
	CreateSession(username string) (SessionToken, error)
	GetUserBySessionToken(sessionToken string) (User, error)
	SignUp(newUser UserSignUp) error
	UpdateUser(username string, newData User) (User, error)
}
