package setting

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type config struct {
	Version   string `yaml:"version"`
	Databases struct {
		Mysql struct {
			Host     string `yaml:"host"`
			Port     int    `yaml:"port"`
			User     string `yaml:"user"`
			Password string `yaml:"password"`
			Db       string `yaml:"db"`
		}
		Redis struct {
			Addr     string `yaml:"addr"`
			Password string `yaml:"password"`
			Db       int    `yaml:"db"`
		}
		Elastic struct {
			Hosts []string `yaml:"hosts"`
		}
	}
}

const DefaultConf = "conf/app.yml"

var (
	AppVer string

	cf         *config
	CustomConf = DefaultConf
	HTTPPort   = "9090"
	HTTPAddr   = "localhost"

	RedisAddr = "localhost:6379"
	RedisPass = ""
	RedisDB   = 0

	EsHosts = []string{"http://127.0.0.1:9200"}

	MysqlHost     = "localhost"
	MysqlPort     = 3306
	MysqlUser     = "root"
	MysqlPassword = ""
	MysqlDB       = ""
)

func NewContext() {
	yf, err := ioutil.ReadFile(CustomConf)
	if err != nil {
		panic(fmt.Sprintf("get conf yaml file with err #%v", err))
	}
	cf = new(config)
	err = yaml.Unmarshal(yf, cf)
	if err != nil {
		panic(fmt.Sprintf("yaml.Unmarshal: %v", err))
	}

	// redis
	RedisAddr = cf.Databases.Redis.Addr
	RedisPass = cf.Databases.Redis.Password
	RedisDB = cf.Databases.Redis.Db

	// elasticsearch
	EsHosts = cf.Databases.Elastic.Hosts

	// mysql
	MysqlHost = cf.Databases.Mysql.Host
	MysqlPort = cf.Databases.Mysql.Port
	MysqlUser = cf.Databases.Mysql.User
	MysqlPassword = cf.Databases.Mysql.Password
	MysqlDB = cf.Databases.Mysql.Db
}
