package notice

import "hei-gin/sdk/infra/db"

type Module struct {
	service *Service
}

var DefaultModule = NewModule()

func NewModule() *Module {
	repo := &repository{db: db.DB}
	svc := &Service{repo: repo}
	return &Module{service: svc}
}

func (m *Module) Service() *Service {
	return m.service
}
