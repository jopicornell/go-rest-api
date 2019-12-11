package services

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	. "github.com/go-jet/jet/postgres"
	"github.com/go-jet/jet/qrm"
	"github.com/jackc/pgx"
	"github.com/jmoiron/sqlx"
	"github.com/jopicornell/go-rest-api/db/entities/palmaactiva/image_gallery/model"
	. "github.com/jopicornell/go-rest-api/db/entities/palmaactiva/image_gallery/table"
	"github.com/jopicornell/go-rest-api/internals/errors"
	"github.com/jopicornell/go-rest-api/internals/models"
	"github.com/jopicornell/go-rest-api/pkg/password"
	"github.com/jopicornell/go-rest-api/pkg/server"
	"github.com/sirupsen/logrus"
	"time"
)

type UsersService interface {
	Login(username string, password string) (*models.Token, error)
	Register(user *model.Customer) (*model.Customer, error)
}

type usersService struct {
	db     *sqlx.DB
	server server.Server
}

func NewAuthService(db *sqlx.DB, server server.Server) UsersService {
	return &usersService{
		db:     db,
		server: server,
	}
}

func (s *usersService) Login(username string, passwd string) (token *models.Token, err error) {
	statement := SELECT(Customer.AllColumns, CustomerHasRoles.AllColumns).FROM(
		Customer.INNER_JOIN(CustomerHasRoles, CustomerHasRoles.UserID.EQ(Customer.UserID)),
	).
		WHERE(Customer.Username.EQ(String(username)))
	var user models.CustomerWithRoles
	logrus.Info(statement.DebugSql())
	if err := statement.Query(s.db, &user); err == nil {
		if err := password.ComparePasswordAndHash(passwd, user.Password); err != nil {
			return nil, errors.AuthUserNotMatched
		}
		hashing := sha512.New()
		if _, err := hashing.Write([]byte(fmt.Sprintf("%d;%s;%d", user.UserID, user.Password, time.Now().Nanosecond()))); err != nil {
			panic(err)
		}
		token := &models.Token{
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

func (s *usersService) Register(user *model.Customer) (_ *model.Customer, err error) {
	tx := s.db.MustBegin()
	statement := Customer.INSERT(
		Customer.FullName,
		Customer.Username,
		Customer.Password,
		Customer.NumComments,
		Customer.NumLikes,
		Customer.NumPictures,
	).
		MODEL(user).
		RETURNING(Customer.AllColumns)
	logrus.Info(statement.DebugSql())
	if err := statement.Query(tx, user); err == nil {
		statement = CustomerHasRoles.INSERT(CustomerHasRoles.AllColumns).MODEL(model.CustomerHasRoles{UserID: user.UserID, Role: models.USER_ROLE})
		logrus.Info(statement.DebugSql())
		if _, err := statement.Exec(tx); err != nil {

			logrus.Error(err)
			err = tx.Rollback()
			logrus.Panicf("error creating roles for a user: %w", err)
		}
		if err = tx.Commit(); err == nil {
			return user, nil
		} else {
			return nil, err
		}
	} else {
		errType, ok := err.(pgx.PgError)
		if ok && errType.ConstraintName == "customer_username_key" && errType.Code == "23505" {
			return nil, errors.UsernameExists
		}
		return nil, err
	}
}
