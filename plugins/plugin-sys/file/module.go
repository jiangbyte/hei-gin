package file

import "hei-gin/sdk/infra/db"

type module struct {
	service *service
}

var defaultModule = newModule()

func newModule() *module {
	repo := &repository{db: db.DB, rdb: db.Redis}
	svc := &service{repo: repo}
	return &module{service: svc}
}
