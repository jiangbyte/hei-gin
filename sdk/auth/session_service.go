package auth

import "context"

type RealmSessionService struct {
	realm *Realm
}

type AggregateSessionService struct {
	realms []*Realm
}

func (s *RealmSessionService) Page(ctx context.Context, current, size int, keyword string) ([]SessionInfo, int64, error) {
	if s == nil || s.realm == nil {
		return []SessionInfo{}, 0, nil
	}
	return s.realm.tool.listSessionInfos(ctx, current, size, keyword)
}

func (s *RealmSessionService) PageByUserIDs(ctx context.Context, userIDs []string, current, size int) ([]SessionInfo, int64, error) {
	if s == nil || s.realm == nil {
		return []SessionInfo{}, 0, nil
	}
	return s.realm.tool.listSessionInfosByUserIDs(ctx, userIDs, current, size)
}

func (s *RealmSessionService) Tokens(ctx context.Context, userID string) ([]SessionTokenInfo, error) {
	if s == nil || s.realm == nil {
		return []SessionTokenInfo{}, nil
	}
	return s.realm.tool.sessionTokens(ctx, userID)
}

func (s *RealmSessionService) Stats(ctx context.Context) (SessionStats, error) {
	if s == nil || s.realm == nil {
		return SessionStats{}, nil
	}
	return s.realm.tool.sessionStats(ctx)
}

func (s *RealmSessionService) Trend(ctx context.Context, days []string) map[string]int {
	if s == nil || s.realm == nil {
		return map[string]int{}
	}
	return s.realm.tool.sessionDailyCounts(ctx, days)
}

func (s *RealmSessionService) KickoutUser(ctx context.Context, userID string) {
	if s == nil || s.realm == nil {
		return
	}
	s.realm.KickoutWithContext(ctx, userID)
}

func (s *RealmSessionService) KickoutToken(ctx context.Context, userID, token string) {
	if s == nil || s.realm == nil {
		return
	}
	s.realm.KickoutTokenWithContext(ctx, userID, token)
}

func Sessions(realms ...*Realm) *AggregateSessionService {
	return &AggregateSessionService{realms: realms}
}

func (s *AggregateSessionService) Stats(ctx context.Context) (SessionStats, error) {
	result := SessionStats{}
	for _, realm := range s.realms {
		if realm == nil {
			continue
		}
		stats, err := realm.Sessions().Stats(ctx)
		if err != nil {
			return SessionStats{}, err
		}
		result.TotalCount += stats.TotalCount
		result.OneHourNewlyAdded += stats.OneHourNewlyAdded
		if stats.MaxTokenCount > result.MaxTokenCount {
			result.MaxTokenCount = stats.MaxTokenCount
		}
	}
	return result, nil
}
