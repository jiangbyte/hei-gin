// internal/modules/iam/role/result.go 出参定义。
//
// Author: Charlie

package role

import (
	"hei-gin/internal/modules/iam/client"
	"hei-gin/internal/modules/iam/relation"
	"hei-gin/internal/modules/iam/resource"
)

// OwnResourceResult è§’è‰²å·²æ‹¥æœ‰ç®¡ç†ç«¯èµ„æºæŽˆæƒç»“æžœã€‚
//
// Author: Charlie
type OwnResourceResult struct {
	ID            string                       `json:"id"`
	Modules       []resource.GrantModule       `json:"modules"`
	GrantInfoList []relation.ResourceGrantInfo `json:"grant_info_list"`
}

// OwnClientResourceResult è§’è‰²å·²æ‹¥æœ‰å®¢æˆ·ç«¯èµ„æºæŽˆæƒç»“æžœã€‚
//
// Author: Charlie
type OwnClientResourceResult struct {
	ID            string                       `json:"id"`
	Modules       []client.GrantModule         `json:"modules"`
	GrantInfoList []relation.ResourceGrantInfo `json:"grant_info_list"`
}
