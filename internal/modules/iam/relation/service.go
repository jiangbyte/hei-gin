package relation

import (
	"context"
	"encoding/json"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"hei-gin/internal/framework/core/security"
	"hei-gin/internal/framework/platform/idgen"
	"hei-gin/internal/framework/platform/module"
	"hei-gin/internal/modules/shared"
)

// New æž„å»º iam.relation æ¨¡å—ã€‚
func New(_ *shared.Deps) module.Module {
	return module.Module{
		Name:   "iam.relation",
		Models: []any{&Relation{}},
	}
}

// Service å…³ç³»æœåŠ¡ï¼šä¸»ä½“-ç›®æ ‡å…³ç³»æŸ¥è¯¢ä¸Žå…¨é‡æ›¿æ¢æŽˆæƒï¼ˆå…ˆåˆ åŽæ’ï¼Œäº‹åŠ¡ï¼‰ã€‚
//
// Author: Charlie
type Service struct{ repo *Repo }

// NewService æž„é€ å…³ç³»æœåŠ¡ã€‚
func NewService(db *gorm.DB) *Service { return &Service{repo: NewRepo(db)} }

// ListTargetIDs åˆ—å‡ºä¸»ä½“å·²å…³è”ç›®æ ‡ IDï¼ˆaccountType ä¸ºç©ºä¸è¿‡æ»¤ï¼‰ã€‚
func (s *Service) ListTargetIDs(ctx context.Context, subjectType, subjectID, relationType, accountType string) ([]string, error) {
	rows, err := s.repo.ListRelations(ctx, subjectType, subjectID, relationType, accountType)
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(rows))
	seen := map[string]struct{}{}
	for _, r := range rows {
		if r.TargetID == "" {
			continue
		}
		if _, ok := seen[r.TargetID]; ok {
			continue
		}
		seen[r.TargetID] = struct{}{}
		out = append(out, r.TargetID)
	}
	return out, nil
}

// ListDeptGrants åˆ—å‡ºè´¦å·å·²æ‹¥æœ‰éƒ¨é—¨æŽˆäºˆæ˜Žç»†ã€‚
func (s *Service) ListDeptGrants(ctx context.Context, accountID, accountType string) ([]DeptGrantInfo, error) {
	rows, err := s.repo.ListRelations(ctx, SubjectAccount, accountID, RelAccountDept, accountType)
	if err != nil {
		return nil, err
	}
	out := make([]DeptGrantInfo, 0, len(rows))
	for _, r := range rows {
		out = append(out, DeptGrantInfo{
			DeptID:             r.TargetID,
			DataScope:          r.DataScope,
			CustomScopeDeptIDs: parseStringList(r.CustomScopeDeptIDs),
		})
	}
	return out, nil
}

// ListResourceGrants åˆ—å‡ºä¸»ä½“å·²æ‹¥æœ‰èµ„æºæŽˆäºˆæ˜Žç»†ï¼ˆç®¡ç†ç«¯/å®¢æˆ·ç«¯ï¼‰ã€‚
func (s *Service) ListResourceGrants(ctx context.Context, subjectType, subjectID, relationType, targetType, accountType string) ([]ResourceGrantInfo, error) {
	rows, err := s.repo.ListRelations(ctx, subjectType, subjectID, relationType, accountType)
	if err != nil {
		return nil, err
	}
	out := make([]ResourceGrantInfo, 0, len(rows))
	for _, r := range rows {
		out = append(out, ResourceGrantInfo{
			ResourceID: r.TargetID,
			GrantMode:  r.GrantMode,
			DataScope:  r.DataScope,
		})
	}
	return out, nil
}

// ReplaceTargetIDs å…ˆåˆ åŽæ’å…¨é‡æ›¿æ¢ä¸»ä½“-ç›®æ ‡ç®€å•å…³ç³»ï¼ˆè§’è‰²/ç”¨æˆ·ç»„ç­‰ï¼‰ã€‚
func (s *Service) ReplaceTargetIDs(ctx context.Context, subjectType, subjectID, relationType, targetType, accountType string, targetIDs []string) error {
	return s.repo.with(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.repo.deleteSubjectRelations(tx, subjectType, subjectID, relationType, accountType); err != nil {
			return err
		}
		rows := make([]Relation, 0, len(targetIDs))
		seen := map[string]struct{}{}
		for _, id := range targetIDs {
			if id == "" {
				continue
			}
			if _, ok := seen[id]; ok {
				continue
			}
			seen[id] = struct{}{}
			rows = append(rows, newRelation(subjectType, subjectID, accountType, relationType, targetType, id))
		}
		if len(rows) == 0 {
			return nil
		}
		return s.repo.CreateInBatches(tx, rows)
	})
}

// ReplaceDeptGrants å…ˆåˆ åŽæ’å…¨é‡æ›¿æ¢è´¦å·-éƒ¨é—¨æŽˆäºˆã€‚
func (s *Service) ReplaceDeptGrants(ctx context.Context, accountID, accountType string, grants []DeptGrantInfo) error {
	return s.repo.with(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.repo.deleteSubjectRelations(tx, SubjectAccount, accountID, RelAccountDept, accountType); err != nil {
			return err
		}
		rows := make([]Relation, 0, len(grants))
		for _, g := range grants {
			if g.DeptID == "" {
				continue
			}
			rel := newRelation(SubjectAccount, accountID, accountType, RelAccountDept, TargetDept, g.DeptID)
			rel.DataScope = orDef(g.DataScope, string(security.DataScopeAll))
			rel.CustomScopeDeptIDs = jsonList(g.CustomScopeDeptIDs)
			rows = append(rows, rel)
		}
		if len(rows) == 0 {
			return nil
		}
		return s.repo.CreateInBatches(tx, rows)
	})
}

