package codegen

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"hei-gin/framework/platform/idgen"
	"hei-gin/framework/platform/module"
	"hei-gin/modules/shared"
)

// Service 代码生成业务服务。
//
// Author: Charlie
type Service struct{ repo *Repo }

// NewService 构造代码生成服务。
func NewService(db *gorm.DB) *Service { return &Service{repo: NewRepo(db)} }

// New 构建 sys.codegen 模块。
func New(d *shared.Deps) module.Module {
	s := NewService(d.DB)
	return module.Module{
		Name:   "sys.codegen",
		Models: []any{&Plan{}, &Field{}},
		Routes: []module.RouteRegistrar{s.registerRoutes(d)},
	}
}

// Create 创建方案并反射同步字段。
func (s *Service) Create(ctx context.Context, req AddParam) error {
	plan := fromAddParam(req)
	if plan.MainPK == "" {
		plan.MainPK = "id"
	}
	if plan.Sort == 0 {
		plan.Sort = 99
	}
	count, err := s.repo.CountByName(ctx, plan.Name, "")
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("codegen plan name already exists")
	}
	if err := s.repo.Create(ctx, plan); err != nil {
		return err
	}
	return s.syncReflectedFields(ctx, plan)
}

// Update 更新方案并重新同步字段。
func (s *Service) Update(ctx context.Context, req EditParam) error {
	plan := fromAddParam(req.AddParam)
	plan.ID = req.ID
	if plan.MainPK == "" {
		plan.MainPK = "id"
	}
	if plan.Sort == 0 {
		plan.Sort = 99
	}
	count, err := s.repo.CountByName(ctx, plan.Name, plan.ID)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("codegen plan name already exists")
	}
	existing, err := s.repo.GetByID(ctx, plan.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("codegen plan not found")
		}
		return err
	}
	updates := planUpdates(plan)
	updates["created_at"] = existing.CreatedAt
	if err := s.repo.Update(ctx, plan.ID, updates); err != nil {
		return err
	}
	return s.syncReflectedFields(ctx, plan)
}

// Delete 批量删除。
func (s *Service) Delete(ctx context.Context, ids []string) error {
	return s.repo.DeleteByIDs(ctx, ids)
}

// Detail 详情。
func (s *Service) Detail(ctx context.Context, id string) (*Plan, error) {
	return s.repo.GetByID(ctx, id)
}

// Page 分页。
func (s *Service) Page(ctx context.Context, q PageParam) (rows []Plan, total int64, current, size int, err error) {
	current, size = q.Normalize()
	rows, total, err = s.repo.Page(ctx, q)
	return rows, total, current, size, err
}

// Tables 数据库表列表。
func (s *Service) Tables(ctx context.Context) ([]DatabaseTableResult, error) {
	rows, err := s.repo.ListTables(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]DatabaseTableResult, 0, len(rows))
	for _, r := range rows {
		out = append(out, DatabaseTableResult{TableName: r.TableName, TableComment: r.TableComment})
	}
	return out, nil
}

// TableColumns 表列元数据。
func (s *Service) TableColumns(ctx context.Context, tableName string) ([]DatabaseColumnResult, error) {
	rows, err := s.repo.ListColumns(ctx, tableName)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("database table not found")
	}
	out := make([]DatabaseColumnResult, 0, len(rows))
	for _, r := range rows {
		py, ts := dbTypeToPythonAndTs(r.DataType, r.UDTName)
		out = append(out, DatabaseColumnResult{
			ColumnName:     r.ColumnName,
			ColumnComment:  r.ColumnComment,
			DBType:         r.UDTName,
			PythonType:     py,
			TypescriptType: ts,
			IsPrimaryKey:   r.IsPrimaryKey,
			IsNullable:     strings.EqualFold(r.IsNullable, "YES"),
			MaxLength:      r.MaxLength,
		})
	}
	return out, nil
}

// Fields 查询方案字段。
func (s *Service) Fields(ctx context.Context, q FieldQuery) ([]Field, error) {
	if _, err := s.repo.GetByID(ctx, q.PlanID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("codegen plan not found")
		}
		return nil, err
	}
	return s.repo.ListFields(ctx, q.PlanID, q.TableRole)
}

