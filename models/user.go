package models

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"time"

	"github.com/olivere/elastic/v7"
)

const (
	userPrefix = "pigo:user:"
)

func passMd5(p string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(p)))
}

type User struct {
	Id        string    `json:"id" dt:"keyword"`
	Name      string    `json:"name" dt:"keyword"`
	Password  string    `json:"password" dt:"keyword"`
	TelNumber string    `json:"telNumberm,omitempty" dt:"keyword"`
	Address   string    `json:"address,omitempty" dt:"text"`
	Birthday  time.Time `json:"birthday,omitempty" dt:"date"`
}

type user struct {
	Etype string `json:"etype"`
	U     *User  `json:"user"`
}

func (u *User) esTypeName() string {
	return "user"
}

func (u *User) checkExists() (bool, error) {
	query := elastic.NewBoolQuery().Must(
		elastic.NewTermQuery("etype", u.esTypeName()),
		elastic.NewTermQuery(u.esTypeName()+".name", u.Name),
	)
	sr, err := esCli.Search().Index(defaultIndex).Query(query).Do(context.Background())
	if err != nil {
		return true, err
	}
	if sr.Hits.TotalHits.Value > 0 {
		return true, nil
	}
	return false, nil
}

func (u *User) Create() (err error) {
	if has, err := u.checkExists(); err != nil {
		return err
	} else if has {
		return fmt.Errorf("User(%s) is exists.", u.Name)
	}
	u.Id = getUUID(u.esTypeName())
	u.Password = passMd5(u.Password)
	body := user{Etype: u.esTypeName(), U: u}
	_, err = esCli.Index().Index(defaultIndex).Id(u.Id).BodyJson(body).Do(context.Background())
	return
}

func (u *User) List() ([]*User, error) {
	fs := elastic.NewFetchSourceContext(true).Exclude(u.esTypeName() + ".password")
	query := elastic.NewTermQuery("etype", u.esTypeName())
	sr, err := esCli.Search(defaultIndex).FetchSourceContext(fs).Query(query).Do(context.Background())
	if err != nil {
		return nil, err
	}
	if sr.Hits.TotalHits.Value < 1 {
		return nil, fmt.Errorf("No user found")
	}
	users := make([]*User, 0)
	for _, hit := range sr.Hits.Hits {
		u := new(user)
		err := json.Unmarshal(hit.Source, u)
		if err != nil {
			return nil, fmt.Errorf("json unmarshal error: %v", err)
		}
		users = append(users, u.U)
	}
	return users, nil
}

type UserVo struct {
	*User
	TokenId string `json:"tokenId"`
}

func (u *UserVo) String() string {
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
	// search in es
	query := elastic.NewBoolQuery().Must(
		elastic.NewTermQuery("etype", u.esTypeName()),
		elastic.NewTermQuery(u.esTypeName()+".name", u.Name),
		elastic.NewTermQuery(u.esTypeName()+".password", passMd5(u.Password)),
	)
	sr, err := esCli.Search().Index(defaultIndex).Query(query).Do(context.Background())
	if err != nil {
		return "", err
	}
	if sr.Hits.TotalHits.Value < 1 {
		return "", fmt.Errorf("username or userpass is not correct.")
	}
	uu := new(user)
	err = json.Unmarshal(sr.Hits.Hits[0].Source, uu)
	if err != nil {
		return "", err
	}
	u.User = uu.U
	tokenid := getUUID("token")
	u.TokenId = tokenid
	logger.Println(u.String())
	_, err = redisCli.Set(userPrefix+tokenid, u.String(), 2*time.Hour).Result()
	if err != nil {
		return "", err
	}
	return tokenid, nil
}

func (u *UserVo) Logout() error {
	return nil
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
