package models

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const userPrefix = "pigo:user:"

type User struct {
	Id       string `json:"id" dt:"keyword"`
	Name     string `json:"name" dt:"keyword"`
	Password string `json:"password" dt:"keyword"`
	tokenId  string
}

func (u *User) setId() (err error) {
	u.Id, err = getUUID("user")
	return
}

// UserLogin handle the login and return a tokenId
func (u *User) Login() (string, error) {
	logger.Printf("Login user #%s", u.Name)
	tokenid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	// todo: check login user
	u.tokenId = tokenid.String()
	// remove Password from u before cache
	u.Password = ""
	uj, err := json.Marshal(*u)
	if err != nil {
		return "", err
	}
	_, err = redisCli.Set(userPrefix+u.tokenId, string(uj), 2*time.Hour).Result()
	if err != nil {
		return "", err
	}
	return u.tokenId, nil
}

func (u *User) Create() (err error) {
	err = u.setId()
	if err != nil {
		return
	}
	u.Password = fmt.Sprintf("%x", md5.Sum([]byte(u.Password)))
	body := map[string]interface{}{"etype": "user", "user": u}
	_, err = esCli.Index().Index(defaultIndex).Id(u.Id).BodyJson(body).Do(context.Background())
	return
}

// ClearUserCache clear all tokens from redis
func ClearUserCache() (dc int64, err error) {
	var keys []string
	keys, err = redisCli.Keys(userPrefix + "*").Result()
	if err != nil {
		return
	}
	if len(keys) > 0 {
		dc, err = redisCli.Del(keys...).Result()
		if err != nil {
			return
		}
		logger.Printf("redis: %d keys is deleted.\n", dc)
	}
	return
}
