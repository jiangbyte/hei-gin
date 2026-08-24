// internal/modules/sys/codegen/service.go 业务服务。
//
// Author: Charlie

package codegen

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/platform/idgen"
	"hei-gin/internal/framework/platform/module"
)

// Service 代码生成业务服务。
//
// Author: Charlie
type Service struct{ repo *Repo }

// NewService 构造代码生成服务。
func NewService(db *gorm.DB) *Service { return &Service{repo: NewRepo(db)} }

// New 构建 sys.codegen 模块。
func New(d *module.Deps) module.Module {
	s := NewService(d.DB)
	return module.Module{
		Name:   "sys.codegen",
		Models: []any{&Plan{}, &Field{}},
		Routes: []module.RouteRegistrar{s.registerRoutes(d)},
	}
}

// Create 创建方案并反射同步字段（先校验方案，对齐 hei-boot CodegenServiceImpl.create）。
func (s *Service) Create(ctx context.Context, req AddParam) error {
	if err := s.validatePlan(ctx, req); err != nil {
		return err
	}
	plan := fromAddParam(req)
	if plan.PKColumn == "" {
		plan.PKColumn = "id"
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

// Update 更新方案并重新同步字段（先校验方案）。
func (s *Service) Update(ctx context.Context, req EditParam) error {
	if err := s.validatePlan(ctx, req.AddParam); err != nil {
		return err
	}
	plan := fromAddParam(req.AddParam)
	plan.ID = req.ID
	if plan.PKColumn == "" {
		plan.PKColumn = "id"
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
			ColumnName: r.ColumnName,
			Label:      r.ColumnComment,
			DBType:     r.UDTName,
			ValueType:  py,
			UIType:     ts,
			PrimaryKey: r.IsPrimaryKey,
			Nullable:   strings.EqualFold(r.IsNullable, "YES"),
			MaxLength:  r.MaxLength,
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
			ID:            idgen.Next(),
			PlanID:        req.PlanID,
			TableRole:     item.TableRole,
			ColumnName:    item.ColumnName,
			Label:         item.Label,
			DBType:        item.DBType,
			ValueType:     def(item.ValueType, "str"),
			UIType:        def(item.UIType, "string"),
			Widget:        def(item.Widget, "input"),
			DictCode:      item.DictCode,
			QueryOperator: item.QueryOperator,
			InTable:       item.InTable,
			InForm:        item.InForm,
			InDetail:      item.InDetail,
			InQuery:       item.InQuery,
			PrimaryKey:    item.PrimaryKey,
			Required:      item.Required,
			UniqueFlag:    item.UniqueFlag,
			Nullable:      item.Nullable,
			MaxLength:     item.MaxLength,
			Sort:          item.Sort,
			CreatedAt:     now,
			UpdatedAt:     now,
		})
	}
	return s.repo.ReplaceFields(ctx, req.PlanID, fields)
}

// ParentResources 父级资源树（对齐 hei-boot ResourceMenuApi.listParentMenus + parentResources）。
func (s *Service) ParentResources(ctx context.Context, moduleID string) ([]ResourceNode, error) {
	db := s.repo.db.WithContext(ctx)
	moduleIDs, err := s.adminModuleIDs(ctx, moduleID)
	if err != nil {
		return nil, err
	}
	if len(moduleIDs) == 0 {
		return []ResourceNode{}, nil
	}
	var rows []ResourceNode
	err = db.Table("sys_resource").
		Select("id, parent_id, name, resource_type, sort, status").
		Where("status = ? AND resource_type IN ? AND module_id IN ?",
			security.StatusEnabled, []string{"CATALOG", "MENU", "PAGE"}, moduleIDs).
		Order("sort asc, id asc").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	normalizeResourceParents(rows)
	return buildResourceTree(rows, nil), nil
}

func (s *Service) adminModuleIDs(ctx context.Context, moduleID string) ([]string, error) {
	db := s.repo.db.WithContext(ctx)
	if moduleID != "" {
		var cnt int64
		if err := db.Table("sys_resource_module").
			Where("id = ? AND client = ? AND status = ?", moduleID, "ADMIN", security.StatusEnabled).
			Count(&cnt).Error; err != nil {
			return nil, err
		}
		if cnt == 0 {
			return nil, nil
		}
		return []string{moduleID}, nil
	}
	var ids []string
	err := db.Table("sys_resource_module").
		Where("client = ? AND status = ?", "ADMIN", security.StatusEnabled).
		Pluck("id", &ids).Error
	return ids, err
}

