package permission

import (
	"encoding/json"
	"sort"

	"github.com/gin-gonic/gin"
)

type Service struct {
	repo *repository
}

// PermissionListModules returns sorted permission module names from Redis.
func (s *Service) ListModules(c *gin.Context) []string {
	ctx := c.Request.Context()
	data, err := s.repo.Cache(ctx)
	if err != nil {
		return []string{}
	}
	var tree map[string]interface{}
	if err := json.Unmarshal([]byte(data), &tree); err != nil {
		return []string{}
	}
	modules := make([]string, 0, len(tree))
	for k := range tree {
		modules = append(modules, k)
	}
	sort.Strings(modules)
	return modules
}

// PermissionListByModule returns permission list for a specific module from Redis.
func (s *Service) ListByModule(c *gin.Context, module string) []interface{} {
	ctx := c.Request.Context()
	data, err := s.repo.Cache(ctx)
	if err != nil {
		return []interface{}{}
	}
	var tree map[string]interface{}
	if err := json.Unmarshal([]byte(data), &tree); err != nil {
		return []interface{}{}
	}
	modulePerms, ok := tree[module].(map[string]interface{})
	if !ok {
		return []interface{}{}
	}
	perms := make([]interface{}, 0, len(modulePerms))
	for _, v := range modulePerms {
		perms = append(perms, v)
	}
	return perms
}
