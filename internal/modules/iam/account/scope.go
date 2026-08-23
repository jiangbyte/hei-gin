// internal/modules/iam/account/scope.go 账号数据范围校验（对齐 hei-boot DataScopeResolver.assertAccountAccessible）。
//
// Author: Charlie

package account

import (
	"context"

	contextx "hei-gin/internal/framework/core/context"
	"hei-gin/internal/framework/core/security/datascope"
)

const accountPagePerm = "iam:account:page"

func (s *Service) assertAccountAccessible(ctx context.Context, accountID string) error {
	sess := contextx.Session(ctx)
	if err := datascope.AssertAccountAccessibleMsg(ctx, s.repo.DB(), sess, accountID, accountPagePerm); err != nil {
		return err
	}
	return nil
}