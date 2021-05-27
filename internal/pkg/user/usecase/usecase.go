package usecase

import (
	"context"
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"liokor_mail/internal/pkg/common"
	session "liokor_mail/internal/pkg/common/protobuf_sessions"
	"liokor_mail/internal/pkg/user"
	"liokor_mail/internal/pkg/user/validators"
	"log"
	"os"
	"strings"
	"time"
)

type UserUseCase struct {
	Repository     user.UserRepository
	SessionManager session.IsAuthClient
	Config         common.Config
}

func (uc *UserUseCase) Login(credentials user.Credentials) error {
	loginUser, err := uc.Repository.GetUserByUsername(credentials.Username)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(loginUser.HashPassword), []byte(credentials.Password))

	if err != nil {
		return common.InvalidUserError{"Invalid credentials"}
	}

	return nil
}

func (uc *UserUseCase) Logout(sessionToken string) error {
	_, err := uc.SessionManager.Delete(
		context.Background(),
		&session.SessionToken{
			SessionToken: sessionToken,
		},
	)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				return common.InvalidSessionError{Message: e.Message()}
			}
		}
		return err
	}

	return nil
}

func (uc *UserUseCase) CreateSession(username string) (common.Session, error) {
	sessionUser, err := uc.Repository.GetUserByUsername(username)
	if err != nil {
		return common.Session{}, err
	}

	newSession := session.Session{
		UserId:       int32(sessionUser.Id),
		SessionToken: common.GenerateRandomString(),
		Expiration:   timestamppb.New(time.Now().Add(10 * 24 * time.Hour)),
	}

	s, err := uc.SessionManager.Create(
		context.Background(),
		&newSession,
	)

	if err != nil {
		return common.Session{}, err
	}

	return common.Session{
		UserId:       int(s.UserId),
		SessionToken: s.SessionToken,
		Expiration:   s.Expiration.AsTime(),
	}, nil
}

func (uc *UserUseCase) SignUp(newUser user.UserSignUp) error {

	if !validators.ValidateUsername(newUser.Username) {
		return user.InvalidUsernameError{"invalid username"}
	}
	if !validators.ValidatePassword(newUser.Password) {
		return user.WeakPasswordError{"password is too weak"}
	}

	hashPSWD, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = uc.Repository.CreateUser(user.User{
		0,
		newUser.Username,
		string(hashPSWD),
		common.NullString{sql.NullString{String: newUser.AvatarURL, Valid: true}},
		newUser.FullName,
		newUser.ReserveEmail,
		time.Now().String(),
		false,
	})
	if err != nil {
		return err
	}

	return nil
}

func (uc *UserUseCase) UpdateUser(username string, newData user.User) (user.User, error) {
	sessionUser, err := uc.Repository.GetUserByUsername(username)
	if err != nil {
		return user.User{}, err
	}
	if newData.FullName != sessionUser.FullName {
		sessionUser.FullName = newData.FullName
	}
	if newData.ReserveEmail != sessionUser.ReserveEmail {
		sessionUser.ReserveEmail = newData.ReserveEmail
	}

	sessionUser, err = uc.Repository.UpdateUser(sessionUser)
	if err != nil {
		return user.User{}, err
	}
	return sessionUser, nil
}

func (uc *UserUseCase) UploadImage(username string, dataUrl string) (string, error) {
	if !strings.HasPrefix(dataUrl, "data:") {
		return "", common.InvalidImageError{"invalid image"}
	}

	imageFileName := common.GenerateRandomString()
	pathToImage, err := common.DataURLToFile(uc.Config.FileStoragePath+imageFileName, dataUrl, 512)
	if err != nil {
		log.Println(err.Error())
		return "", common.InvalidImageError{"invalid image"}
	}

	err = uc.Repository.AddUploadedFile(username, pathToImage)
	if err != nil {
		_ = os.Remove(pathToImage)
		return "", err
	}

	return pathToImage, nil
}

func (uc *UserUseCase) UpdateAvatar(username string, newAvatar string) (user.User, error) {
	sessionUser, err := uc.Repository.GetUserByUsername(username)
	if err != nil {
		return user.User{}, err
	}

	if strings.HasPrefix(newAvatar, "data:") {
		avatarFileName := common.GenerateRandomString()
		pathToAvatar, err := common.DataURLToFile(uc.Config.AvatarStoragePath+avatarFileName, newAvatar, 500)
		if err != nil {
			log.Println(err.Error())
			return sessionUser, common.InvalidImageError{"invalid image"}
		}
		if len(sessionUser.AvatarURL.String) > 0 {
			_ = os.Remove(sessionUser.AvatarURL.String)
		}
		sessionUser.AvatarURL.String = pathToAvatar
		sessionUser.AvatarURL.Valid = true
	} else {
		return sessionUser, common.InvalidImageError{"invalid image"}
	}

	sessionUser, err = uc.Repository.UpdateAvatar(username, sessionUser.AvatarURL)
	if err != nil {
		return user.User{}, err
	}
	return sessionUser, nil
}

func (uc *UserUseCase) GetUserByUsername(username string) (user.User, error) {
	requestedUser, err := uc.Repository.GetUserByUsername(username)
	if err != nil {
		return user.User{}, err
	}

	return requestedUser, nil
}

func (uc *UserUseCase) GetUserById(id int) (user.User, error) {
	requestedUser, err := uc.Repository.GetUserById(id)
	if err != nil {
		return user.User{}, err
	}

	return requestedUser, nil
}

func (uc *UserUseCase) ChangePassword(sessionUser user.User, changePSWD user.ChangePassword) error {
	if !validators.ValidatePassword(changePSWD.NewPassword) {
		return common.InvalidUserError{"invalid password"}
	}

	err := bcrypt.CompareHashAndPassword([]byte(sessionUser.HashPassword), []byte(changePSWD.OldPassword))

	if err != nil {
		return common.InvalidUserError{"Invalid password"}
	}

	hashPSWD, err := bcrypt.GenerateFromPassword([]byte(changePSWD.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return uc.Repository.ChangePassword(sessionUser.Username, string(hashPSWD))

}
