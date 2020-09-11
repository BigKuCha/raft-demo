package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/raft"
	"study/raft/pkgs/fsm"
)

var (
	localID string // 节点名称
	address string // 地址
	/*
	* 1. 是否作为主节点启动；true时启动集群，否则不启动集群，被动等待leader把当前节点加入集群
	* 2. 参考ETCD集群启动配置；ETCD_INITIAL_CLUSTER_STATE="new"  初始集群状态，new为新建集群，exist为子节点
	 */
	leaderNode bool
)

func init() {
	flag.StringVar(&localID, "id", "local", "localid")
	flag.StringVar(&address, "addr", "127.0.0.1:8880", "address")
	flag.BoolVar(&leaderNode, "l", false, "leader or not")
}

func main() {
	flag.Parse()
	cfg := raft.DefaultConfig()
	cfg.LocalID = raft.ServerID(localID)
	// 新建内存日志存储
	logstore := raft.NewInmemStore()
	address := address
	// 节点间通讯地址
	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		panic(err)
	}

	transport, err := raft.NewTCPTransport(address, tcpAddr, 3, 10*time.Second, os.Stderr)
	if err != nil {
		panic(err)
	}
	hlog := hclog.New(&hclog.LoggerOptions{})
	fss, err := raft.NewFileSnapshotStoreWithLogger("", 1, hlog)
	if err != nil {
		panic(err)
	}
	// 自定义有限状态机，接收日志变更事件
	mfsm := &fsm.MyFsm{LocalID: string(cfg.LocalID)}
	r, err := raft.NewRaft(cfg, mfsm, logstore, logstore, fss, transport)
	if err != nil {
		panic(err)
	}
	if leaderNode {
		/*
		* 集群节点。
		* 1. 集群启动时，检测所有server，投票数超过一半才能选举出leader，子节点数不足一半，集群启动失败(没有leader)
		* 2. 不配置所有servers时，通过r.AddVoter() 动态添加节点.
		* 3. 只配置当前节点，则当前节点为leader。动态添加节点时leader不变
		* 3. 集群搭建时尽量配置所有节点，集群重启时不用一个个添加子节点
		 */
		servers := []raft.Server{
			// 当前节点
			{
				ID:      cfg.LocalID,
				Address: transport.LocalAddr(),
			},
			// 其他节点
			{
				ID:      "local1",
				Address: "127.0.0.1:8881",
			},
			{
				ID:      "local2",
				Address: "127.0.0.1:8882",
			},
		}
		// 启动集群
		r.BootstrapCluster(raft.Configuration{Servers: servers})
	}

	// time.Sleep(4 * time.Second)
	// r.AddVoter("voter", "127.0.0.1:8881", 0, 0)
	// r.Apply()
	tk := time.NewTicker(time.Second * 4)
	for {
		select {
		case <-tk.C:
			if r.State() == raft.Leader {
				fmt.Println("apply log", r.Stats())
				r.Apply([]byte(time.Now().String()), time.Second)
			}
			// meta, _, _ := r.Snapshot().Open()
			// fmt.Println("meta", meta)
			// fmt.Println(r.State(), time.Now().String())
		}
	}
}
