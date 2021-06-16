package global

import (
	"github.com/bwmarrin/snowflake"
	"github.com/rs/zerolog/log"
)

// snowflake ID节点实例
var SnowflakeNode *snowflake.Node

// 设置snowflake节点
func SetSnowflake() (err error) {
	SnowflakeNode, err = snowflake.NewNode(0)
	log.Err(err).Caller().Msg("设置snowflake节点失败")
	return
}
