package home

type QuickActionVO struct {
	ID         string `json:"id"`
	ResourceID string `json:"resource_id"`
	ParentID   string `json:"parent_id"`
	Type       string `json:"type"`
	Name       string `json:"name"`
	Icon       string `json:"icon"`
	RoutePath  string `json:"route_path"`
	SortCode   int    `json:"sort_code"`
}

type HomeNotice struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Level     string `json:"level"`
	CreatedAt string `json:"created_at"`
}

type HomeStats struct {
	TotalUsers int `json:"total_users"`
}

type HomeVO struct {
	QuickActions       []QuickActionVO `json:"quick_actions"`
	AvailableResources []QuickActionVO `json:"available_resources"`
	Notices            []HomeNotice    `json:"notices"`
	Stats              HomeStats       `json:"stats"`
}

type AddQuickActionParam struct {
	ResourceID string `json:"resource_id"`
}

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
