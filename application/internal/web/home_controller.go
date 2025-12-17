package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func HomeController(c *gin.Context) {
	c.HTML(http.StatusOK, "home.tmpl", gin.H{
		"title": "Main website",
	})
}
