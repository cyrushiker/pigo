package models

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"time"

	"github.com/olivere/elastic/v7"
)

const userPrefix = "pigo:user:"

type User struct {
	Id       string `json:"id" dt:"keyword"`
	Name     string `json:"name" dt:"keyword"`
	Password string `json:"password" dt:"keyword"`
}

type user struct {
	Etype string `json:"etype"`
	U     *User  `json:"user"`
}

type UserVo struct {
	User
	TokenId string `json:"tokenId"`
}

func (u *User) String() string {
	// remove Password from u before cache
	u.Password = ""
	uj, err := json.Marshal(*u)
	if err != nil {
		panic(fmt.Sprintf("json marshal with error: #%v", err))
	}
	return string(uj)
}

// UserLogin handle the login and return a tokenId
func (u *UserVo) Login() (string, error) {
	logger.Printf("Login user #%s", u.Name)

	query := elastic.NewBoolQuery().Must(
		elastic.NewTermQuery("etype", "user"),
		elastic.NewTermQuery("user.name", u.Name),
		elastic.NewTermQuery("user.password", fmt.Sprintf("%x", md5.Sum([]byte(u.Password)))),
	)
	sr, err := esCli.Search().Index(defaultIndex).Query(query).Do(context.Background())
	if err != nil {
		return "", err
	}
	if sr.Hits.TotalHits.Value < 1 {
		return "", fmt.Errorf("username or userpass is not correct.")
	}
	tokenid := userPrefix + getUUID("token")
	u.TokenId = tokenid
	_, err = redisCli.Set(tokenid, u.String(), 2*time.Hour).Result()
	if err != nil {
		return "", err
	}
	return tokenid, nil
}

func (u *UserVo) Logout() error {
	return nil
}

func (u *User) Create() (err error) {
	u.Id = getUUID("user")
	u.Password = fmt.Sprintf("%x", md5.Sum([]byte(u.Password)))
	body := user{Etype: "user", U: u}
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
