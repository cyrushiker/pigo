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
