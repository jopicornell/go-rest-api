package services

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	. "github.com/go-jet/jet/postgres"
	"github.com/go-jet/jet/qrm"
	"github.com/jackc/pgx"
	"github.com/jmoiron/sqlx"
	"github.com/jopicornell/go-rest-api/api/users/errors"
	"github.com/jopicornell/go-rest-api/api/users/models"
	"github.com/jopicornell/go-rest-api/api/users/responses"
	"github.com/jopicornell/go-rest-api/db/entities/palmaactiva/image_gallery/model"
	. "github.com/jopicornell/go-rest-api/db/entities/palmaactiva/image_gallery/table"
	"github.com/jopicornell/go-rest-api/pkg/password"
	"github.com/jopicornell/go-rest-api/pkg/server"
	"github.com/sirupsen/logrus"
	"time"
)

type UsersService interface {
	Login(username string, password string) (*responses.Token, error)
	Register(user *model.User) (*responses.User, error)
	GetUser(id uint) (*responses.User, error)
	UpdateUser(id uint, user *model.User) (*responses.User, error)
	CheckUsername(username string) error
	CheckUserAccess(user *models.UserWithRoles, userToModify *responses.User) bool
}

type usersService struct {
	db     *sqlx.DB
	server server.Server
}

func NewAuthService(server server.Server) UsersService {
	return &usersService{
		db:     server.GetRelationalDatabase(),
		server: server,
	}
}

func (s *usersService) Login(username string, passwd string) (token *responses.Token, err error) {
	statement := SELECT(User.AllColumns, UserHasRoles.AllColumns).FROM(
		User.INNER_JOIN(UserHasRoles, UserHasRoles.UserID.EQ(User.UserID)),
	).
		WHERE(User.Username.EQ(String(username)))
	user := new(models.UserWithRoles)
	logrus.Info(statement.DebugSql())
	if err := statement.Query(s.db, user); err == nil {
		if err := password.ComparePasswordAndHash(passwd, user.Password); err != nil {
			return nil, errors.AuthUserNotMatched
		}
		hashing := sha512.New()
		if _, err := hashing.Write([]byte(fmt.Sprintf("%d;%s;%d", user.UserID, user.Password, time.Now().Nanosecond()))); err != nil {
			panic(err)
		}
		token := &responses.Token{
			Token:  base64.StdEncoding.EncodeToString(hashing.Sum(nil)),
			UserID: int(user.UserID),
			Roles:  user.GetRoles(),
		}
		s.server.GetCache().SetStruct(token.Token, user)
		return token, nil
	} else {
		if err == qrm.ErrNoRows {
			return nil, errors.AuthUserNotMatched
		}
		return nil, err
	}
}

func (s *usersService) Register(user *model.User) (_ *responses.User, err error) {
	tx := s.db.MustBegin()
	statement := User.INSERT(
		User.FullName,
		User.Username,
		User.Password,
		User.NumComments,
		User.NumLikes,
		User.NumPictures,
	).
		MODEL(user).
		RETURNING(
			User.UserID,
			User.FullName,
			User.Username,
		)
	logrus.Info(statement.DebugSql())
	responseUser := new(responses.User)
	if err := statement.Query(tx, responseUser); err == nil {
		statement = UserHasRoles.INSERT(UserHasRoles.AllColumns).MODEL(model.UserHasRoles{UserID: responseUser.UserID, Role: models.USER_ROLE})
		logrus.Info(statement.DebugSql())
		if _, err := statement.Exec(tx); err != nil {

			logrus.Error(err)
			err = tx.Rollback()
			logrus.Panicf("error creating roles for a user: %w", err)
		}
		if err = tx.Commit(); err == nil {
			return responseUser, nil
		} else {
			return nil, err
		}
	} else {
		errType, ok := err.(pgx.PgError)
		if ok && errType.ConstraintName == "user_username_key" && errType.Code == "23505" {
			return nil, errors.UsernameExists
		}
		return nil, err
	}
}

func (s *usersService) GetUser(id uint) (*responses.User, error) {
	user := new(responses.User)
	statement := SELECT(
		User.UserID,
		User.Username,
		User.FullName,
		User.NumPictures,
		User.NumLikes,
		User.NumComments,
		User.AvatarID,
		Image.AllColumns,
	).FROM(User.LEFT_JOIN(Image, Image.ImageID.EQ(User.AvatarID))).
		WHERE(User.UserID.EQ(Int(int64(id))))
	logrus.Info(statement.DebugSql())
	if err := statement.Query(s.db, user); err != nil {
		if err == qrm.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (s *usersService) UpdateUser(id uint, user *model.User) (*responses.User, error) {
	updateStatement := User.
		UPDATE(User.AvatarID, User.FullName, User.Username).
		MODEL(user).
		WHERE(User.UserID.EQ(Int(int64(id))))
	result, err := updateStatement.Exec(s.db)
	if err != nil {
		return nil, err
	}
	if numRows, err := result.RowsAffected(); err != nil {
		return nil, err
	} else if numRows == 0 {
		return nil, errors.UserNotFound
	}
	return s.GetUser(id)
}

func (s *usersService) CheckUsername(username string) error {
	user := new(responses.User)
	statement := SELECT(
		User.UserID,
	).FROM(User).
		WHERE(User.Username.EQ(String(username)))
	logrus.Info(statement.DebugSql())
	if err := statement.Query(s.db, user); err != nil {
		if err == qrm.ErrNoRows {
			return errors.UserNotFound
		}
	}
	return nil
}

func (s *usersService) CheckUserAccess(user *models.UserWithRoles, userToModify *responses.User) bool {
	if user.HasRole(models.ADMIN_ROLE) {
		return true
	}
	return user.UserID == userToModify.UserID

}
