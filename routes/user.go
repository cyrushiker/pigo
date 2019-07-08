package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/cyrushiker/pigo/models"
)

func auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("tokenId")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "User not logged in")
			return
		}
		log.Printf("token: %s", token)
		// check tokenId from redis
		uv, err := models.CurrentUserWithToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		}
		c.Set("user", uv)
		c.Next()
	}
}

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
	err := u.Create()
	if err != nil {
		c.String(http.StatusForbidden, err.Error())
		return
	}
	c.String(http.StatusOK, u.Id)
}

func userList(c *gin.Context) {
	cu, e := c.Get("user")
	if e != true {
		c.String(http.StatusInternalServerError, "user not found in context")
		return
	}
	datas, err := cu.(*models.UserVo).List()
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
	user.GET("", auth(), userList)
}
