package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"sort"
	"time"

	"hei-gin/sdk/config"
	"hei-gin/sdk/infra/db"
	"hei-gin/sdk/shared/contracts"

	"github.com/gin-gonic/gin"

	"github.com/redis/go-redis/v9"
)

// baseAuthTool is the shared authentication implementation for both BUSINESS and CONSUMER.
// It is parameterized by loginType and generates Redis key prefixes accordingly.
type baseAuthTool struct {
	expire    int
	tokenName string

	realmID RealmID
}

func newBaseAuthTool(realmID RealmID) *baseAuthTool {
	t := &baseAuthTool{realmID: realmID}
	t.ensureConfig()
	return t
}

// ensureConfig initializes default values from the global config if not already set.
func (t *baseAuthTool) ensureConfig() {
	if config.C == nil {
		return
	}
	t.expire = config.C.Token.ExpireSeconds
	t.tokenName = config.C.Token.TokenName

}

func (t *baseAuthTool) tokenURLSafe(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func (t *baseAuthTool) getRedis() *redis.Client {
	return db.Redis
}

func (t *baseAuthTool) getTokenKey(token string) string {
	return "hei:auth:" + string(t.realmID) + ":token:" + token
}

func (t *baseAuthTool) getSessionKey(userID string) string {
	return "hei:auth:" + string(t.realmID) + ":session:" + userID
}

func (t *baseAuthTool) getDisableKey(loginID string) string {
	return "hei:auth:" + string(t.realmID) + ":disable:" + loginID
}

func requestContext(c *gin.Context) context.Context {
	if c != nil && c.Request != nil {
		return c.Request.Context()
	}
	return context.Background()
}

// Init sets custom expire and token name. Falls back to config if values are zero/empty.
func (t *baseAuthTool) Init(expire int, tokenName string) {
	t.ensureConfig()
	if expire > 0 {
		t.expire = expire
	}
	if tokenName != "" {
		t.tokenName = tokenName
	}
}

// GetLoginType returns the login type identifier.
func (t *baseAuthTool) GetLoginType() string {
	t.ensureConfig()
	return string(t.realmID)
}

// GetTokenName returns the HTTP header name used to carry the token.
func (t *baseAuthTool) GetTokenName() string {
	t.ensureConfig()
	return t.tokenName
}

// GetTokenValue extracts the token from the request header.
func (t *baseAuthTool) GetTokenValue(c *gin.Context) string {
	t.ensureConfig()
	if c == nil {
		return ""
	}
	return c.GetHeader(t.tokenName)
}

// Login authenticates a user by user ID, stores token data in Redis, and returns the token.
func (t *baseAuthTool) Login(c *gin.Context, id string, extra map[string]any) (string, error) {
	ctx := requestContext(c)
	t.ensureConfig()

	now := time.Now()
	signedToken := t.tokenURLSafe(32)
	claims, err := t.buildClaims(ctx, id, now, extra)
	if err != nil {
		return "", err
	}
	tokenDataJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	redisClient := t.getRedis()

	err = redisClient.SetEx(ctx, t.getTokenKey(signedToken), tokenDataJSON, time.Duration(t.expire)*time.Second).Err()
	if err != nil {
		return "", err
	}

	sessionKey := t.getSessionKey(id)

	// Clean expired tokens from the session set
	existingTokens, _ := redisClient.SMembers(ctx, sessionKey).Result()
	for _, existingToken := range existingTokens {
		if _, ok := t.getClaimsWithContext(ctx, existingToken); !ok {
			_ = redisClient.SRem(ctx, sessionKey, existingToken).Err()
		}
	}

	err = redisClient.SAdd(ctx, sessionKey, signedToken).Err()
	if err != nil {
		return "", err
	}
	err = redisClient.Expire(ctx, sessionKey, time.Duration(t.expire)*time.Second).Err()
	if err != nil {
		return "", err
	}

	t.trackLoginSession(ctx, id, signedToken, now, time.Duration(t.expire)*time.Second)

	return signedToken, nil
}

func (t *baseAuthTool) buildClaims(ctx context.Context, userID string, now time.Time, extra map[string]any) (*SessionClaims, error) {
	if extra == nil {
		extra = map[string]any{}
	}
	acl, err := t.loadACL(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &SessionClaims{
		UserID:    userID,
		RealmID:   t.realmID,
		CreatedAt: now.Format("2006-01-02 15:04:05"),
		Extra:     extra,
		ACL:       acl,
	}, nil
}

func (t *baseAuthTool) loadACL(ctx context.Context, userID string) (ACLSnapshot, error) {
	if PermissionDelegate == nil || userID == "" {
		return ACLSnapshot{
			Permissions: []string{},
			Roles:       []string{},
			ScopeMap:    map[string]ScopeInfo{},
		}, nil
	}

	perms, err := PermissionDelegate.GetPermissionList(ctx, contracts.RealmID(t.realmID), userID)
	if err != nil {
		return ACLSnapshot{}, err
	}
	roles, err := PermissionDelegate.GetRoleList(ctx, contracts.RealmID(t.realmID), userID)
	if err != nil {
		return ACLSnapshot{}, err
	}
	scopeMap, err := PermissionDelegate.GetPermissionScopeMap(ctx, contracts.RealmID(t.realmID), userID)
	if err != nil {
		return ACLSnapshot{}, err
	}

	sort.Strings(perms)
	sort.Strings(roles)
	if scopeMap == nil {
		scopeMap = map[string]ScopeInfo{}
	}

	return ACLSnapshot{
		Permissions: perms,
		Roles:       roles,
		ScopeMap:    scopeMap,
	}, nil
}

// Logout invalidates the current session. If loginID is provided, it kicks out all sessions for that user.
func (t *baseAuthTool) Logout(c *gin.Context, loginID ...string) {
	ctx := requestContext(c)
	t.ensureConfig()

	if len(loginID) > 0 {
		t.kickoutWithContext(ctx, loginID[0])
		return
	}

	token := t.GetTokenValue(c)
	if token == "" {
		return
	}

	claims, ok := t.getClaimsWithContext(ctx, token)
	userID := ""
	if ok {
		userID = claims.UserID
	}
	if userID != "" {
		redisClient := t.getRedis()
		sessionKey := t.getSessionKey(userID)
		_ = redisClient.SRem(ctx, sessionKey, token).Err()
	}

	redisClient := t.getRedis()
	tokenKey := t.getTokenKey(token)
	_ = redisClient.Del(ctx, tokenKey).Err()
	t.untrackToken(ctx, userID, token)
}

// Kickout deletes all tokens and session data for the given login ID.
func (t *baseAuthTool) Kickout(loginID string) {
	t.kickoutWithContext(context.Background(), loginID)
}

func (t *baseAuthTool) KickoutWithContext(ctx context.Context, loginID string) {
	t.kickoutWithContext(ctx, loginID)
}

func (t *baseAuthTool) kickoutWithContext(ctx context.Context, loginID string) {
	t.ensureConfig()

	redisClient := t.getRedis()
	sessionKey := t.getSessionKey(loginID)

	tokens, err := redisClient.SMembers(ctx, sessionKey).Result()
	if err != nil {
		return
	}

	if len(tokens) > 0 {
		pipe := redisClient.Pipeline()
		for _, token := range tokens {
			pipe.Del(ctx, t.getTokenKey(token))
		}
		_, _ = pipe.Exec(ctx)
		t.removeTokenIndexes(ctx, tokens)
	}

	_ = redisClient.Del(ctx, sessionKey).Err()
	t.removeSessionIndexes(ctx, loginID)
}

// KickoutToken removes a specific token from the user's session set and deletes its data.
func (t *baseAuthTool) KickoutToken(loginID, token string) {
	t.KickoutTokenWithContext(context.Background(), loginID, token)
}

func (t *baseAuthTool) KickoutTokenWithContext(ctx context.Context, loginID, token string) {
	t.ensureConfig()

	redisClient := t.getRedis()
	sessionKey := t.getSessionKey(loginID)
	tokenKey := t.getTokenKey(token)
	if loginID == "" || token == "" {
		return
	}
	owner, ownerErr := redisClient.HGet(ctx, t.getTokenOwnerKey(), token).Result()
	if ownerErr == nil {
		if owner != loginID {
			return
		}
	} else if ownerErr == redis.Nil {
		if isMember, err := redisClient.SIsMember(ctx, sessionKey, token).Result(); err != nil || !isMember {
			return
		}
	} else {
		return
	}

	_ = redisClient.SRem(ctx, sessionKey, token).Err()
	_ = redisClient.Del(ctx, tokenKey).Err()
	t.untrackToken(ctx, loginID, token)
}

// IsLogin checks whether the current request carries a valid token.
func (t *baseAuthTool) IsLogin(c *gin.Context) bool {
	loginID := t.GetLoginIDDefaultNull(c)
	return loginID != ""
}

// CheckLogin returns an error if the current request is not authenticated.
func (t *baseAuthTool) CheckLogin(c *gin.Context) error {
	if !t.IsLogin(c) {
		return errors.New("未授权/未登录")
	}
	return nil
}

// GetLoginID extracts and returns the login ID from the current request's token.
func (t *baseAuthTool) GetLoginID(c *gin.Context) string {
	return t.GetLoginIDDefaultNull(c)
}

// GetLoginIDDefaultNull returns the login ID from the token, or empty string if not logged in.
func (t *baseAuthTool) GetLoginIDDefaultNull(c *gin.Context) string {
	token := t.GetTokenValue(c)
	if token == "" {
		return ""
	}
	claims, ok := t.decodeClaims(c, token)
	if !ok {
		return ""
	}
	return claims.UserID
}

// GetLoginIDByToken extracts the login ID from the given token value.
func (t *baseAuthTool) GetLoginIDByToken(token string) string {
	if token == "" {
		return ""
	}
	claims, ok := t.getClaims(token)
	if !ok {
		return ""
	}
	return claims.UserID
}

func (t *baseAuthTool) decodeClaims(c *gin.Context, token string) (*SessionClaims, bool) {
	if token == "" {
		return nil, false
	}

	return t.getClaimsForRequest(c, token)
}

func (t *baseAuthTool) getClaims(token string) (*SessionClaims, bool) {
	return t.getClaimsWithContext(context.Background(), token)
}

func (t *baseAuthTool) getClaimsForRequest(c *gin.Context, token string) (*SessionClaims, bool) {
	if c == nil {
		return t.getClaims(token)
	}
	cacheKey := "_auth_claims:" + string(t.realmID) + ":" + token
	if cached, exists := c.Get(cacheKey); exists {
		if claims, ok := cached.(*SessionClaims); ok {
			return claims, true
		}
	}
	claims, ok := t.getClaimsWithContext(requestContext(c), token)
	if ok {
		c.Set(cacheKey, claims)
	}
	return claims, ok
}

func (t *baseAuthTool) getClaimsWithContext(ctx context.Context, token string) (*SessionClaims, bool) {
	if token == "" {
		return nil, false
	}

	redisClient := t.getRedis()
	tokenKey := t.getTokenKey(token)

	data, err := redisClient.Get(ctx, tokenKey).Result()
	if err == redis.Nil {
		return nil, false
	}
	if err != nil {
		return nil, false
	}

	var result SessionClaims
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		return nil, false
	}
	if result.Extra == nil {
		result.Extra = map[string]any{}
	}
	if result.ACL.ScopeMap == nil {
		result.ACL.ScopeMap = map[string]ScopeInfo{}
	}
	return &result, true
}

// GetExtra returns a specific extra field from the token data.
func (t *baseAuthTool) GetExtra(c *gin.Context, key string) any {
	claims, ok := t.CurrentClaims(c)
	if ok {
		return claims.Extra[key]
	}
	return nil
}

// GetSession returns the full token payload for the current request.
func (t *baseAuthTool) GetSession(c *gin.Context) map[string]any {
	token := t.GetTokenValue(c)
	if token == "" {
		return nil
	}
	claims, ok := t.getClaimsForRequest(c, token)
	if !ok {
		return nil
	}
	return claimsToMap(claims)
}

func (t *baseAuthTool) CurrentClaims(c *gin.Context) (*SessionClaims, bool) {
	token := t.GetTokenValue(c)
	if token == "" {
		return nil, false
	}
	return t.getClaimsForRequest(c, token)
}

func (t *baseAuthTool) GetClaims(c *gin.Context) (*SessionClaims, bool) {
	return t.CurrentClaims(c)
}

func (t *baseAuthTool) refreshUserSessionsACL(ctx context.Context, userID string) error {
	if userID == "" {
		return nil
	}
	redisClient := t.getRedis()
	if redisClient == nil {
		return nil
	}

	tokens, err := redisClient.SMembers(ctx, t.getSessionKey(userID)).Result()
	if err != nil || len(tokens) == 0 {
		return err
	}

	acl, err := t.loadACL(ctx, userID)
	if err != nil {
		return err
	}

	staleTokens := make([]string, 0)
	pipe := redisClient.Pipeline()
	for _, token := range tokens {
		claims, ok := t.getClaimsWithContext(ctx, token)
		if !ok {
			staleTokens = append(staleTokens, token)
			continue
		}
		claims.ACL = acl
		payload, marshalErr := json.Marshal(claims)
		if marshalErr != nil {
			return marshalErr
		}
		ttl, ttlErr := redisClient.TTL(ctx, t.getTokenKey(token)).Result()
		if ttlErr != nil || ttl <= 0 {
			staleTokens = append(staleTokens, token)
			continue
		}
		pipe.Set(ctx, t.getTokenKey(token), payload, ttl)
	}
	if _, err := pipe.Exec(ctx); err != nil {
		return err
	}
	if len(staleTokens) > 0 {
		t.removeStaleTokens(ctx, userID, staleTokens)
	}
	t.refreshSessionIndexes(ctx, userID)
	return nil
}

// RenewTimeout extends the token and session timeouts.
func (t *baseAuthTool) RenewTimeout(c *gin.Context, timeout ...int) {
	t.ensureConfig()

	token := t.GetTokenValue(c)
	if token == "" {
		return
	}

	newTimeout := t.expire
	if len(timeout) > 0 && timeout[0] > 0 {
		newTimeout = timeout[0]
	}

	redisClient := t.getRedis()
	ctx := requestContext(c)
	tokenKey := t.getTokenKey(token)
	_ = redisClient.Expire(ctx, tokenKey, time.Duration(newTimeout)*time.Second).Err()

	loginID := t.GetLoginIDDefaultNull(c)
	if loginID != "" {
		sessionKey := t.getSessionKey(loginID)
		_ = redisClient.Expire(ctx, sessionKey, time.Duration(newTimeout)*time.Second).Err()
		t.updateTokenExpiryIndex(ctx, loginID, token, time.Duration(newTimeout)*time.Second)
	}
}

// GetTokenTimeout returns the remaining TTL (in seconds) of the current token. Returns -1 if not logged in.
func (t *baseAuthTool) GetTokenTimeout(c *gin.Context) int {
	token := t.GetTokenValue(c)
	if token == "" {
		return -1
	}

	redisClient := t.getRedis()
	ctx := requestContext(c)
	tokenKey := t.getTokenKey(token)

	ttl, err := redisClient.TTL(ctx, tokenKey).Result()
	if err != nil || ttl < 0 {
		return -1
	}
	return int(ttl.Seconds())
}

// GetSessionTimeout returns the remaining TTL (in seconds) of the current session. Returns -1 if not logged in.
func (t *baseAuthTool) GetSessionTimeout(c *gin.Context) int {
	loginID := t.GetLoginIDDefaultNull(c)
	if loginID == "" {
		return -1
	}

	redisClient := t.getRedis()
	ctx := requestContext(c)
	sessionKey := t.getSessionKey(loginID)

	ttl, err := redisClient.TTL(ctx, sessionKey).Result()
	if err != nil || ttl < 0 {
		return -1
	}
	return int(ttl.Seconds())
}

// GetTokenValueByLoginID returns one token for the given login ID.
func (t *baseAuthTool) GetTokenValueByLoginID(loginID string) string {
	redisClient := t.getRedis()
	ctx := context.Background()
	sessionKey := t.getSessionKey(loginID)

	tokens, err := redisClient.SMembers(ctx, sessionKey).Result()
	if err != nil || len(tokens) == 0 {
		return ""
	}
	return tokens[0]
}

// GetTokenValuesByLoginID returns all tokens for the given login ID.
func (t *baseAuthTool) GetTokenValuesByLoginID(loginID string) []string {
	redisClient := t.getRedis()
	ctx := context.Background()
	sessionKey := t.getSessionKey(loginID)

	tokens, err := redisClient.SMembers(ctx, sessionKey).Result()
	if err != nil {
		return nil
	}
	return tokens
}

// Disable marks a login ID as disabled for the specified duration (in seconds).
func (t *baseAuthTool) Disable(loginID string, timeSeconds int) {
	redisClient := t.getRedis()
	ctx := context.Background()
	disableKey := t.getDisableKey(loginID)
	_ = redisClient.SetEx(ctx, disableKey, "1", time.Duration(timeSeconds)*time.Second).Err()
}

// IsDisable checks whether a login ID is currently disabled.
func (t *baseAuthTool) IsDisable(loginID string) bool {
	redisClient := t.getRedis()
	ctx := context.Background()
	disableKey := t.getDisableKey(loginID)

	exists, err := redisClient.Exists(ctx, disableKey).Result()
	if err != nil {
		return false
	}
	return exists > 0
}

// CheckDisable returns an error if the login ID is currently disabled.
func (t *baseAuthTool) CheckDisable(loginID string) error {
	if t.IsDisable(loginID) {
		return errors.New("账号已被禁用")
	}
	return nil
}

// GetDisableTime returns the remaining disable time (in seconds). Returns -1 if not disabled.
func (t *baseAuthTool) GetDisableTime(loginID string) int {
	redisClient := t.getRedis()
	ctx := context.Background()
	disableKey := t.getDisableKey(loginID)

	ttl, err := redisClient.TTL(ctx, disableKey).Result()
	if err != nil || ttl < 0 {
		return -1
	}
	return int(ttl.Seconds())
}

// UntieDisable removes the disabled status from a login ID.
func (t *baseAuthTool) UntieDisable(loginID string) {
	redisClient := t.getRedis()
	ctx := context.Background()
	disableKey := t.getDisableKey(loginID)
	_ = redisClient.Del(ctx, disableKey).Err()
}

func claimsToMap(claims *SessionClaims) map[string]any {
	if claims == nil {
		return nil
	}
	return map[string]any{
		"user_id":    claims.UserID,
		"realm_id":   string(claims.RealmID),
		"created_at": claims.CreatedAt,
		"extra":      claims.Extra,
		"acl": map[string]any{
			"permissions": claims.ACL.Permissions,
			"roles":       claims.ACL.Roles,
			"scope_map":   claims.ACL.ScopeMap,
		},
	}
}
