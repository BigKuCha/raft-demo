package main

import (
	"flag"

	"study/raft/internal/conf"
	"study/raft/internal/http"
	"study/raft/internal/store"
)

func main() {
	flag.Parse()
	// 配置初始化
	conf.Init()
	// raft初始化
	store.Init()
	// http接口
	httpSrv := http.NewHttpServer()
	httpSrv.Run()
}
