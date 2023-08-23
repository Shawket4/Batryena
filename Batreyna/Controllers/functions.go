package Controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ReturnErr(c *gin.Context, err error) {
	log.Println(err.Error())
	c.JSON(http.StatusBadRequest, err)
}
