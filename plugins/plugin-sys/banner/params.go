package banner

// BannerVO is the view object for a banner, used for create/modify requests and API responses.
type BannerVO struct {
	ID          string  `json:"id" validate:"required"`
	Title       string  `json:"title" validate:"required"`
	Image       string  `json:"image" validate:"required"`
	URL         *string `json:"url"`
	LinkType    string  `json:"link_type" validate:"oneof=URL ROUTE"`
	Summary     *string `json:"summary"`
	Description *string `json:"description"`
	Category    string  `json:"category" validate:"required"`
	Type        string  `json:"type" validate:"required"`
	Position    string  `json:"position" validate:"required"`
	SortCode    int     `json:"sort_code"`
	ViewCount   int     `json:"view_count"`
	ClickCount  int     `json:"click_count"`
	CreatedAt   *string `json:"created_at"`
	CreatedBy   *string `json:"created_by"`
	UpdatedAt   *string `json:"updated_at"`
	UpdatedBy   *string `json:"updated_by"`
}

// BannerPageParam holds pagination parameters for the banner page query.
type BannerPageParam struct {
	Current int `json:"current" form:"current"`
	Size    int `json:"size" form:"size"`
}

func (p *BannerPageParam) GetCurrent() int { return p.Current }
func (p *BannerPageParam) GetSize() int    { return p.Size }
