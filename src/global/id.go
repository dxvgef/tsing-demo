package global

import "github.com/bwmarrin/snowflake"

// IDnode ID节点实例
var IDnode *snowflake.Node

// SetIDnode 设置ID节点
func SetIDnode() (err error) {
	if LocalConfig.Snowflake.Epoch > 0 {
		snowflake.Epoch = LocalConfig.Snowflake.Epoch
	}
	IDnode, err = snowflake.NewNode(LocalConfig.Snowflake.Node)
	return
}
