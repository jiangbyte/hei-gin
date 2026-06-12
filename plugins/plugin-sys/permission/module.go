package permission

import (
	"hei-gin/sdk/infra/db"
)

var DefaultModule = NewModule()

type Module struct {
	service *Service
}

func NewModule() *Module {
	repo := &repository{rdb: db.Redis}
	svc := &Service{repo: repo}
	return &Module{service: svc}
}

func (m *Module) Service() *Service {
	return m.service
}
