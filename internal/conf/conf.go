package conf

import (
	"flag"
	"log"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type raftCfg struct {
	LocalID string // 节点名称
	Addr    string // 地址
	/*
	* 1. 是否作为主节点启动；true时启动集群，否则不启动集群，被动等待leader把当前节点加入集群
	* 2. 参考ETCD集群启动配置；ETCD_INITIAL_CLUSTER_STATE="new"  初始集群状态，new为新建集群，exist为子节点
	 */
	State string
}
type app struct {
	HttpAddr string `json:"http_port"` // http 监听端口
}
type AppConf struct {
	App  app `json:"app" toml:"app"`
	Raft raftCfg
}

var (
	configFile string
	Conf       AppConf
)

func init() {
	flag.StringVar(&configFile, "f", "", "configFile file")
}

func Init() {
	if configFile == "" {
		panic("config file must assign")
	}
	file, err := filepath.Abs(configFile)
	if err != nil {
		panic(err)
	}
	_, err = toml.DecodeFile(file, &Conf)
	if err != nil {
		panic(err)
	}
	log.Printf("conf:%+v", Conf)
}
