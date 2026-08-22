// internal/modules/workspace/param.go 入参定义。
//
// Author: Charlie

package workspace

// ShortcutSaveParam 保存工作台个人快捷应用入参（整体替换）。
//
// Author: Charlie
type ShortcutSaveParam struct {
	ResourceIDs []string `json:"resource_ids" binding:"max=16,dive"`
}
