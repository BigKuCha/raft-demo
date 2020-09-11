package fsm

import (
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/hashicorp/go-msgpack/codec"
	"github.com/hashicorp/raft"
)

type MyFsm struct {
	sync.Mutex
	LocalID string
}

func (m *MyFsm) Apply(log *raft.Log) interface{} {
	fmt.Printf("==== %s | %s | %s \n", m.LocalID, log.Data, time.Now().String())
	return nil
}

func (m *MyFsm) Snapshot() (raft.FSMSnapshot, error) {
	return &raft.MockSnapshot{}, nil
}

func (m *MyFsm) Restore(irc io.ReadCloser) error {
	m.Lock()
	defer m.Unlock()
	defer irc.Close()
	hd := codec.MsgpackHandle{}
	dec := codec.NewDecoder(irc, &hd)

	return dec.Decode("")
}
