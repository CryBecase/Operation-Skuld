package xmiddleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Cors(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "*")
	if c.Request.Method == http.MethodOptions {
		c.AbortWithStatus(200)
	}
}
