// Package scope 为 biz 示例模块提供数据范围辅助（对齐 fastapi data_scope.py）。
//
// Author: Charlie
package scope

import (
	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/core/security/datascope"
)

// DefaultOwnerDeptID 新建行默认归属部门（会话首个 dept_id）。
func DefaultOwnerDeptID(sess *security.SessionPayload) *string {
	if sess == nil || len(sess.DeptIDs) == 0 {
		return nil
	}
	id := sess.DeptIDs[0]
	return &id
}

// Assert 写操作数据范围校验。
func Assert(sess *security.SessionPayload, ownerDeptID, createdBy *string) error {
	dept, acct := "", ""
	if ownerDeptID != nil {
		dept = *ownerDeptID
	}
	if createdBy != nil {
		acct = *createdBy
	}
	return datascope.Assert(sess, dept, acct)
}
