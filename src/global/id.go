package global

import "github.com/bwmarrin/snowflake"

// IDnode ID节点实例
var IDnode *snowflake.Node

// SetIDnode 设置ID节点
func SetIDnode() (err error) {
	IDnode, err = snowflake.NewNode(0)
	return
}
