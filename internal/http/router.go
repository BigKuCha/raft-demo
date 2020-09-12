package http

import (
	"github.com/gin-gonic/gin"
)

func registerRouter(e *gin.Engine) {
	{
		base := e.Group("")
		base.GET("", func(c *gin.Context) {
			c.JSON(200, "hello")
		})
	}
	{
		raft := e.Group("raft")
		raft.GET("stats", stats)
	}
}
