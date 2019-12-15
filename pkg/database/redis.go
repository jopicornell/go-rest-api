package database

import (
	"encoding/json"
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

func (r *Redis) InitializeClient() {
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
	r.client = client
}

func (r *Redis) GetStruct(key string, ifc interface{}) {
	cmd := r.client.Get(key)
	if cmd.Err() != nil {
		if cmd.Err() != redis.Nil {
			panic(cmd.Err())
		}
		return
	}
	if err := json.Unmarshal([]byte(cmd.Val()), ifc); err != nil {
		panic(err)
	}
}

func (r *Redis) SetStruct(key string, m interface{}) {
	var cmd *redis.StatusCmd
	if bytes, err := json.Marshal(&m); err != nil {
		panic(err)
	} else {
		cmd = r.client.Set(key, string(bytes), time.Hour*24*7)
	}
	if cmd.Err() != nil {
		panic(cmd.Err())
	}
}
