package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/cyrushiker/pigo/models"
)

func userLogin(c *gin.Context) {
	var u models.UserVo
	if c.ShouldBindJSON(&u) != nil {
		c.JSON(http.StatusForbidden, "登录参数无效")
		return
	}
	tid, err := u.Login()
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}
	c.String(http.StatusOK, tid)
}

func userCreate(c *gin.Context) {
	var u models.User
	if c.ShouldBindJSON(&u) != nil {
		c.JSON(http.StatusForbidden, "用户信息无效")
		return
	}
	err := u.Create()
	if err != nil {
		c.String(http.StatusForbidden, err.Error())
		return
	}
	c.String(http.StatusOK, u.Id)
}

func userCacheClear(c *gin.Context) {
	dc, err := models.ClearUserCache()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.String(http.StatusOK, fmt.Sprintf("%v login user clear", dc))
}

func RegUser(user *gin.RouterGroup) {
	user.POST("/login", userLogin)
	user.POST("/clear", userCacheClear)
	user.POST("/create", userCreate)
}
