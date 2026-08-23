// internal/modules/iam/relation/model.go 数据模型。
//
// Author: Charlie

package relation

import (
	"time"

	"gorm.io/datatypes"
)

// Relation 映射 sys_iam_relation 授权关系边。
//
// Author: Charlie
type Relation struct {
	ID                 string         `gorm:"column:id;primaryKey;size:64" json:"id"`
	SubjectType        string         `gorm:"column:subject_type;size:32" json:"subject_type"`
	SubjectID          string         `gorm:"column:subject_id;size:64" json:"subject_id"`
	AccountType        string         `gorm:"column:account_type;size:32" json:"account_type"`
	RelationType       string         `gorm:"column:relation_type;size:64" json:"relation_type"`
	TargetType         string         `gorm:"column:target_type;size:32" json:"target_type"`
	TargetID           string         `gorm:"column:target_id;size:64" json:"target_id"`
	TargetKey          string         `gorm:"column:target_key;size:128" json:"target_key"`
	GrantMode          string         `gorm:"column:grant_mode;size:32" json:"grant_mode"`
	DataScope          string         `gorm:"column:data_scope;size:32" json:"data_scope"`
	CustomScopeDeptIDs datatypes.JSON `gorm:"column:custom_scope_dept_ids;type:json" json:"custom_scope_dept_ids"`
	IsPrimary          bool           `gorm:"column:is_primary" json:"is_primary"`
	Sort               int            `gorm:"column:sort" json:"sort"`
	Status             string         `gorm:"column:status;size:32" json:"status"`
	Description        *string        `gorm:"column:description" json:"description"`
	Reason             *string        `gorm:"column:reason" json:"reason"`
	ExpiredAt          *time.Time     `gorm:"column:expired_at" json:"expired_at"`
	Extra              datatypes.JSON `gorm:"column:extra;type:json" json:"extra"`
	CreatedAt          time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	CreatedBy          *string        `gorm:"column:created_by;size:64" json:"created_by"`
	UpdatedAt          time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	UpdatedBy          *string        `gorm:"column:updated_by;size:64" json:"updated_by"`
}

// TableName 返回表名。
func (Relation) TableName() string { return "sys_iam_relation" }

// 主体类型与关系类型常量。
const (
	SubjectAccount        = "ACCOUNT"
	SubjectRole           = "ROLE"
	SubjectGroup          = "GROUP"
	SubjectResource       = "RESOURCE"
	SubjectClientResource = "CLIENT_RESOURCE"

	RelAccountRole              = "ACCOUNT_ROLE"
	RelAccountDept              = "ACCOUNT_DEPT"
	RelAccountGroup             = "ACCOUNT_GROUP"
	RelSubjectResourceGrant     = "SUBJECT_RESOURCE_GRANT"
	RelSubjectClientResourceGrant = "SUBJECT_CLIENT_RESOURCE_GRANT"
	RelGroupRole                = "GROUP_ROLE"
	RelResourcePermission       = "RESOURCE_PERMISSION"
	RelClientResourcePermission = "CLIENT_RESOURCE_PERMISSION"
)

// 目标类型常量。
const (
	TargetAccount        = "ACCOUNT"
	TargetRole           = "ROLE"
	TargetDept           = "DEPT"
	TargetGroup          = "GROUP"
	TargetResource       = "RESOURCE"
	TargetClientResource = "CLIENT_RESOURCE"
	TargetPermission     = "PERMISSION"
)

// 授予模式常量。
const (
	GrantDirect  = "DIRECT"
	GrantCascade = "CASCADE"
)
