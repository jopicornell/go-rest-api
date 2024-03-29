package services

import (
	"database/sql"
	"errors"
	"fmt"
	. "github.com/go-jet/jet/postgres"
	"github.com/go-jet/jet/qrm"
	"github.com/jackc/pgx"
	"github.com/jmoiron/sqlx"
	errors2 "github.com/jopicornell/go-rest-api/api/pictures/errors"
	"github.com/jopicornell/go-rest-api/api/pictures/responses"
	"github.com/jopicornell/go-rest-api/api/users/models"
	"github.com/jopicornell/go-rest-api/db/entities/palmaactiva/image_gallery/model"
	. "github.com/jopicornell/go-rest-api/db/entities/palmaactiva/image_gallery/table"
	"github.com/jopicornell/go-rest-api/pkg/server"
	"github.com/sirupsen/logrus"
)

type PicturesService interface {
	GetPictures() ([]responses.PictureWithImages, error)
	GetPicture(uint) (*responses.PictureWithImages, error)
	GetPictureComments(uint) ([]responses.Comment, error)
	GetPictureLikes(uint) ([]responses.Like, error)
	UpdatePicture(uint, *model.Picture) (*responses.PictureWithImages, error)
	CreatePicture(*model.Picture, *models.UserWithRoles) (*responses.Picture, error)
	DeletePicture(id uint, user *models.UserWithRoles) error
	CreatePictureComment(int32, *model.Comment) (*responses.Comment, error)
	CreatePictureLike(pictureId int32, userID int32) error
	DeletePictureLike(id int32, userId int32) error
	DeletePictureComment(id int32, commentId int32) error
	CheckUserAccess(user *models.UserWithRoles, model interface{}) bool
}

type pictureService struct {
	server server.Server
	db     *sqlx.DB
}

type GroupedUsers struct {
	UserId      int32
	NumComments int32
	NumPictures int32
	NumLikes    int32
}

var PictureNullError = errors.New("picture should not be null")

func NewPictureService(server server.Server) PicturesService {
	return &pictureService{
		server: server,
		db:     server.GetRelationalDatabase(),
	}
}

func (s *pictureService) GetPictures() (pictures []responses.PictureWithImages, err error) {
	pictures = []responses.PictureWithImages{}
	statement := SELECT(Picture.AllColumns, Image.AllColumns).FROM(
		Picture.INNER_JOIN(Image, Picture.ImageID.EQ(Image.ImageID)),
	).
		ORDER_BY(Picture.Created.DESC())
	sqlQuery := statement.DebugSql()
	logrus.Info(sqlQuery)
	if err = statement.Query(s.db, &pictures); err != nil {
		return nil, err
	}
	return pictures, nil
}

func (s *pictureService) GetPicture(id uint) (*responses.PictureWithImages, error) {
	picture := new(responses.PictureWithImages)
	statement := SELECT(
		Picture.AllColumns,
		Image.AllColumns,
		User.AS("user_without_avatar").FullName,
		User.AS("user_without_avatar").Username,
		User.AS("user_without_avatar").UserID,
	).
		FROM(Picture.
			INNER_JOIN(Image, Picture.ImageID.EQ(Image.ImageID)).
			INNER_JOIN(User.AS("user_without_avatar"), User.AS("user_without_avatar").UserID.EQ(Picture.UserID))).
		WHERE(Picture.PictureID.EQ(Int(int64(id))))
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
	statement := SELECT(Comment.AllColumns, User.UserID, User.AvatarID, User.FullName, Image.AllColumns).
		FROM(Comment.
			INNER_JOIN(User, User.UserID.EQ(Comment.UserID)).
			LEFT_JOIN(Image, Image.ImageID.EQ(User.AvatarID))).
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

func (s *pictureService) GetComment(id uint) (*responses.Comment, error) {
	comment := new(responses.Comment)
	statement := SELECT(Comment.AllColumns, User.UserID, User.AvatarID, User.FullName, Image.AllColumns).
		FROM(Comment.
			INNER_JOIN(User, User.UserID.EQ(Comment.UserID)).
			LEFT_JOIN(Image, Image.ImageID.EQ(User.AvatarID))).
		WHERE(Comment.CommentID.EQ(Int(int64(id))))
	sqlQuery := statement.DebugSql()
	logrus.Info(sqlQuery)
	if err := statement.Query(s.db, comment); err != nil {
		if err == qrm.ErrNoRows {
			return nil, errors2.CommentNotFound
		}
		return nil, err
	}
	return comment, nil
}

func (s *pictureService) CreatePicture(picture *model.Picture, user *models.UserWithRoles) (*responses.Picture, error) {
	tx := s.db.MustBegin()
	statement := Picture.INSERT(Picture.ImageID, Picture.UserID, Picture.Title, Picture.Description).
		VALUES(picture.ImageID, picture.UserID, picture.Title, picture.Description).
		RETURNING(Picture.AllColumns)
	createdPicture := new(responses.Picture)
	if err := statement.Query(tx, createdPicture); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, fmt.Errorf("error creating picture %v and additional rollback error %v", rollbackErr, err)
		}
		return nil, err
	}
	updateUser := User.UPDATE(User.NumPictures).
		SET(User.NumPictures.ADD(Int(1))).
		WHERE(User.UserID.EQ(Int(int64(user.UserID))))
	if _, err := updateUser.Exec(tx); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return createdPicture, nil
}

