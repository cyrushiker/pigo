package models

import (
	"encoding/json"
	// "fmt"
	// "os"
	"testing"
	"time"

	"github.com/dchest/captcha"
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
		"valueUnit6": null,
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
	for _, g := range tkeys {
		if d, ok := d[g.key]; ok {
			t.Logf("Before::%#v --- %T", d, d)
			rv, err := verifySwitch(d, g.key, g.kt, g.rr, g.unit)
			if err != nil {
				t.Errorf("After:: rv = %v ~~~ err = %v", rv, err)
			} else {
				t.Logf("After:: rv = %v ~~~ err = %v", rv, err)
			}
		}
	}
}
