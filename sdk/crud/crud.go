package crud

import (
	"gorm.io/gorm"

	"hei-gin/sdk/db"
	"hei-gin/sdk/exception"
	"hei-gin/sdk/result"

	"github.com/gin-gonic/gin"
)

// PageParam is the interface for pagination parameters.
type PageParam interface {
	GetCurrent() int
	GetSize() int
}

func pageNum(p PageParam) (current, size int) {
	current = p.GetCurrent()
	if current < 1 {
		current = 1
	}
	size = p.GetSize()
	if size < 1 || size > 100 {
		size = 10
	}
	return
}

// Page executes a standard paginated query.
// model: pointer type (e.g. &SysBanner{})
// toVO: func(entity *T) any
func Page[T any, P any](c *gin.Context, model T, param *P, buildQuery func(q *gorm.DB) *gorm.DB, order string, toVO func(entity T) any) {
	ctx := c.Request.Context()
	// P must be *T-like for crud.PageParams, so use type assertion
	var pp PageParam
	if p, ok := any(param).(PageParam); ok {
		pp = p
	}
	current, size := 1, 10
	if pp != nil {
		current, size = pageNum(pp)
	}
	offset := (current - 1) * size

	// Use a pointer's base type for gorm queries
	query := db.DB.WithContext(ctx).Model(model)
	if buildQuery != nil {
		query = buildQuery(query)
	}

	var total int64
	query.Count(&total)

	var records []T
	if order == "" {
		order = "created_at DESC"
	}
	query.Order(order).Limit(size).Offset(offset).Find(&records)

	vos := make([]any, len(records))
	for i := range records {
		vos[i] = toVO(records[i])
	}
	result.PageDataResult(c, vos, total, current, size)
}

// Detail retrieves a single record by ID.
// model should be a pointer type (e.g. &SysBanner{}), it will be populated.
func Detail[T any](c *gin.Context, model T, id, name string) {
	ctx := c.Request.Context()
	if err := db.DB.WithContext(ctx).Where("id = ?", id).First(model).Error; err != nil {
		panic(exception.NewBusinessError(name+"不存在: "+err.Error(), 500))
	}
}

// Options returns a sorted list of all records.
func Options[T any](c *gin.Context, model T, order string, toVO func(entity T) any) []any {
	ctx := c.Request.Context()
	var records []T
	if order == "" {
		order = "sort_code ASC"
	}
	db.DB.WithContext(ctx).Order(order).Find(&records)
	vos := make([]any, len(records))
	for i := range records {
		vos[i] = toVO(records[i])
	}
	return vos
}

// Remove deletes records by IDs.
func Remove[T any](c *gin.Context, model T, ids []string) {
	if len(ids) == 0 {
		return
	}
	ctx := c.Request.Context()
	db.DB.WithContext(ctx).Where("id IN ?", ids).Delete(model)
}
