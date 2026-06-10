package notice

type NoticeVO struct {
	ID        string  `json:"id"`
	Title     string  `json:"title"`
	Category  string  `json:"category"`
	Type      string  `json:"type"`
	Summary   *string `json:"summary"`
	Content   *string `json:"content"`
	Cover     *string `json:"cover"`
	Level     string  `json:"level"`
	IsTop     string  `json:"is_top"`
	Status    string  `json:"status"`
	SortCode  int     `json:"sort_code"`
	Author    *string `json:"author"`
	PublishAt *string `json:"publish_at"`
	ExpireAt  *string `json:"expire_at"`
	CreatedAt string  `json:"created_at"`
	CreatedBy *string `json:"created_by"`
	UpdatedAt string  `json:"updated_at"`
	UpdatedBy *string `json:"updated_by"`
}

type NoticeLatestParam struct {
	Size int `json:"size" form:"size"`
}

type NoticePageParam struct {
	Current  int    `json:"current" form:"current"`
	Size     int    `json:"size" form:"size"`
	Keyword  string `json:"keyword" form:"keyword"`
	Category string `json:"category" form:"category"`
	Status   string `json:"status" form:"status"`
}
