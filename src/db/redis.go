package db

import (
	"context"
	"encoding/json"

	"github.com/codemunsta/risevest-test/src/models"
	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

type RedisSession struct {
	AuthToken string
	User      models.User
}

type RedisAdminSession struct {
	AuthToken string
	Admin     models.Admin
}

func GetSession(authToken string) (RedisSession, bool) {
	sessionJson, err := RedisClient.Get(context.Background(), "Session:"+authToken).Result()

	sessionFound := false
	var session RedisSession
	if err == redis.Nil || err != nil {
		sessionFound = false
	} else {
		sessionFound = true
	}
	if sessionFound {
		err2 := json.Unmarshal([]byte(sessionJson), &session)
		if err2 != nil {
			return session, false
		}
		return session, true
	} else {
		return session, false
	}
}

func CreateSession(authToken string, user models.User) error {
	var session RedisSession
	session.AuthToken = authToken
	session.User = user
	sessionJSON, err := json.Marshal(session)
	if err != nil {
		return err
	} else {
		err = RedisClient.Set(context.Background(), "Session:"+session.AuthToken, string(sessionJSON), 600000000000).Err()
		return err
	}
}

func GetAdminSession(authToken string) (RedisSession, bool) {
	sessionJson, err := RedisClient.Get(context.Background(), "AdminSession:"+authToken).Result()

	sessionFound := false
	var session RedisSession
	if err == redis.Nil || err != nil {
		sessionFound = false
	} else {
		sessionFound = true
	}
	if sessionFound {
		err2 := json.Unmarshal([]byte(sessionJson), &session)
		if err2 != nil {
			return session, false
		}
		return session, true
	} else {
		return session, false
	}
}

func CreateAdminSession(authToken string, admin models.Admin) error {
	var session RedisAdminSession
	session.AuthToken = authToken
	session.Admin = admin
	sessionJSON, err := json.Marshal(session)
	if err != nil {
		return err
	} else {
		err = RedisClient.Set(context.Background(), "AdminSession:"+session.AuthToken, string(sessionJSON), 600000000000).Err()
		return err
	}
}
