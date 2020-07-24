package global

import (
	"math/rand"
	"time"

	"github.com/bwmarrin/snowflake"
)

// IDnode ID节点实例
var IDnode *snowflake.Node

// SetIDnode 设置ID节点
func SetIDnode() (err error) {
	snowflake.Epoch = time.Now().Unix()
	rand.Seed(rand.Int63n(time.Now().UnixNano()))
	node := 0 + rand.Int63n(1023)
	IDnode, err = snowflake.NewNode(node)
	return
}
