package global

import "github.com/bwmarrin/snowflake"

// IDnode ID节点实例
var IDnode *snowflake.Node

// SetIDnode 设置ID节点
func SetIDnode() (err error) {
	if Config.Snowflake.Epoch > 0 {
		snowflake.Epoch = Config.Snowflake.Epoch
	}
	IDnode, err = snowflake.NewNode(Config.Snowflake.Node)
	return
}
