// Package idgen 基于雪花算法生成分布式字符串 ID。
package idgen

import (
	"fmt"
	"sync"

	"github.com/bwmarrin/snowflake"
)

var (
	node    *snowflake.Node
	once    sync.Once
	initErr error
)

// Init 用 workerID 与 datacenterID 初始化雪花节点（仅生效一次）。
func Init(workerID, datacenterID int64) error {
	once.Do(func() {
		// 将 datacenter 折入 worker 高位，便于简单部署
		id := workerID + datacenterID*32
		node, initErr = snowflake.NewNode(id)
	})
	return initErr
}

// Next 生成下一个雪花 ID 字符串；未 Init 时 panic。
func Next() string {
	if node == nil {
		panic("idgen not initialized")
	}
	return fmt.Sprintf("%d", node.Generate().Int64())
}
