package services

import (
	"database/sql"
	"errors"
	. "github.com/go-jet/jet/postgres"
	"github.com/jmoiron/sqlx"
	. "github.com/jopicornell/go-rest-api/db/entities/palmaactiva/image_gallery/table"
	"github.com/jopicornell/go-rest-api/internals/models"
	"github.com/sirupsen/logrus"
)

type PicturesService interface {
	GetPictures() ([]models.Picture, error)
	GetPicture(uint) (*models.Picture, error)
	UpdatePicture(uint, *models.Picture) (*models.Picture, error)
	CreatePicture(*models.Picture, *models.User) (*models.Picture, error)
	DeletePicture(uint) error
}

type pictureService struct {
	db *sqlx.DB
}

var PictureNullError = errors.New("picture should not be null")

func NewPictureService(db *sqlx.DB) PicturesService {
	return &pictureService{
		db: db,
	}
}

func (s *pictureService) GetPictures() (pictures []models.Picture, err error) {
	pictures = []models.Picture{}
	statement := SELECT(Picture.AllColumns).FROM(
		Picture.INNER_JOIN(User.AS("usr"), Picture.UserID.EQ(User.AS("usr").UserID)),
	)
	sqlQuery := statement.DebugSql()
	logrus.Info(sqlQuery)
	if err = statement.Query(s.db, &pictures); err != nil {
		return nil, err
	}
	return pictures, nil
}

func (s *pictureService) GetPicture(id uint) (*models.Picture, error) {
	var picture models.Picture
	if err := s.db.Get(&picture, "SELECT * from image_gallery.pictures where id = ?", id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &picture, nil
}

func (s *pictureService) CreatePicture(picture *models.Picture, user *models.User) (*models.Picture, error) {
	insertQuery := "INSERT INTO pictures (start_date, duration, end_date, status, user_id) VALUES (?, ?, ?, ?, ?)"
	tx := s.db.MustBegin()
	if _, err := tx.Exec(insertQuery, picture.StartDate, picture.Duration, picture.EndDate, models.StatusPending, user.ID); err == nil {
		err = tx.Commit()
		return picture, err
	} else {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return nil, err
		}
		return nil, err
	}
}

func (s *pictureService) UpdatePicture(id uint, picture *models.Picture) (*models.Picture, error) {
	if picture == nil {
		return nil, PictureNullError
	}
	updateQuery := "UPDATE pictures SET start_date=?, duration=?, end_date=?, status=? where id = ?"
	tx := s.db.MustBegin()
	if _, err := tx.Exec(updateQuery, picture.StartDate, picture.Duration, picture.EndDate, picture.Status, picture.ID); err == nil {
		err = tx.Commit()
		return picture, err
	} else {
		errRollback := tx.Rollback()
		if errRollback != nil {
			panic(err)
		}
		return nil, err
	}
}

func (s *pictureService) DeletePicture(id uint) (err error) {
	deleteQuery := "DELETE FROM pictures WHERE id = ?"
	tx := s.db.MustBegin()
	if _, err := tx.Exec(deleteQuery, id); err == nil {
		err = tx.Commit()
		return err
	} else {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return errRollback
		}
		return err
	}
}
