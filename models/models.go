package models

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/olivere/elastic/v7"

	"github.com/cyrushiker/pigo/pkg/setting"
)

var (
	redisCli *redis.Client
	esCli    *elastic.Client
	db       *gorm.DB

	indexTypes []esType
	tables     []interface{}
)

var MyTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
	return "pg_" + defaultTableName
}

const redisNil = redis.Nil

const defaultIndex = ".pigo"

func init() {
	indexTypes = append(indexTypes, new(Doct), new(User))
	tables = append(tables, new(Atom), new(Molecule), new(Compound))
}

func CreateIndex(flag bool) error {
	// check exists
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := esCli.IndexExists(defaultIndex).Do(ctx)
	if err != nil {
		return fmt.Errorf("Exists index with error: #%v", err)
	}
	if res {
		// exists
		if flag {
			// delete index
			res, err := esCli.DeleteIndex(defaultIndex).Do(ctx)
			if err != nil {
				return fmt.Errorf("Delete index with error: #%v", err)
			}
			if res.Acknowledged {
				return createIndex(ctx)
			} else {
				return fmt.Errorf("Delete index with result: #%v", res)
			}
		} else {
			return fmt.Errorf("Index is already exists\n")
		}
	} else {
		return createIndex(ctx)
	}
}

func createIndex(ctx context.Context) error {
	// create or reset index mapping
	indexMapping := make(map[string]interface{})
	properties := make(map[string]interface{})
	properties["etype"] = prop("keyword", "")
	for _, it := range indexTypes {
		properties[it.esTypeName()] = mapping(it)
	}
	indexMapping["mappings"] = map[string]interface{}{"properties": properties}

	createIndex, err := esCli.CreateIndex(defaultIndex).BodyJson(indexMapping).Do(ctx)
	if err != nil {
		return fmt.Errorf("Init index mapping with error: #%v", err)
	}
	if createIndex == nil {
		return fmt.Errorf("expected put mapping response; got: %v", createIndex)
	}
	logger.Printf("Create index (%s) successfully", defaultIndex)
	return nil
}

// NewRedisCli new redis cli in global scope
func NewRedisCli() {
	redisCli = redis.NewClient(&redis.Options{
		Addr:     setting.RedisAddr,
		Password: setting.RedisPass, // no password set
		DB:       setting.RedisDB,   // use default DB
	})
}

// NewEsCli new es cli in global scope
func NewEsCli() {
	var err error
	logger.Printf("Connect to ES: %v", setting.EsHosts)
	esCli, err = elastic.NewClient(
		elastic.SetURL(setting.EsHosts...),
		elastic.SetSniff(false),
		elastic.SetErrorLog(logger),
	)
	if err != nil {
		panic(fmt.Sprintf("New elastic client error: #%v", err))
	}
}

func NewGormDB() {
	var err error
	gorm.DefaultTableNameHandler = MyTableNameHandler
	db, err = gorm.Open("mysql", "root:root@tcp(localhost:3306)/pigo?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}
	db.LogMode(true)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	db.DB().SetConnMaxLifetime(time.Hour)
}

func CreateTables(drop bool) error {
	if drop {
		return db.DropTableIfExists(tables...).CreateTable(tables...).Error
	}
	return db.CreateTable(tables...).Error
}

// esType declare the type name
type esType interface {
	esTypeName() string
}

func mapping(i esType) map[string]interface{} {
	fmap := make(map[string]interface{})
	v := reflect.TypeOf(i).Elem()
	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Field(i)
		tag := fieldInfo.Tag
		dt := tag.Get("dt")
		jn := tag.Get("json")
		ns := strings.Split(jn, ",")
		if len(ns) > 0 && dt != "" {
			p := prop(dt, "")
			if p != nil {
				fmap[ns[0]] = p
			}
		}
	}
	return map[string]interface{}{"properties": fmap}
}

func prop(dt, analyzer string) map[string]interface{} {
	switch dt {
	case "date":
		return map[string]interface{}{
			"type":   "date",
			"format": "strict_date_optional_time||epoch_millis",
		}
	case "keyword":
		return map[string]interface{}{
			"type": "keyword",
		}
	case "text":
		return map[string]interface{}{
			"type": "text",
		}
	default:
		return nil
	}
}

func getUUID(tname string) string {
	id, err := uuid.NewRandom()
	if err != nil {
		panic(fmt.Sprintf("google uuid with error: #%v", err))
	}
	if tname == "" {
		return id.String()
	}
	return fmt.Sprintf("%s_%s", tname, id)
}

// cache captcha to redis
type captchaStore struct {
	m          *redis.Client
	expiration time.Duration
}

func NewCaptchaStore(expiration time.Duration) *captchaStore {
	s := new(captchaStore)
	s.expiration = expiration
	s.m = redisCli
	return s
}

func (s *captchaStore) Id(id string) string {
	return fmt.Sprintf("pigo:captcha:%s", id)
}

func (s *captchaStore) Set(id string, digits []byte) {
	_, err := s.m.Set(s.Id(id), digits, s.expiration).Result()
	if err != nil {
		panic(fmt.Sprintf("captchaStore Set: %v", err))
	}
}

func (s *captchaStore) Get(id string, clear bool) (digits []byte) {
	digits, err := s.m.Get(s.Id(id)).Bytes()
	if err != nil {
		return
	}
	if clear {
		defer s.Del(id)
	}
	return
}

func (s *captchaStore) Del(id string) {
	_, err := s.m.Del(s.Id(id)).Result()
	if err != nil {
		panic(fmt.Sprintf("captchaStore Del: %v", err))
	}
}