// 生成类型常量（对齐 hei-boot CodegenServiceImpl TREE_TYPES / RELATION_TYPES）。
var (
	treeGenTypes     = map[string]bool{"TREE": true, "LEFT_TREE_TABLE": true}
	relationGenTypes = map[string]bool{"LEFT_TREE_TABLE": true, "MASTER_DETAIL": true}
)

// validatePlan 校验方案：主键/树字段/子表配置必须存在于表结构（对齐 hei-boot validatePlan）。
func (s *Service) validatePlan(ctx context.Context, req AddParam) error {
	mainNames, err := s.columnNameSet(ctx, req.Table)
	if err != nil {
		return err
	}
	mainPK := req.PKColumn
	if mainPK == "" {
		mainPK = "id"
	}
	if !mainNames[mainPK] {
		return fmt.Errorf("main primary key field does not exist")
	}
	if treeGenTypes[req.GenType] {
		if req.TreeParentField == nil || !mainNames[*req.TreeParentField] {
			return fmt.Errorf("tree parent field does not exist")
		}
		if req.TreeLabelField == nil || !mainNames[*req.TreeLabelField] {
			return fmt.Errorf("tree label field does not exist")
		}
	}
	if relationGenTypes[req.GenType] {
		if req.SubTable == nil || *req.SubTable == "" || req.SubPK == nil || *req.SubPK == "" || req.SubForeignKey == nil || *req.SubForeignKey == "" {
			return fmt.Errorf("sub table configuration is incomplete")
		}
		subNames, err := s.columnNameSet(ctx, *req.SubTable)
		if err != nil {
			return err
		}
		if !subNames[*req.SubPK] {
			return fmt.Errorf("sub primary key field does not exist")
		}
		if !subNames[*req.SubForeignKey] {
			return fmt.Errorf("sub foreign key field does not exist")
		}
	}
	return nil
}

// columnNameSet 查询表列名集合。
func (s *Service) columnNameSet(ctx context.Context, tableName string) (map[string]bool, error) {
	cols, err := s.repo.ListColumns(ctx, tableName)
	if err != nil {
		return nil, err
	}
	out := make(map[string]bool, len(cols))
	for _, c := range cols {
		out[c.ColumnName] = true
	}
	return out, nil
}

// syncReflectedFields 按表结构反射同步 MAIN/SUB 字段。
func (s *Service) syncReflectedFields(ctx context.Context, plan *Plan) error {
	if err := s.upsertReflected(ctx, plan.ID, "MAIN", plan.Table); err != nil {
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
				InTable: !isAudit, InForm: !c.IsPrimaryKey && !isAudit, InDetail: true,
				InQuery: isDefaultQuery(c.ColumnName),
				Widget:  widget,
				DictCode: defaultDictCode(c.ColumnName),
				CreatedAt: now,
			}
		}
		f.Label = blankToNil(c.ColumnComment)
		f.DBType = c.UDTName
		f.ValueType = py
		f.UIType = ts
		f.PrimaryKey = c.IsPrimaryKey
		f.Required = !strings.EqualFold(c.IsNullable, "YES") && !c.IsPrimaryKey && !isAudit
		f.UniqueFlag = false
		f.Nullable = strings.EqualFold(c.IsNullable, "YES")
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
			"label": f.Label, "db_type": f.DBType, "value_type": f.ValueType,
			"ui_type": f.UIType, "primary_key": f.PrimaryKey, "required": f.Required,
			"unique_flag": f.UniqueFlag, "nullable": f.Nullable, "max_length": f.MaxLength,
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
		Table: req.Table, PKColumn: req.PKColumn, EntityName: req.EntityName,
		ModulePath: req.ModulePath, BusinessName: req.BusinessName,
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
		"table_name": p.Table, "pk_column": p.PKColumn, "entity_name": p.EntityName,
		"module_path": p.ModulePath, "business_name": p.BusinessName,
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
	ID           string         `json:"id" gorm:"column:id"`
	ParentID     *string        `json:"parent_id" gorm:"column:parent_id"`
	Name         string         `json:"name" gorm:"column:name"`
	ResourceType string         `json:"resource_type" gorm:"column:resource_type"`
	Sort         int            `json:"sort" gorm:"column:sort"`
	Status       string         `json:"status" gorm:"column:status"`
	Children     []ResourceNode `json:"children" gorm:"-"`
}

func normalizeResourceParents(rows []ResourceNode) {
	ids := make(map[string]struct{}, len(rows))
	for _, r := range rows {
		ids[r.ID] = struct{}{}
	}
	for i := range rows {
		if rows[i].ParentID == nil || strings.TrimSpace(*rows[i].ParentID) == "" {
			rows[i].ParentID = nil
			continue
		}
		if _, ok := ids[*rows[i].ParentID]; !ok {
			rows[i].ParentID = nil
		}
	}
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
