package services

import (
	"database/sql"
	"errors"
	. "github.com/go-jet/jet/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/jopicornell/go-rest-api/db/entities/palmaactiva/image_gallery/model"
	. "github.com/jopicornell/go-rest-api/db/entities/palmaactiva/image_gallery/table"
	"github.com/jopicornell/go-rest-api/internals/pictures/responses"
	"github.com/sirupsen/logrus"
)

type PicturesService interface {
	GetPictures() ([]responses.PictureWithImages, error)
	GetPicture(uint) (*responses.PictureWithImages, error)
	GetPictureComments(uint) ([]responses.Comment, error)
	UpdatePicture(uint, *model.Picture) (*model.Picture, error)
	CreatePicture(*model.Picture, *model.Customer) (*model.Picture, error)
	DeletePicture(uint) error
	CreatePictureComment(int32, *model.Comment) (*model.Comment, error)
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

func (s *pictureService) GetPictures() (pictures []responses.PictureWithImages, err error) {
	pictures = []responses.PictureWithImages{}
	statement := SELECT(Picture.AllColumns, Image.AllColumns).FROM(
		Picture.INNER_JOIN(Image, Picture.ImageID.EQ(Image.ImageID)),
	)
	sqlQuery := statement.DebugSql()
	logrus.Info(sqlQuery)
	if err = statement.Query(s.db, &pictures); err != nil {
		return nil, err
	}
	return pictures, nil
}

func (s *pictureService) GetPicture(id uint) (*responses.PictureWithImages, error) {
	picture := new(responses.PictureWithImages)
	statement := SELECT(Picture.AllColumns, Image.AllColumns).FROM(
		Picture.INNER_JOIN(Image, Picture.ImageID.EQ(Image.ImageID)),
	).WHERE(Picture.PictureID.EQ(Int(int64(id))))
	sqlQuery := statement.DebugSql()
	logrus.Info(sqlQuery)
	if err := statement.Query(s.db, picture); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return picture, nil
}

func (s *pictureService) GetPictureComments(id uint) ([]responses.Comment, error) {
	comments := make([]responses.Comment, 0)
	statement := SELECT(Comment.AllColumns).
		FROM(Comment).
		WHERE(Comment.PictureID.EQ(Int(int64(id))))
	sqlQuery := statement.DebugSql()
	logrus.Info(sqlQuery)
	if err := statement.Query(s.db, &comments); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return comments, nil
}

func (s *pictureService) CreatePicture(picture *model.Picture, user *model.Customer) (*model.Picture, error) {
	statement := Picture.INSERT(Picture.ImageID, Picture.UserID, Picture.Title, Picture.Description).
		VALUES(picture.ImageID, picture.UserID, picture.Title, picture.Description).
		RETURNING(Picture.AllColumns)
	createdPicture := new(model.Picture)
	if err := statement.Query(s.db, createdPicture); err == nil {
		return createdPicture, nil
	} else {
		return nil, err
	}
}

func (s *pictureService) UpdatePicture(id uint, picture *model.Picture) (*model.Picture, error) {
	if picture == nil {
		return nil, PictureNullError
	}
	return nil, nil
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

func (s *pictureService) CreatePictureComment(id int32, comment *model.Comment) (*model.Comment, error) {
	returnComment := new(model.Comment)
	statement := Comment.INSERT(Comment.PictureID, Comment.UserID, Comment.Comment).
		VALUES(id, comment.UserID, comment.Comment).
		RETURNING(Comment.AllColumns)
	logrus.Info(statement.DebugSql())
	tx := s.db.MustBegin()
	if err := statement.Query(tx, returnComment); err != nil {
		if errRollBack := tx.Rollback(); errRollBack != nil {
			logrus.Errorf("Error rolling back", errRollBack)
		}
		return nil, err
	}
	selectPicture := SELECT(Picture.NumComments).
		FROM(Picture).
		WHERE(Picture.PictureID.EQ(Int(int64(id)))).
		FOR(UPDATE())
	picture := new(model.Picture)
	if err := selectPicture.Query(tx, picture); err != nil {
		if errRollBack := tx.Rollback(); errRollBack != nil {
			logrus.Errorf("Error rolling back", errRollBack)
		}
		return nil, err
	}
	picture.NumComments++
	updatePicture := Picture.UPDATE(Picture.NumComments).
		SET(picture.NumComments).
		WHERE(Picture.PictureID.EQ(Int(int64(id))))
	if _, err := updatePicture.Exec(tx); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return comment, nil
}
