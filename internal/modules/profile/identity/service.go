// internal/modules/profile/identity/service.go 实名认证业务服务。
//
// Author: Charlie
package identity

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"

	"hei-gin/internal/framework/platform/idgen"
	"hei-gin/internal/framework/platform/storage"
	"hei-gin/internal/modules/sys/file"
)

var documentTypes = []string{"ID_CARD", "PASSPORT", "EID"}

// Service 实名认证服务。
type Service struct {
	repo      *Repo
	crypto    *FieldCrypto
	providers *ProviderRegistry
	storage   *storage.Manager
}

// NewService 构造实名认证服务。
func NewService(db *gorm.DB, crypto *FieldCrypto, sto *storage.Manager) *Service {
	return &Service{
		repo:      NewRepo(db),
		crypto:    crypto,
		providers: NewProviderRegistry(MockVerifyProvider{}),
		storage:   sto,
	}
}

// IsVerified 账号是否已完成实名认证。
func (s *Service) IsVerified(ctx context.Context, accountID string) bool {
	row, err := s.repo.GetIdentity(ctx, accountID)
	return err == nil && row != nil && row.Status == StatusVerified
}

// GetUserStatus 用户侧实名状态（脱敏）。
func (s *Service) GetUserStatus(ctx context.Context, accountID string) (*IdentityStatusResult, error) {
	out, err := s.getStatus(ctx, accountID)
	if err != nil {
		return nil, err
	}
	return sanitizeStatus(out), nil
}

// GetStatus 管理侧实名状态（完整字段）。
func (s *Service) GetStatus(ctx context.Context, accountID string) (*IdentityStatusResult, error) {
	return s.getStatus(ctx, accountID)
}

