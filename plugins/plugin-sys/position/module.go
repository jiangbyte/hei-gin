package position

import "hei-gin/sdk/infra/db"

type module struct {
	service *service
}

var defaultModule = newModule()

func newModule() *module {
	repo := &repository{db: db.DB}
	svc := &service{repo: repo}
	return &module{service: svc}
}
