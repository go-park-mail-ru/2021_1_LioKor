package repository

import (
	"errors"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
	"liokor_mail/internal/pkg/common"
)

type GormPostgresSessionRepository struct {
	DBInstance common.GormPostgresDataBase
}

func (gsr *GormPostgresSessionRepository) Create(session common.Session) error {
	result := gsr.DBInstance.DB.Create(&session)
	if err := result.Error; err != nil {
		if pgerr, ok := err.(*pgconn.PgError); ok {
			if pgerr.ConstraintName == "sessions_user_id_fkey" {
				return common.InvalidUserError{"user doesn't exist"}
			} else if pgerr.ConstraintName == "sessions_pkey" {
				return common.InvalidSessionError{"sessionToken exists"}
			}
		}
		return err
	}
	return nil
}
func (gsr *GormPostgresSessionRepository) Get(token string) (common.Session, error) {
	var s common.Session
	err := gsr.DBInstance.DB.Where("token = ?", token).First(&s).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.Session{}, common.InvalidSessionError{"session doesn't exist"}
		} else {
			return common.Session{}, err
		}
	}
	return s, nil
}
func (gsr *GormPostgresSessionRepository) Delete(token string) error {
	s, err := gsr.Get(token)
	if err != nil{
		return err
	}
	gsr.DBInstance.DB.Where("token = ?", token).Delete(&s)
	if err = gsr.DBInstance.DB.Error; err != nil || gsr.DBInstance.DB.RowsAffected == 0{
		return err
	}
	return nil
}
