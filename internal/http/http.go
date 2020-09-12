package http

import (
	"flag"
	"fmt"

	"github.com/gin-gonic/gin"
	"study/raft/internal/conf"
)

type Server struct {
	engine *gin.Engine
}

func NewHttpServer() *Server {
	fmt.Printf("conf:%+v", conf.Conf)
	return &Server{
		engine: newEngine(),
	}
}

func (h *Server) Run() {
	if !flag.Parsed() {
		fmt.Println("******")
		flag.Parse()
	}
	err := h.engine.Run(conf.Conf.App.HttpAddr)
	if err != nil {
		panic(err)
	}
}

func newEngine() *gin.Engine {
	e := gin.Default()
	registerRouter(e)
	return e
}