// UpdateFieldsBatch 批量更新方案字段。
func (s *Service) UpdateFieldsBatch(ctx context.Context, req FieldsUpdateBatchParam) error {
	if _, err := s.repo.GetByID(ctx, req.PlanID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("codegen plan not found")
		}
		return err
	}
	now := time.Now().UTC()
	fields := make([]Field, 0, len(req.Fields))
	for _, item := range req.Fields {
		fields = append(fields, Field{
			ID:             idgen.Next(),
			PlanID:         req.PlanID,
			TableRole:      item.TableRole,
			ColumnName:     item.ColumnName,
			ColumnComment:  item.ColumnComment,
			DBType:         item.DBType,
			PythonType:     def(item.PythonType, "str"),
			TypescriptType: def(item.TypescriptType, "string"),
			FormWidget:     def(item.FormWidget, "input"),
			DictCode:       item.DictCode,
			QueryOperator:  item.QueryOperator,
			ShowInTable:    item.ShowInTable,
			ShowInForm:     item.ShowInForm,
			ShowInDetail:   item.ShowInDetail,
			ShowInQuery:    item.ShowInQuery,
			IsPrimaryKey:   item.IsPrimaryKey,
			IsRequired:     item.IsRequired,
			IsUnique:       item.IsUnique,
			IsNullable:     item.IsNullable,
			MaxLength:      item.MaxLength,
			Sort:           item.Sort,
			CreatedAt:      now,
			UpdatedAt:      now,
		})
	}
	return s.repo.ReplaceFields(ctx, req.PlanID, fields)
}

// ParentResources 父级资源树（CATALOG 类型，可选模块过滤）。
func (s *Service) ParentResources(ctx context.Context, moduleID string) ([]ResourceNode, error) {
	db := s.repo.db.WithContext(ctx)
	q := db.Table("sys_resource").
		Select("id, parent_id, name, resource_type, sort, status").
		Where("resource_type IN ?", []string{"CATALOG", "MENU"})
	if moduleID != "" {
		q = q.Where("module_id = ?", moduleID)
	}
	var rows []ResourceNode
	if err := q.Order("sort asc, id asc").Find(&rows).Error; err != nil {
		return nil, err
	}
	return buildResourceTree(rows, nil), nil
}

// syncReflectedFields 按表结构反射同步 MAIN/SUB 字段。
func (s *Service) syncReflectedFields(ctx context.Context, plan *Plan) error {
	if err := s.upsertReflected(ctx, plan.ID, "MAIN", plan.MainTable); err != nil {
		return err
	}
	if isRelationType(plan.GenType) && plan.SubTable != nil && *plan.SubTable != "" {
		return s.upsertReflected(ctx, plan.ID, "SUB", *plan.SubTable)
	}
	return nil
}

func (s *Service) upsertReflected(ctx context.Context, planID, tableRole, tableName string) error {
	cols, err := s.repo.ListColumns(ctx, tableName)
	if err != nil {
		return err
	}
	existing := map[string]Field{}
	rows, err := s.repo.ListFields(ctx, planID, tableRole)
	if err != nil {
		return err
	}
	for _, f := range rows {
		existing[f.ColumnName] = f
	}
	now := time.Now().UTC()
	toCreate := make([]Field, 0, len(cols))
	toUpdate := make([]Field, 0, len(cols))
	index := 1
	for _, c := range cols {
		py, ts := dbTypeToPythonAndTs(c.DataType, c.UDTName)
		isAudit := isAuditColumn(c.ColumnName)
		widget := defaultWidget(c.ColumnName, py)
		f, ok := existing[c.ColumnName]
		if !ok {
			f = Field{
				ID: idgen.Next(), PlanID: planID, TableRole: tableRole, ColumnName: c.ColumnName,
				ShowInTable: !isAudit, ShowInForm: !c.IsPrimaryKey && !isAudit, ShowInDetail: true,
				ShowInQuery: isDefaultQuery(c.ColumnName),
				FormWidget:  widget,
				DictCode:    defaultDictCode(c.ColumnName),
				CreatedAt:   now,
			}
		}
		f.ColumnComment = blankToNil(c.ColumnComment)
		f.DBType = c.UDTName
		f.PythonType = py
		f.TypescriptType = ts
		f.IsPrimaryKey = c.IsPrimaryKey
		f.IsRequired = !strings.EqualFold(c.IsNullable, "YES") && !c.IsPrimaryKey && !isAudit
		f.IsUnique = false
		f.IsNullable = strings.EqualFold(c.IsNullable, "YES")
		f.MaxLength = c.MaxLength
		f.Sort = c.Sort
		f.UpdatedAt = now
		if f.ID == "" {
			toCreate = append(toCreate, f)
		} else {
			toUpdate = append(toUpdate, f)
		}
		index++
	}
	if len(toCreate) > 0 {
		if err := s.repo.with(ctx).Create(&toCreate).Error; err != nil {
			return err
		}
	}
	for _, f := range toUpdate {
		if err := s.repo.with(ctx).Model(&Field{}).Where("id = ?", f.ID).Updates(map[string]any{
			"column_comment": f.ColumnComment, "db_type": f.DBType, "python_type": f.PythonType,
			"typescript_type": f.TypescriptType, "is_primary_key": f.IsPrimaryKey, "is_required": f.IsRequired,
			"is_unique": f.IsUnique, "is_nullable": f.IsNullable, "max_length": f.MaxLength,
			"sort": f.Sort, "updated_at": f.UpdatedAt,
		}).Error; err != nil {
			return err
		}
	}
	return nil
}

