package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/cyrushiker/pigo/models"
)

func userLogin(c *gin.Context) {
	var u models.UserVo
	if err := c.ShouldBindJSON(&u); err != nil {
		c.String(http.StatusForbidden, fmt.Sprintf("parse user info error: #%v", err))
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
	u := new(models.User)
	if err := c.ShouldBindJSON(u); err != nil {
		c.String(http.StatusForbidden, fmt.Sprintf("parse user info error: #%v", err))
		return
	}
	// todo: how to handle the zero value of time.Time
	err := u.Create()
	if err != nil {
		c.String(http.StatusForbidden, err.Error())
		return
	}
	c.String(http.StatusOK, u.Id)
}

func userList(c *gin.Context) {
	u := new(models.User)
	datas, err := u.List()
	if err != nil {
		c.String(http.StatusForbidden, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{"data": datas})
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
	user.GET("", userList)
}
