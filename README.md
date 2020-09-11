## Raft
[hashicorp/raft](https://github.com/hashicorp/raft) 

## 启动集群
```
$ ./raft -id local -addr 127.0.0.1:8880 -l true # 启动主节点
$ ./raft -id local1 addr 127.0.0.1:8881         # 启动节点1
$ ./raft -id local2 addr 127.0.0.1:8882         # 启动节点2
```