func (s *Service) getStatus(ctx context.Context, accountID string) (*IdentityStatusResult, error) {
	out := &IdentityStatusResult{Status: StatusUnverified}
	row, err := s.repo.GetIdentity(ctx, accountID)
	if err == nil && row != nil {
		out.Status = row.Status
		out.DocumentType = row.DocumentType
		out.VerifyChannel = row.VerifyChannel
		out.Provider = row.Provider
		out.VerifiedAt = row.VerifiedAt
		out.RevokedAt = row.RevokedAt
		if row.RealNameCipher != nil && *row.RealNameCipher != "" {
			plain, err := s.crypto.Decrypt(*row.RealNameCipher)
			if err == nil {
				if masked := s.crypto.MaskRealName(plain); masked != nil {
					out.RealNameMasked = *masked
				}
			}
		}
		if row.DocumentNoCipher != nil && *row.DocumentNoCipher != "" {
			plain, err := s.crypto.Decrypt(*row.DocumentNoCipher)
			if err == nil {
				if masked := s.crypto.MaskDocumentNo(plain); masked != nil {
					out.DocumentNoMasked = *masked
				}
			}
		}
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	pending, err := s.repo.FindLatestPendingCase(ctx, accountID)
	if err == nil && pending != nil {
		summary := s.toSummary(pending)
		out.PendingCase = &summary
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return out, nil
}

// Options 实名认证可选项。
func (s *Service) Options(_ context.Context) *RealNameCaseOptionsResult {
	return &RealNameCaseOptionsResult{
		BusinessTypes: []RealNameBusinessOptionResult{{
			BusinessType: BusinessAccountVerify,
			Label:        "账号实名认证",
			Channels:     []string{ChannelManual, ChannelThirdParty},
		}},
		DocumentTypes: append([]string(nil), documentTypes...),
	}
}

// Submit 提交人工实名认证工单。
func (s *Service) Submit(ctx context.Context, accountID string, req RealNameCaseSubmitParam) error {
	businessType := normalizeBusinessType(req.BusinessType)
	if err := s.validateSubmit(ctx, accountID, businessType, req.DocumentType, req.DocumentNo); err != nil {
		return err
	}
	attachments := s.normalizeAttachmentIDs(req.AttachmentIDs)
	if len(attachments) == 0 {
		return bizErr(400, 400, "请上传证件材料")
	}
	entity, err := s.newCaseEntity(accountID, businessType, ChannelManual, req.DocumentType, req.RealName, req.DocumentNo)
	if err != nil {
		return err
	}
	entity.AttachmentIDs = attachments
	if strings.TrimSpace(req.ApplicantContact) != "" {
		cipher, err := s.crypto.Encrypt(req.ApplicantContact)
		if err != nil {
			return err
		}
		entity.ApplicantContactCipher = &cipher
	}
	submitter := accountID
	entity.SubmitterID = &submitter
	entity.CreatedBy = &accountID
	entity.UpdatedBy = &accountID
	if err := s.repo.CreateCase(ctx, entity); err != nil {
		return err
	}
	return s.repo.AppendRecord(ctx, entity, "SUBMIT", "", entity.Status, accountID, "")
}

// InitThirdParty 发起第三方实名认证。
func (s *Service) InitThirdParty(ctx context.Context, accountID string, req RealNameCaseInitThirdPartyParam) (*RealNameCaseInitResult, error) {
	businessType := normalizeBusinessType(req.BusinessType)
	if err := s.validateSubmit(ctx, accountID, businessType, req.DocumentType, req.DocumentNo); err != nil {
		return nil, err
	}
	entity, err := s.newCaseEntity(accountID, businessType, ChannelThirdParty, req.DocumentType, req.RealName, req.DocumentNo)
	if err != nil {
		return nil, err
	}
	submitter := accountID
	entity.SubmitterID = &submitter
	entity.CreatedBy = &accountID
	entity.UpdatedBy = &accountID
	if err := s.repo.CreateCase(ctx, entity); err != nil {
		return nil, err
	}
	provider, err := s.providers.Resolve(ChannelThirdParty, req.DocumentType, req.Provider)
	if err != nil {
		return nil, err
	}
	initResult := provider.InitVerify(entity, req)
	entity.Provider = &initResult.Provider
	entity.ProviderOrderNo = &initResult.ProviderOrderNo
	if err := s.repo.UpdateCase(ctx, entity); err != nil {
		return nil, err
	}
	if err := s.repo.AppendRecord(ctx, entity, "INIT_THIRD_PARTY", "", entity.Status, accountID, ""); err != nil {
		return nil, err
	}
	return &initResult, nil
}

// Callback 第三方实名认证回调。
func (s *Service) Callback(ctx context.Context, req RealNameCaseCallbackParam) error {
	entity, err := s.requireCase(ctx, req.CaseID)
	if err != nil {
		return err
	}
	if entity.Status != CaseStatusPending {
		return bizErr(400, 400, "Case is not pending")
	}
	provider, err := s.providers.Resolve(entity.VerifyChannel, deref(entity.DocumentType), deref(entity.Provider))
	if err != nil {
		return err
	}
	provider.HandleCallback(entity, req)
	before := entity.Status
	success := req.Success != nil && *req.Success
	now := time.Now().UTC()
	if success {
		entity.Status = CaseStatusApproved
		entity.ReviewedAt = &now
		if err := s.repo.UpdateCase(ctx, entity); err != nil {
			return err
		}
		if err := s.upsertOnApprove(ctx, entity); err != nil {
			return err
		}
		return s.repo.AppendRecord(ctx, entity, "CALLBACK", before, entity.Status, "SYSTEM", req.Message)
	}
	reason := strings.TrimSpace(req.Message)
	if reason == "" {
		reason = "Third-party verification failed"
	}
	entity.Status = CaseStatusRejected
	entity.ReviewedAt = &now
	entity.RejectReason = &reason
	if err := s.repo.UpdateCase(ctx, entity); err != nil {
		return err
	}
	return s.repo.AppendRecord(ctx, entity, "CALLBACK", before, entity.Status, "SYSTEM", reason)
}

// MyPage 我的实名工单分页。
func (s *Service) MyPage(ctx context.Context, accountID string, q RealNameCaseMyPageParam) ([]RealNameCaseSummaryResult, int64, int, int, error) {
	rows, total, current, size, err := s.repo.PageCasesByAccount(ctx, accountID, q)
	if err != nil {
		return nil, 0, current, size, err
	}
	out := make([]RealNameCaseSummaryResult, 0, len(rows))
	for i := range rows {
		summary := sanitizeSummary(s.toSummary(&rows[i]))
		out = append(out, summary)
	}
	return out, total, current, size, nil
}

// ReviewPage 审核队列分页。
func (s *Service) ReviewPage(ctx context.Context, q RealNameCaseReviewPageParam) ([]RealNameCaseSummaryResult, int64, int, int, error) {
	rows, total, current, size, err := s.repo.PageReviewCases(ctx, q)
	if err != nil {
		return nil, 0, current, size, err
	}
	out := make([]RealNameCaseSummaryResult, 0, len(rows))
	for i := range rows {
		out = append(out, s.toSummary(&rows[i]))
	}
	return out, total, current, size, nil
}

// Detail 工单详情。
func (s *Service) Detail(ctx context.Context, caseID string) (*RealNameCaseDetailResult, error) {
	entity, err := s.requireCase(ctx, caseID)
	if err != nil {
		return nil, err
	}
	summary := s.toSummary(entity)
	out := &RealNameCaseDetailResult{
		RealNameCaseSummaryResult: summary,
		Provider:                  entity.Provider,
		ProviderOrderNo:           entity.ProviderOrderNo,
		SubmitterID:               entity.SubmitterID,
		ReviewerID:                entity.ReviewerID,
		Attachments:             s.resolveAttachments(ctx, entity.AttachmentIDs),
	}
	return out, nil
}

// Approve 审核通过。
func (s *Service) Approve(ctx context.Context, reviewerID string, req RealNameCaseApproveParam) error {
	entity, err := s.requireCase(ctx, req.CaseID)
	if err != nil {
		return err
	}
	if entity.Status != CaseStatusPending {
		return bizErr(400, 400, "Case is not pending")
	}
	before := entity.Status
	now := time.Now().UTC()
	entity.Status = CaseStatusApproved
	entity.ReviewerID = &reviewerID
	entity.ReviewedAt = &now
	if err := s.repo.UpdateCase(ctx, entity); err != nil {
		return err
	}
	if err := s.upsertOnApprove(ctx, entity); err != nil {
		return err
	}
	return s.repo.AppendRecord(ctx, entity, "APPROVE", before, entity.Status, reviewerID, req.Remark)
}

// Reject 审核驳回。
func (s *Service) Reject(ctx context.Context, reviewerID string, req RealNameCaseRejectParam) error {
	entity, err := s.requireCase(ctx, req.CaseID)
	if err != nil {
		return err
	}
	if entity.Status != CaseStatusPending {
		return bizErr(400, 400, "Case is not pending")
	}
	before := entity.Status
	now := time.Now().UTC()
	reason := strings.TrimSpace(req.RejectReason)
	entity.Status = CaseStatusRejected
	entity.ReviewerID = &reviewerID
	entity.ReviewedAt = &now
	entity.RejectReason = &reason
	if err := s.repo.UpdateCase(ctx, entity); err != nil {
		return err
	}
	return s.repo.AppendRecord(ctx, entity, "REJECT", before, entity.Status, reviewerID, reason)
}

// IdentityPage 实名快照分页。
func (s *Service) IdentityPage(ctx context.Context, q IdentityPageParam) ([]IdentityPageResult, int64, int, int, error) {
	rows, total, current, size, err := s.repo.PageIdentity(ctx, q)
	if err != nil {
		return nil, 0, current, size, err
	}
	out := make([]IdentityPageResult, 0, len(rows))
	for i := range rows {
		out = append(out, s.toPageResult(&rows[i]))
	}
	return out, total, current, size, nil
}

// Revoke 撤销实名快照。
func (s *Service) Revoke(ctx context.Context, operatorID string, req IdentityRevokeParam) error {
	row, err := s.repo.GetIdentity(ctx, req.AccountID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return bizErr(404, 404, "Verified identity not found")
		}
		return err
	}
	if row.Status != StatusVerified {
		return bizErr(404, 404, "Verified identity not found")
	}
	now := time.Now().UTC()
	row.Status = StatusRevoked
	row.RevokedAt = &now
	row.RevokedBy = &operatorID
	row.UpdatedBy = &operatorID
	return s.repo.SaveIdentity(ctx, row)
}

func (s *Service) validateSubmit(ctx context.Context, accountID, businessType, documentType, documentNo string) error {
	if businessType != BusinessAccountVerify {
		return bizErr(400, 400, "Unsupported business_type: "+businessType)
	}
	identity, err := s.repo.GetIdentity(ctx, accountID)
	if err == nil && identity != nil && identity.Status == StatusVerified {
		return bizErr(400, 400, "账号已完成实名认证")
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	pending, err := s.repo.CountPendingCases(ctx, accountID, businessType)
	if err != nil {
		return err
	}
	if pending > 0 {
		return bizErr(400, 400, "已有进行中的实名认证申请")
	}
	return s.assertDocumentAvailable(ctx, documentType, documentNo, accountID)
}

func (s *Service) assertDocumentAvailable(ctx context.Context, documentType, documentNo, excludeAccountID string) error {
	if strings.TrimSpace(documentNo) == "" {
		return bizErr(400, 400, "证件号码不能为空")
	}
	hash := s.crypto.HashDocumentNo(documentType, documentNo)
	if _, err := s.repo.FindVerifiedByDocumentHash(ctx, hash, excludeAccountID); err == nil {
		return bizErr(400, 400, "该证件已被其他账号绑定")
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if _, err := s.repo.FindPendingCaseByDocumentHash(ctx, hash, excludeAccountID); err == nil {
		return bizErr(400, 400, "该证件已有进行中的认证申请")
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return nil
}

func (s *Service) newCaseEntity(accountID, businessType, channel, documentType, realName, documentNo string) (*RealNameCase, error) {
	docType := strings.ToUpper(strings.TrimSpace(documentType))
	nameCipher, err := s.crypto.Encrypt(realName)
	if err != nil {
		return nil, err
	}
	noCipher, err := s.crypto.Encrypt(documentNo)
	if err != nil {
		return nil, err
	}
	hash := s.crypto.HashDocumentNo(docType, documentNo)
	acc := accountID
	return &RealNameCase{
		CaseID:           idgen.Next(),
		BusinessType:     businessType,
		VerifyChannel:    channel,
		Status:           CaseStatusPending,
		AccountID:        &acc,
		DocumentType:     &docType,
		RealNameCipher:   &nameCipher,
		DocumentNoCipher: &noCipher,
		DocumentNoHash:   &hash,
		AttachmentIDs:    StringList{},
	}, nil
}

func (s *Service) upsertOnApprove(ctx context.Context, caseEntity *RealNameCase) error {
	if caseEntity.AccountID == nil {
		return bizErr(400, 400, "case account_id required")
	}
	accountID := *caseEntity.AccountID
	now := time.Now().UTC()
	row, err := s.repo.GetIdentity(ctx, accountID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if row == nil {
		row = &ProfileIdentity{AccountID: accountID, CreatedBy: &accountID}
	}
	row.Status = StatusVerified
	row.DocumentType = caseEntity.DocumentType
	row.RealNameCipher = caseEntity.RealNameCipher
	row.DocumentNoCipher = caseEntity.DocumentNoCipher
	row.DocumentNoHash = caseEntity.DocumentNoHash
	row.VerifyChannel = &caseEntity.VerifyChannel
	row.Provider = caseEntity.Provider
	row.ProviderOrderNo = caseEntity.ProviderOrderNo
	row.VerifiedAt = &now
	caseID := caseEntity.CaseID
	row.SourceCaseID = &caseID
	row.RevokedAt = nil
	row.RevokedBy = nil
	row.UpdatedBy = &accountID
	return s.repo.SaveIdentity(ctx, row)
}

func (s *Service) requireCase(ctx context.Context, caseID string) (*RealNameCase, error) {
	row, err := s.repo.GetCase(ctx, caseID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, bizErr(404, 404, "Real-name case not found")
		}
		return nil, err
	}
	return row, nil
}

func (s *Service) toSummary(entity *RealNameCase) RealNameCaseSummaryResult {
	out := RealNameCaseSummaryResult{
		CaseID:        entity.CaseID,
		AccountID:     entity.AccountID,
		BusinessType:  entity.BusinessType,
		VerifyChannel: entity.VerifyChannel,
		Status:        entity.Status,
		DocumentType:  entity.DocumentType,
		CreatedAt:     entity.CreatedAt,
		ReviewedAt:    entity.ReviewedAt,
		RejectReason:  entity.RejectReason,
	}
	if entity.RealNameCipher != nil && *entity.RealNameCipher != "" {
		if plain, err := s.crypto.Decrypt(*entity.RealNameCipher); err == nil {
			out.RealNameMasked = s.crypto.MaskRealName(plain)
		}
	}
	if entity.DocumentNoCipher != nil && *entity.DocumentNoCipher != "" {
		if plain, err := s.crypto.Decrypt(*entity.DocumentNoCipher); err == nil {
			out.DocumentNoMasked = s.crypto.MaskDocumentNo(plain)
		}
	}
	return out
}

func (s *Service) toPageResult(row *ProfileIdentity) IdentityPageResult {
	out := IdentityPageResult{
		AccountID:     row.AccountID,
		Status:        row.Status,
		DocumentType:  row.DocumentType,
		VerifyChannel: row.VerifyChannel,
		Provider:      row.Provider,
		VerifiedAt:    row.VerifiedAt,
		RevokedAt:     row.RevokedAt,
	}
	if row.RealNameCipher != nil && *row.RealNameCipher != "" {
		if plain, err := s.crypto.Decrypt(*row.RealNameCipher); err == nil {
			out.RealNameMasked = s.crypto.MaskRealName(plain)
		}
	}
	if row.DocumentNoCipher != nil && *row.DocumentNoCipher != "" {
		if plain, err := s.crypto.Decrypt(*row.DocumentNoCipher); err == nil {
			out.DocumentNoMasked = s.crypto.MaskDocumentNo(plain)
		}
	}
	return out
}

func (s *Service) normalizeAttachmentIDs(raw []string) []string {
	names := make([]string, 0, len(raw))
	seen := map[string]struct{}{}
	for _, v := range raw {
		n := file.NormalizeObjectName(v)
		if n == "" {
			continue
		}
		if _, ok := seen[n]; ok {
			continue
		}
		seen[n] = struct{}{}
		names = append(names, n)
	}
	return names
}

func (s *Service) resolveAttachments(ctx context.Context, attachmentIDs StringList) []RealNameCaseAttachmentResult {
	if len(attachmentIDs) == 0 {
		return []RealNameCaseAttachmentResult{}
	}
	files, _ := s.repo.ListFilesByObjectNames(ctx, attachmentIDs)
	fileMap := make(map[string]file.File, len(files))
	for i := range files {
		fileMap[files[i].ObjectName] = files[i]
	}
	out := make([]RealNameCaseAttachmentResult, 0, len(attachmentIDs))
	for _, name := range attachmentIDs {
		item := RealNameCaseAttachmentResult{ObjectName: name}
		if f, ok := fileMap[name]; ok {
			item.ID = &f.ID
			item.OriginalName = &f.OriginalName
			item.ContentType = &f.ContentType
			item.Size = &f.Size
		}
		if s.storage != nil {
			if u := s.storage.ResolveURL(ctx, name); u != "" {
				item.URL = &u
			}
		}
		out = append(out, item)
	}
	return out
}

func sanitizeStatus(source *IdentityStatusResult) *IdentityStatusResult {
	if source == nil {
		return nil
	}
	source.VerifyChannel = nil
	source.Provider = nil
	if source.PendingCase != nil {
		summary := sanitizeSummary(*source.PendingCase)
		source.PendingCase = &summary
	}
	return source
}

func sanitizeSummary(summary RealNameCaseSummaryResult) RealNameCaseSummaryResult {
	summary.VerifyChannel = ""
	summary.RealNameMasked = nil
	summary.DocumentNoMasked = nil
	return summary
}

func normalizeBusinessType(businessType string) string {
	bt := strings.TrimSpace(businessType)
	if bt == "" {
		return BusinessAccountVerify
	}
	return strings.ToUpper(bt)
}

func deref(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}