// fromAddParam 入参 → Plan 实体。
func fromAddParam(req AddParam) *Plan {
	return &Plan{
		ID: idgen.Next(), Name: req.Name, GenType: req.GenType, Author: req.Author, Description: req.Description,
		MainTable: req.MainTable, MainPK: req.MainPK, MainEntityName: req.MainEntityName,
		MainModulePath: req.MainModulePath, MainBusinessName: req.MainBusinessName,
		APIPrefix: req.APIPrefix, PermissionPrefix: req.PermissionPrefix,
		ResourceModuleID: req.ResourceModuleID, ParentResourceID: req.ParentResourceID,
		MenuName: req.MenuName, MenuPath: req.MenuPath, ComponentPath: req.ComponentPath,
		Icon: req.Icon, Sort: req.Sort,
		TreeParentField: req.TreeParentField, TreeLabelField: req.TreeLabelField,
		SubTable: req.SubTable, SubPK: req.SubPK, SubForeignKey: req.SubForeignKey,
		SubEntityName: req.SubEntityName, SubBusinessName: req.SubBusinessName,
	}
}

func planUpdates(p *Plan) map[string]any {
	return map[string]any{
		"name": p.Name, "gen_type": p.GenType, "author": p.Author, "description": p.Description,
		"main_table": p.MainTable, "main_pk": p.MainPK, "main_entity_name": p.MainEntityName,
		"main_module_path": p.MainModulePath, "main_business_name": p.MainBusinessName,
		"api_prefix": p.APIPrefix, "permission_prefix": p.PermissionPrefix,
		"resource_module_id": p.ResourceModuleID, "parent_resource_id": p.ParentResourceID,
		"menu_name": p.MenuName, "menu_path": p.MenuPath, "component_path": p.ComponentPath,
		"icon": p.Icon, "sort": p.Sort,
		"tree_parent_field": p.TreeParentField, "tree_label_field": p.TreeLabelField,
		"sub_table": p.SubTable, "sub_pk": p.SubPK, "sub_foreign_key": p.SubForeignKey,
		"sub_entity_name": p.SubEntityName, "sub_business_name": p.SubBusinessName,
	}
}

// ResourceNode 父级资源树节点。
//
// Author: Charlie
type ResourceNode struct {
	ID           string         `json:"id"`
	ParentID     *string        `json:"parent_id"`
	Name         string         `json:"name"`
	ResourceType string         `json:"resource_type"`
	Sort         int            `json:"sort"`
	Status       string         `json:"status"`
	Children     []ResourceNode `json:"children"`
}

func buildResourceTree(rows []ResourceNode, parent *string) []ResourceNode {
	out := make([]ResourceNode, 0)
	for _, r := range rows {
		same := (r.ParentID == nil && parent == nil) || (r.ParentID != nil && parent != nil && *r.ParentID == *parent)
		if same {
			r.Children = buildResourceTree(rows, &r.ID)
			out = append(out, r)
		}
	}
	return out
}

// dbTypeToPythonAndTs 数据库类型 → Python/TS 语义类型。
func dbTypeToPythonAndTs(dataType, udtName string) (string, string) {
	t := strings.ToLower(udtName)
	if t == "" {
		t = strings.ToLower(dataType)
	}
	switch t {
	case "int2", "int4", "int8", "integer", "bigint", "smallint", "serial", "bigserial":
		return "int", "number"
	case "numeric", "decimal", "float4", "float8", "double precision", "real", "money":
		return "float", "number"
	case "bool", "boolean":
		return "bool", "boolean"
	case "json", "jsonb":
		return "dict", "Record<string, any>"
	case "timestamp", "timestamptz", "timestamp without time zone", "timestamp with time zone",
		"date", "time", "timetz":
		return "datetime", "string"
	default:
		return "str", "string"
	}
}

func isAuditColumn(name string) bool {
	switch name {
	case "created_at", "created_by", "updated_at", "updated_by":
		return true
	}
	return false
}

func isRelationType(genType string) bool {
	return genType == "LEFT_TREE_TABLE" || genType == "MASTER_DETAIL"
}

func defaultWidget(columnName, pythonType string) string {
	if columnName == "status" {
		return "dict"
	}
	if pythonType == "int" || pythonType == "float" {
		return "number"
	}
	if pythonType == "bool" {
		return "switch"
	}
	if strings.Contains(columnName, "content") || strings.Contains(columnName, "description") || strings.Contains(columnName, "remark") {
		return "textarea"
	}
	return "input"
}

func defaultDictCode(columnName string) *string {
	if columnName == "status" {
		v := "COMMON_STATUS"
		return &v
	}
	return nil
}

func isDefaultQuery(columnName string) bool {
	switch columnName {
	case "name", "title", "code", "status", "category", "type":
		return true
	}
	return false
}

func blankToNil(v string) *string {
	if strings.TrimSpace(v) == "" {
		return nil
	}
	return &v
}

func def(v, fallback string) string {
	if strings.TrimSpace(v) == "" {
		return fallback
	}
	return v
}