func (s *pictureService) UpdatePicture(id uint, picture *model.Picture) (*responses.PictureWithImages, error) {
	updateStatement := Picture.UPDATE(Picture.Title, Picture.Description).
		MODEL(picture).
		WHERE(Picture.PictureID.EQ(Int(int64(id))))
	logrus.Info(updateStatement.DebugSql())
	if _, err := updateStatement.Exec(s.db); err != nil {
		return nil, err
	}
	return s.GetPicture(id)
}

func (s *pictureService) DeletePicture(id uint, user *models.UserWithRoles) (err error) {
	tx := s.db.MustBegin()
	pictureSelect := SELECT(Picture.AllColumns).
		FROM(Picture).
		WHERE(Picture.PictureID.EQ(Int(int64(id)))).
		FOR(UPDATE())
	picture := new(model.Picture)
	if err := pictureSelect.Query(tx, picture); err != nil {
		tx.Rollback()
		if err == qrm.ErrNoRows {
			return errors2.PictureNotFound
		}
		return err
	}
	pictureComments := getGroupedUsersForDelete(int64(id))
	logrus.Info(pictureComments.DebugSql())
	commentsAndLikes := make([]GroupedUsers, 0)
	if err := pictureComments.Query(tx, &commentsAndLikes); err != nil {
		tx.Rollback()
		return err
	}
	for _, grouped := range commentsAndLikes {
		updateStatement := User.UPDATE(User.NumComments, User.NumLikes, User.NumPictures).
			SET(grouped.NumComments, grouped.NumLikes, grouped.NumPictures).
			WHERE(User.UserID.EQ(Int(int64(grouped.UserId))))
		if _, err := updateStatement.Exec(tx); err != nil {
			errRollback := tx.Rollback()
			if errRollback != nil {
				return errRollback
			}
			return err
		}
	}
	deletePictureStatement := Picture.DELETE().
		WHERE(Picture.PictureID.EQ(Int(int64(id))))
	if _, err := deletePictureStatement.Exec(tx); err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return errRollback
		}
		return err
	}
	deleteImageStatement := Image.DELETE().WHERE(Image.ImageID.EQ(Int(int64(picture.ImageID))))
	if _, err := deleteImageStatement.Exec(tx); err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return errRollback
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *pictureService) CreatePictureComment(id int32, comment *model.Comment) (*responses.Comment, error) {
	returnComment := new(responses.Comment)
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
	updatePicture := Picture.UPDATE(Picture.NumComments).
		SET(Picture.NumComments.ADD(Int(1))).
		WHERE(Picture.PictureID.EQ(Int(int64(id))))
	if _, err := updatePicture.Exec(tx); err != nil {
		if errRollBack := tx.Rollback(); errRollBack != nil {
			logrus.Errorf("Error rolling back", errRollBack)
		}
		return nil, err
	}
	updateUser := User.UPDATE(User.NumComments).
		SET(User.NumComments.ADD(Int(1))).
		WHERE(User.UserID.EQ(Int(int64(comment.UserID))))
	if _, err := updateUser.Exec(tx); err != nil {
		if errRollBack := tx.Rollback(); errRollBack != nil {
			logrus.Errorf("Error rolling back", errRollBack)
		}
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return returnComment, nil
}

func (s *pictureService) GetPictureLikes(id uint) ([]responses.Like, error) {
	likes := make([]responses.Like, 0)
	statement := SELECT(Like.UserID, Like.PictureID).
		FROM(Like.INNER_JOIN(User, User.UserID.EQ(Like.UserID))).
		WHERE(Like.PictureID.EQ(Int(int64(id))))
	sqlQuery := statement.DebugSql()
	logrus.Info(sqlQuery)
	if err := statement.Query(s.db, &likes); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return likes, nil
}

func (s *pictureService) CreatePictureLike(id int32, userId int32) error {
	statement := Like.INSERT(Like.PictureID, Like.UserID).
		VALUES(id, userId)
	logrus.Info(statement.DebugSql())
	tx := s.db.MustBegin()
	if _, err := statement.Exec(tx); err != nil {
		errType, ok := err.(pgx.PgError)
		if ok && errType.Code == "23505" {
			return errors2.UserAlreadyLikedPicture
		}
		if errRollBack := tx.Rollback(); errRollBack != nil {
			logrus.Errorf("Error rolling back", errRollBack)
		}
		return err
	}
	updatePicture := Picture.UPDATE(Picture.NumLikes).
		SET(Picture.NumLikes.ADD(Int(1))).
		WHERE(Picture.PictureID.EQ(Int(int64(id))))
	if _, err := updatePicture.Exec(tx); err != nil {
		return err
	}
	updateUser := User.UPDATE(User.NumLikes).
		SET(User.NumLikes.ADD(Int(1))).
		WHERE(User.UserID.EQ(Int(int64(userId))))
	if _, err := updateUser.Exec(tx); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (s *pictureService) DeletePictureLike(id int32, userId int32) error {
	statement := Like.DELETE().
		WHERE(Like.UserID.EQ(Int(int64(userId))).
			AND(Like.PictureID.EQ(Int(int64(id)))))
	logrus.Info(statement.DebugSql())
	tx := s.db.MustBegin()
	result, err := statement.Exec(tx)
	numRows, _ := result.RowsAffected()
	if err != nil || numRows == 0 {
		if errRollBack := tx.Rollback(); errRollBack != nil {
			logrus.Errorf("Error rolling back", errRollBack)
		}
		if numRows == 0 {
			return errors2.PictureNotFound
		}
		return err
	}
	updatePicture := Picture.UPDATE(Picture.NumLikes).
		SET(Picture.NumLikes.SUB(Int(1))).
		WHERE(Picture.PictureID.EQ(Int(int64(id))))
	if _, err := updatePicture.Exec(tx); err != nil {
		return err
	}
	updateUser := User.UPDATE(User.NumLikes).
		SET(User.NumLikes.SUB(Int(1))).
		WHERE(User.UserID.EQ(Int(int64(userId))))
	if _, err := updateUser.Exec(tx); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (s *pictureService) DeletePictureComment(id int32, commentId int32) error {
	comment, err := s.GetComment(uint(commentId))
	if err != nil {
		if err == qrm.ErrNoRows {
			return errors2.PictureNotFound
		}
		return err
	}
	statement := Comment.DELETE().
		WHERE(Comment.CommentID.EQ(Int(int64(commentId))).
			AND(Comment.PictureID.EQ(Int(int64(id)))))
	logrus.Info(statement.DebugSql())
	tx := s.db.MustBegin()
	result, err := statement.Exec(tx)
	numRows, _ := result.RowsAffected()
	if err != nil || numRows == 0 {
		if errRollBack := tx.Rollback(); errRollBack != nil {
			logrus.Errorf("Error rolling back", errRollBack)
		}
		if numRows == 0 {
			return errors2.PictureNotFound
		}
		return err
	}
	updatePicture := Picture.UPDATE(Picture.NumLikes).
		SET(Picture.NumLikes.SUB(Int(1))).
		WHERE(Picture.PictureID.EQ(Int(int64(id))))
	if _, err := updatePicture.Exec(tx); err != nil {
		return err
	}
	updateUser := User.UPDATE(User.NumLikes).
		SET(User.NumLikes.SUB(Int(1))).
		WHERE(User.UserID.EQ(Int(int64(comment.UserID))))
	if _, err := updateUser.Exec(tx); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (s *pictureService) CheckUserAccess(user *models.UserWithRoles, modelCheck interface{}) bool {
	if user.HasRole(models.ADMIN_ROLE) {
		return true
	}
	switch modelCheck.(type) {
	case *responses.PictureWithImages:
		return user.UserID == modelCheck.(*responses.PictureWithImages).UserID
	case *model.Comment:
		return user.UserID == modelCheck.(*model.Comment).UserID
	case *model.Like:
		return user.UserID == modelCheck.(*model.Like).UserID
	}
	return false
}

func getGroupedUsersForDelete(id int64) SelectStatement {
	return SELECT(
		User.UserID.
			AS("grouped_users.user_id"),
		User.NumComments.SUB(COUNT(Comment.CommentID)).
			AS("grouped_users.num_comments"),
		CASE().
			WHEN(COUNT(Like.PictureID).GT(Int(0))).
			THEN(User.NumLikes.SUB(Int(1))).
			ELSE(User.NumLikes).
			AS("grouped_users.num_likes"),
		CASE().
			WHEN(Picture.UserID.EQ(User.UserID)).
			THEN(User.NumPictures.SUB(Int(1))).
			ELSE(User.NumPictures).
			AS("grouped_users.num_pictures"),
	).
		FROM(
			Picture.
				LEFT_JOIN(Comment, Comment.PictureID.EQ(Picture.PictureID)).
				LEFT_JOIN(Like, Like.PictureID.EQ(Picture.PictureID)).
				LEFT_JOIN(User, User.UserID.EQ(Comment.UserID).OR(User.UserID.EQ(Like.UserID)))).
		WHERE(Picture.PictureID.EQ(Int(id))).
		GROUP_BY(User.UserID, Picture.UserID, User.NumComments, User.NumLikes, User.NumPictures).
		HAVING(User.UserID.IS_NOT_NULL())
}
