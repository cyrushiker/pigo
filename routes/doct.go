package routes

import (
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

func RegDoct(doct *gin.RouterGroup) {
	doct.POST("/add", addDoct)
}
