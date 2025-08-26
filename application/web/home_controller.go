package web

import (
	"github.com/gin-gonic/gin"
)

func HomeController(c *gin.Context) {
	c.String(200, "TODO flights home controller")
}
