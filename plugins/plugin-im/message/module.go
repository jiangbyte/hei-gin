package message

type Module struct {
	service *Service
}

var DefaultModule = NewModule()

func NewModule() *Module {
	repo := &repository{}
	svc := &Service{repo: repo}
	return &Module{service: svc}
}

func (m *Module) Service() *Service {
	return m.service
}
