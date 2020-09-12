package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"study/raft/internal/store"
)

func stats(c *gin.Context) {
	stats := store.Raft.Stats()
	stats["leader"] = string(store.Raft.Leader())
	stats["http_port"] = store.Raft.Meta["http_port"]
	c.JSON(http.StatusOK, stats)
}
