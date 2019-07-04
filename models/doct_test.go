package models

import (
	"fmt"
	"os"
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
	out, _ := os.Create(fmt.Sprintf("/tmp/%s.png", cid))
	captcha.WriteImage(out, cid, 320, 80)

	vr := captcha.VerifyString(cid, "121111")
	t.Logf("verify result: %v", vr)
}