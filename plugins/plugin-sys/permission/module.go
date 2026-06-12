package permission

import (
	"hei-gin/sdk/db"
)

type module struct {
	service *service
}

var defaultModule = newModule()

func newModule() *module {
	repo := &repository{rdb: db.Redis}
	svc := &service{repo: repo}
	return &module{service: svc}
}
