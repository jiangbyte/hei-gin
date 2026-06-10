package home

// SysQuickActionToQuickActionVO 将 home.SysQuickAction 映射到 home.QuickActionVO
func SysQuickActionToQuickActionVO(src *SysQuickAction) *QuickActionVO {
	if src == nil {
		return nil
	}

	dst := &QuickActionVO{}

	dst.ID = src.ID
	dst.ResourceID = src.ResourceID
	dst.SortCode = src.SortCode

	return dst
}

// QuickActionVOToSysQuickAction 将 home.QuickActionVO 映射到 home.SysQuickAction
func QuickActionVOToSysQuickAction(src *QuickActionVO) *SysQuickAction {
	if src == nil {
		return nil
	}

	dst := &SysQuickAction{}

	dst.ID = src.ID
	dst.ResourceID = src.ResourceID
	dst.SortCode = src.SortCode

	return dst
}
