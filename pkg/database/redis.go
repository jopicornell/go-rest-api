package database

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"log"
	"time"
)

type Redis struct {
	client   *redis.Client
	Host     string
	Password string
}

func (m *Redis) InitializeClient() {
	var client *redis.Client
	err := Retry(func() (err error) {
		client = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("localhost:6379"),
			Password: "",
			DB:       0,
		})
		err = client.Ping().Err()

		if err != nil {
			logrus.Errorf("Error connecting to client %s", err)
		}
		return err
	}, time.Second*15, time.Minute*5)
	if err != nil {
		wrapError := errors.Wrap(err, "some problem with initializing relational client")
		log.Fatal(wrapError.Error())
	}
	log.Println("New connection to redis")
	m.client = client
}
