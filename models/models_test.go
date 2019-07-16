package models

import (
	"context"
	"encoding/json"
	// "os"
	"testing"
	"time"

	"github.com/dchest/captcha"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestDoctMapping(t *testing.T) {
	d := new(Doct)
	t.Log(mapping(d))

	NewRedisCli()
	cs := NewCaptchaStore(15 * time.Minute)
	captcha.SetCustomStore(cs)

	cid := captcha.New()
	// out, _ := os.Create(fmt.Sprintf("/tmp/%s.png", cid))
	// captcha.WriteImage(out, cid, 320, 80)

	vr := captcha.VerifyString(cid, "121111")
	t.Logf("verify result: %v", vr)
	t.Logf("time now: %s", time.Now().Format("2006-01-02 15:04:05"))
	u := new(User)
	us, _ := json.Marshal(u)
	t.Log(string(us))
}

func TestDoctAdd(t *testing.T) {
	data := `{
		"a": 1,
		"b": {
			"b1": "11"
		},
		"c": [{
			"c1": 11,
			"c2": "12",
			"c3": null
		},{
			"c1": 12,
			"c2": "22"
		}]
	}`
	d := make(map[string]interface{})
	json.Unmarshal([]byte(data), &d)
	t.Logf("#%v", d)
	t.Logf(`d["b"] type is: #%T`, d["b"])
	t.Logf(`d["c"] type is: #%T`, d["c"])
	for _, g := range []string{"a", "b", "c"} {
		if d, ok := d[g]; ok {
			if _, ok := d.([]interface{}); ok {
				t.Logf(`d["%s"] is an array`, g)
			} else if _, ok := d.(map[string]interface{}); ok {
				t.Logf(`d["%s"] is a map`, g)
			} else {
				t.Logf(`d["%s"] is not a map or an array`, g)
			}
		}
	}
}

func TestDoctVerify(t *testing.T) {
	data := `{
		"stringKey1": "   脑残病  ",
		"stringKey2": {"value": 11, "unit": "ml"},
		"arrayStringKey1": ["你", "是", "谁", "?"],
		"multipleKey1": ["你", "是", "谁", "?"],
		"multipleKey2": [1, 3, 11],
		"numberKey1": 111,
		"dateKey1": "2017-01-09 11:12",
		"dateKey2": "2017-01-09 11:12:11",
		"dateKey3": "2017-01-09",
		"dateKey4": "2017-01-09 ",
		"dateKey5": "13452231235123",
		"valueUnit1": {"value": 11, "unit": "mg/ml", "origin": ">=11"},
		"valueUnit2": "123",
		"valueUnit3": "133 mg/gg",
		"valueUnit4": "23 mg/gg gt",
		"valueUnit5": {"value": "14", "unit": "mg/ml", "origin": "null"},
		"valueUnit6": "null",
		"valueUnit7": 444,
		"regExpKey1": "123456"
	}`
	d := make(map[string]interface{})
	err := json.Unmarshal([]byte(data), &d)
	if err != nil {
		t.Fatal(err)
	}

	tkeys := []struct {
		key, kt, rr, unit string
	}{
		{"stringKey1", "String", "", ""},
		{"stringKey2", "Object", "", ""},
		{"arrayStringKey1", "ArrayString", "", ""},
		{"multipleKey1", "Multiple", "", ""},
		{"multipleKey2", "Multiple", "", ""},
		{"numberKey1", "Number", "", ""},
		{"dateKey1", "Date", "", ""},
		{"dateKey2", "Date", "", ""},
		{"dateKey3", "Date", "", ""},
		{"dateKey4", "Date", "", ""},
		{"dateKey5", "Date", "", ""},
		{"valueUnit1", "ValueUnit", "", ""},
		{"valueUnit2", "ValueUnit", "", ""},
		{"valueUnit3", "ValueUnit", "", ""},
		{"valueUnit4", "ValueUnit", "", ""},
		{"valueUnit5", "ValueUnit", "", ""},
		{"valueUnit6", "ValueUnit", "", ""},
		{"valueUnit7", "ValueUnit", "", ""},
		{"regExpKey1", "RegExp", `/^\d+$/`, ""},
	}
	indoc := make(map[string]interface{})
	for _, g := range tkeys {
		if d, ok := d[g.key]; ok {
			t.Logf("Before::%#v --- %T", d, d)
			rv, err := verifySwitch(d, g.key, g.kt, g.rr, g.unit)
			if err != nil {
				t.Errorf("After:: rv = %v ~~~ err = %v", rv, err)
			} else {
				indoc[g.key] = rv
				t.Logf("After:: rv = %v ~~~ err = %v", rv, err)
			}
		}
	}

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		t.Fatal(err)
	}
	collection := client.Database("test").Collection("verify")
	insertResult, err := collection.InsertOne(context.TODO(), indoc)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Inserted a single document: ", insertResult.InsertedID)

	result := make(map[string]interface{})
	err = collection.FindOne(context.TODO(), bson.M{}).Decode(&result)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Found a single document: %#v\n", result)
	for k, v := range result {
		t.Logf("Search: Key=%s *** Val=%v *** VType=%T ##", k, v, v)
	}

	indoc2 := make(map[string]interface{})
	collection2 := client.Database("test").Collection("verify2")
	for _, k := range []string{
		"stringKey1",
		"arrayStringKey1",
		"multipleKey1",
		"numberKey1",
		"dateKey1",
		"valueUnit1",
		"regExpKey1",
	} {
		indoc2[k] = result[k]
	}
	collection2.InsertOne(context.TODO(), indoc2)
}

func TestGorm(t *testing.T) {
	NewGormDB()
	defer db.Close()
	// err := db.DropTableIfExists(&Atom{}).CreateTable(&Atom{}).Error
	// if err != nil {
	// 	t.Log(err)
	// }

	a := Atom{Key: "key4", Name: "测试Key3", Keyword: "关键字 第一个 测试字段", Level: 1}
	if db.NewRecord(a) {
		dbn := db.Create(&a)
		if dbn.Error != nil {
			t.Fatal(dbn.Error)
		}
		f, _ := json.MarshalIndent(a, "", "  ")
		t.Log(string(f))
	}

	// mkr := new(Atom)
	// if err := db.First(mkr).Error; err == nil {
	// 	// t.Log(mkr)
	// 	mkrj, _ := json.MarshalIndent(mkr, "", "  ")
	// 	t.Log(string(mkrj))
	// }
}

func TestHP(t *testing.T) {
	girls, _ := Extract("https://www.meizitu.com")

	t.Logf("%#v", girls)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := DownloadPics(ctx, "/tmp/meizitu/", girls, 10)
	t.Logf("%v", err)
}

func TestCtx(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	r, err := longOpDo(ctx)
	if err != nil {
		t.Fatal(err, r)
	}
}
