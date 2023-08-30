package db

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/codemunsta/risevest-test/src/models"
	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client = redis.NewClient(&redis.Options{
	Addr: "redis:6379",
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

func GetAdminSession(authToken string) (RedisAdminSession, bool) {
	sessionJson, err := RedisClient.Get(context.Background(), "AdminSession:"+authToken).Result()

	sessionFound := false
	var session RedisAdminSession
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

func GetOrCreateFolderRedis(folder models.Folder) (models.Folder, error) {
	folderString, err := RedisClient.Get(context.Background(), "folder:"+fmt.Sprintf("%v", folder.UserID)+folder.Name).Result()
	if err == redis.Nil {
		folderJson, err := json.Marshal(folder)
		if err != nil {
			var empty models.Folder
			return empty, err
		} else {
			err = RedisClient.Set(context.Background(), "folder:"+fmt.Sprintf("%v", folder.UserID)+folder.Name, string(folderJson), 900000000000).Err()
			return folder, err
		}
	} else if err != nil {
		var empty models.Folder
		return empty, err
	} else {
		var folder_ models.Folder
		err2 := json.Unmarshal([]byte(folderString), &folder_)
		if err2 != nil {
			var empty models.Folder
			return empty, err
		} else {
			return folder_, nil
		}
	}
}
