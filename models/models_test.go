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
