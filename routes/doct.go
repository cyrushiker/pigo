package routes

import (
	"github.com/cyrushiker/pigo/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func addDoct(c *gin.Context) {
	doctBody := make(map[string]interface{})
	if err := c.ShouldBindJSON(&doctBody); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	log.Printf("#%v", doctBody)
	c.String(http.StatusOK, `the body should be map`)
}

func addMetaKey(c *gin.Context) {
	a := new(models.Atom)
	if err := c.ShouldBindJSON(&a); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if err := a.Save(); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
}

func RegDoct(doct *gin.RouterGroup) {
	doct.POST("/add", addDoct)
	doct.PUT("/metakey", addMetaKey)
}
