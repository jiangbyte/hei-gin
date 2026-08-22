// internal/framework/platform/audit/labels.go 审计展示文案（对齐 hei-boot AuditLabelCatalog 子集）。
//
// Author: Charlie
package audit

import (
	"strings"
)

// ModuleLabel 由 resource_type 推导模块展示名。
func ModuleLabel(resourceType string) string {
	key := strings.ToLower(strings.TrimSpace(resourceType))
	switch key {
	case "auth", "account":
		return "认证 - 账号"
	case "auth_session":
		return "认证 - 会话"
	case "iam_account":
		return "权限 - 账号"
	case "iam_role":
		return "权限 - 角色"
	case "iam_dept":
		return "权限 - 部门"
	case "iam_group":
		return "权限 - 用户组"
	case "iam_position":
		return "权限 - 岗位"
	case "iam_resource", "resources":
		return "权限 - 资源"
	case "sys_notice":
		return "系统 - 消息"
	case "sys_banner":
		return "系统 - 展示图"
	case "sys_file":
		return "系统 - 文件"
	case "sys_config":
		return "系统 - 配置"
	case "sys_dict":
		return "系统 - 字典"
	case "sys_feedback":
		return "系统 - 反馈"
	case "profile_center":
		return "个人中心"
	case "real_name_case":
		return "实名认证 - 工单"
	case "profile_identity":
		return "实名认证 - 身份"
	case "workspace_shortcut":
		return "工作台 - 快捷应用"
	default:
		if strings.HasPrefix(key, "biz_") {
			return "业务 - " + key[4:]
		}
		if strings.HasPrefix(key, "sys_") {
			return "系统 - " + key[4:]
		}
		if strings.HasPrefix(key, "iam_") {
			return "权限 - " + key[4:]
		}
		if key == "" {
			return "系统"
		}
		return resourceType
	}
}

func entityShortName(resourceType string) string {
	module := ModuleLabel(resourceType)
	if idx := strings.Index(module, " - "); idx >= 0 {
		return module[idx+3:]
	}
	return module
}

// ActionName 推导操作展示名。
func ActionName(resourceType, action, explicit string) string {
	if strings.TrimSpace(explicit) != "" {
		return strings.TrimSpace(explicit)
	}
	act := normalizeAction(action)
	short := entityShortName(resourceType)
	switch act {
	case "create":
		return "创建" + short
	case "update":
		return "更新" + short
	case "delete":
		return "删除" + short
	case "login":
		return "登录"
	case "logout":
		return "退出登录"
	case "register":
		return "注册"
	case "upload":
		return "上传文件"
	case "publish":
		return "发布" + short
	case "read", "read_all", "read-all":
		return "阅读" + short
	case "submit":
		return "提交" + short
	case "batch-save", "batch_save":
		return "批量保存" + short
	default:
		if act != "" {
			return act
		}
		return "操作"
	}
}

// ActionType 推导操作类型枚举。
func ActionType(action, explicit string) string {
	if strings.TrimSpace(explicit) != "" {
		return strings.ToUpper(strings.TrimSpace(explicit))
	}
	act := normalizeAction(action)
	switch act {
	case "create":
		return "CREATE"
	case "update":
		return "UPDATE"
	case "delete":
		return "DELETE"
	case "login":
		return "LOGIN"
	case "logout":
		return "LOGOUT"
	case "register":
		return "CREATE"
	case "upload":
		return "CREATE"
	case "export":
		return "EXPORT"
	case "read", "read_all", "read-all":
		return "QUERY"
	default:
		return "OTHER"
	}
}

func normalizeAction(action string) string {
	return strings.ToLower(strings.TrimSpace(strings.ReplaceAll(action, "_", "-")))
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// EnrichActivityLabels 回填审计日志展示字段（DB 为空时按 catalog 推导）。
func EnrichActivityLabels(module, resourceType, action string, actionName, actionType, moduleLabel **string) {
	rt := resourceType
	if rt == "" {
		rt = module
	}
	if moduleLabel != nil && (*moduleLabel == nil || **moduleLabel == "") {
		*moduleLabel = strPtr(ModuleLabel(rt))
	}
	if actionName != nil && (*actionName == nil || **actionName == "") {
		*actionName = strPtr(ActionName(rt, action, ""))
	}
	if actionType != nil && (*actionType == nil || **actionType == "") {
		*actionType = strPtr(ActionType(action, ""))
	}
}
