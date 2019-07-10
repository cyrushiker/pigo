package routes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"

	"github.com/cyrushiker/pigo/models"
	"github.com/cyrushiker/pigo/pkg/setting"
)

func GlobalInit() {
	setting.NewContext()

	// init db clients
	models.NewRedisCli()
	models.NewEsCli()
	models.NewGormDB()

	// set captcha store to redis
	captcha.SetCustomStore(models.NewCaptchaStore(15 * time.Minute))
}

func captchaCode(c *gin.Context) {
	cid := captcha.New()
	c.Header("Captcha-Id", cid)
	captcha.WriteImage(c.Writer, cid, 160, 60)
}

func captchaVerify(c *gin.Context) {
	cid := c.Params.ByName("cid")
	cval := c.Params.ByName("cval")
	verify := captcha.VerifyString(cid, cval)
	c.String(http.StatusOK, fmt.Sprintf("%t", verify))
}

func RegHome(home *gin.RouterGroup) {
	home.GET("/", func(c *gin.Context) {
		data := map[string]interface{}{
			"author": fmt.Sprintf("Cyrushiker"),
			"email":  "cyrushiker@outlook.com",
		}
		c.AsciiJSON(http.StatusOK, data)
	})

	home.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	home.GET("/captcha", captchaCode)
	home.POST("/captcha/:cid/:cval", captchaVerify)
}