// ReplaceResourceGrants å…ˆåˆ åŽæ’å…¨é‡æ›¿æ¢ä¸»ä½“-èµ„æºæŽˆäºˆï¼ˆç®¡ç†ç«¯/å®¢æˆ·ç«¯ï¼‰ã€‚
func (s *Service) ReplaceResourceGrants(ctx context.Context, subjectType, subjectID, relationType, targetType, accountType string, grants []ResourceGrantInfo) error {
	return s.repo.with(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.repo.deleteSubjectRelations(tx, subjectType, subjectID, relationType, accountType); err != nil {
			return err
		}
		rows := make([]Relation, 0, len(grants))
		for _, g := range grants {
			if g.ResourceID == "" {
				continue
			}
			rel := newRelation(subjectType, subjectID, accountType, relationType, targetType, g.ResourceID)
			rel.GrantMode = orDef(g.GrantMode, GrantCascade)
			rel.DataScope = orDef(g.DataScope, string(security.DataScopeAll))
			rows = append(rows, rel)
		}
		if len(rows) == 0 {
			return nil
		}
		return s.repo.CreateInBatches(tx, rows)
	})
}

// ReplaceSubjectAccounts å…ˆåˆ åŽæ’å…¨é‡æ›¿æ¢ä¸»ä½“-è´¦å·æˆå‘˜ï¼ˆGROUP_USER/ROLE_USERï¼‰ã€‚
func (s *Service) ReplaceSubjectAccounts(ctx context.Context, subjectType, subjectID, relationType string, accountIDs []string, accountTypes map[string]string) error {
	return s.repo.with(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.repo.deleteSubjectRelations(tx, subjectType, subjectID, relationType, ""); err != nil {
			return err
		}
		rows := make([]Relation, 0, len(accountIDs))
		seen := map[string]struct{}{}
		for _, id := range accountIDs {
			if id == "" {
				continue
			}
			if _, ok := seen[id]; ok {
				continue
			}
			seen[id] = struct{}{}
			accType := accountTypes[id]
			if accType == "" {
				return gorm.ErrRecordNotFound
			}
			rows = append(rows, newRelation(subjectType, subjectID, accType, relationType, TargetAccount, id))
		}
		if len(rows) == 0 {
			return nil
		}
		return s.repo.CreateInBatches(tx, rows)
	})
}

// BindResourcePermissions å…ˆåˆ åŽæ’ä¸ºèµ„æºç»‘å®šæƒé™é”®ï¼ˆRESOURCE_PERMISSION/CLIENT_RESOURCE_PERMISSIONï¼‰ã€‚
func (s *Service) BindResourcePermissions(ctx context.Context, subjectType, subjectID, relationType, accountType string, permissionKeys []string) error {
	return s.repo.with(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.repo.deleteSubjectRelations(tx, subjectType, subjectID, relationType, accountType); err != nil {
			return err
		}
		rows := make([]Relation, 0, len(permissionKeys))
		for _, key := range permissionKeys {
			if key == "" {
				continue
			}
			rel := newRelation(subjectType, subjectID, accountType, relationType, TargetPermission, "")
			rel.TargetKey = key
			rel.GrantMode = GrantCascade
			rel.DataScope = string(security.DataScopeAll)
			rows = append(rows, rel)
		}
		if len(rows) == 0 {
			return nil
		}
		return s.repo.CreateInBatches(tx, rows)
	})
}

// DeleteBySubjectIDs æŒ‰ä¸»ä½“ id é›†åˆåˆ é™¤æŒ‡å®šå…³ç³»ç±»åž‹çš„å…³ç³»ã€‚
func (s *Service) DeleteBySubjectIDs(ctx context.Context, subjectType string, subjectIDs []string, relationType string) error {
	return s.repo.DeleteBySubjectIDs(ctx, subjectType, subjectIDs, relationType)
}

// newRelation æž„é€ é»˜è®¤å¯ç”¨å…³ç³»è¡Œã€‚
func newRelation(subjectType, subjectID, accountType, relationType, targetType, targetID string) Relation {
	return Relation{
		ID:                 idgen.Next(),
		SubjectType:        subjectType,
		SubjectID:          subjectID,
		AccountType:        accountType,
		RelationType:       relationType,
		TargetType:         targetType,
		TargetID:           targetID,
		GrantMode:          GrantCascade,
		DataScope:          string(security.DataScopeAll),
		CustomScopeDeptIDs: datatypes.JSON([]byte("[]")),
		Status:             security.StatusEnabled,
		Extra:              datatypes.JSON([]byte("{}")),
	}
}

func orDef(s, d string) string {
	if s == "" {
		return d
	}
	return s
}

// parseStringList è§£æž jsonb å­—ç¬¦ä¸²æ•°ç»„ã€‚
func parseStringList(raw datatypes.JSON) []string {
	var out []string
	if len(raw) == 0 {
		return []string{}
	}
	if err := json.Unmarshal(raw, &out); err != nil {
		return []string{}
	}
	if out == nil {
		out = []string{}
	}
	return out
}

// jsonList å°†å­—ç¬¦ä¸²æ•°ç»„ç¼–ç ä¸º jsonbã€‚
func jsonList(items []string) datatypes.JSON {
	if items == nil {
		items = []string{}
	}
	raw, err := json.Marshal(items)
	if err != nil {
		return datatypes.JSON([]byte("[]"))
	}
	return datatypes.JSON(raw)
}
