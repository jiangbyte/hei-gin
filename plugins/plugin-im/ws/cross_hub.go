package ws

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"hei-gin/sdk/auth"
	"hei-gin/sdk/config"
	"hei-gin/sdk/infra/eventbus"

	"github.com/redis/go-redis/v9"
)

// ─── Config helpers ───────────────────────────────────────────────────

func hbInterval() time.Duration {
	if getConfig().HeartbeatInterval > 0 {
		return time.Duration(getConfig().HeartbeatInterval) * time.Second
	}
	return 15 * time.Second
}

func instTTL() time.Duration {
	if getConfig().InstanceTTL > 0 {
		return time.Duration(getConfig().InstanceTTL) * time.Second
	}
	return 60 * time.Second
}

func staleClean() time.Duration {
	if getConfig().StaleCleanInterval > 0 {
		return time.Duration(getConfig().StaleCleanInterval) * time.Minute
	}
	return 5 * time.Minute
}

func rlWindow() time.Duration {
	if getConfig().RateLimitWindow > 0 {
		return time.Duration(getConfig().RateLimitWindow) * time.Second
	}
	return 10 * time.Second
}

func rlMax() int64 {
	if getConfig().RateLimitMax > 0 {
		return int64(getConfig().RateLimitMax)
	}
	return 30
}

func dedupTTL() time.Duration {
	if getConfig().DedupTTL > 0 {
		return time.Duration(getConfig().DedupTTL) * time.Second
	}
	return 30 * time.Second
}

func pollTO() time.Duration {
	if getConfig().PollTimeout > 0 {
		return time.Duration(getConfig().PollTimeout) * time.Second
	}
	return 2 * time.Second
}

// msgListTTL is the TTL for per-instance message lists, preventing stale message accumulation.
const msgListTTL = 5 * time.Minute

// ─── Types ────────────────────────────────────────────────────────────

type crossInstanceMessage struct {
	ToUserID   string          `json:"to_user_id"`
	ToUserType string          `json:"to_user_type"`
	Message    json.RawMessage `json:"message"`
	MessageID  string          `json:"message_id,omitempty"`
	Timestamp  int64           `json:"timestamp"`
}

// CrossHub wraps the local Hub with Redis-backed cross-instance IM delivery.
type CrossHub struct {
	local      *Hub
	rdb        *redis.Client
	instanceID string
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
	closeOnce  sync.Once
}

// NewCrossHub creates a CrossHub. If rdb is nil, runs in single-instance mode.
func NewCrossHub(local *Hub, rdb *redis.Client) *CrossHub {
	instanceID := strconv.FormatInt(config.C.Snowflake.Instance, 10)

	ctx, cancel := context.WithCancel(context.Background())
	ch := &CrossHub{
		local:      local,
		rdb:        rdb,
		instanceID: instanceID,
		ctx:        ctx,
		cancel:     cancel,
	}

	local.OnClientRegistered = ch.onClientRegistered
	local.OnClientUnregistered = ch.onClientUnregistered

	if rdb != nil {
		ch.wg.Add(4)
		go ch.pollLoop()
		go ch.heartbeatLoop()
		go ch.staleCleanupLoop()
		go ch.msgListCleanupLoop()
		log.Printf("[CrossHub] Cross-instance mode enabled, instance=%s", instanceID)
	} else {
		log.Printf("[CrossHub] Redis not configured, single-instance mode")
	}

	return ch
}

// ─── Hub lifecycle hooks ──────────────────────────────────────────────

func (ch *CrossHub) onClientRegistered(c *Client) {
	ch.TrackConnection(c.UserID, c.UserType)
	ch.broadcastPresence(c.UserID, c.UserType, true)
	eventbus.DefaultBus.Publish(eventbus.TopicUserConnected, c)
	switch c.UserType {
	case auth.BusinessID:
		ch.local.SendToUser(c.UserID, Message{Type: MsgUnreadCount})
	case auth.ConsumerID:
		ch.local.SendToConsumer(c.UserID, Message{Type: MsgUnreadCount})
	}
}

func (ch *CrossHub) onClientUnregistered(c *Client) {
	ch.UntrackConnection(c.UserID, c.UserType)
	if !ch.IsUserOnlineAnywhere(c.UserID, c.UserType) {
		ch.broadcastPresence(c.UserID, c.UserType, false)
	}
	eventbus.DefaultBus.Publish(eventbus.TopicUserDisconnected, c)
}

// ─── Redis key helpers ────────────────────────────────────────────────

func (ch *CrossHub) userSetKey(userType auth.RealmID, userID string) string {
	return "ws:user:" + string(userType) + ":" + userID
}

func (ch *CrossHub) msgListKey() string {
	return "ws:messages:" + ch.instanceID
}

func (ch *CrossHub) instanceKey() string {
	return "ws:instance:" + ch.instanceID
}

func (ch *CrossHub) rateLimitKey(userID string, userType auth.RealmID) string {
	return "ws:ratelimit:" + string(userType) + ":" + userID
}

func (ch *CrossHub) userCountKey(userType auth.RealmID, userID string) string {
	return "ws:usercnt:" + string(userType) + ":" + userID
}

func (ch *CrossHub) dedupKey(messageID string) string {
	return "ws:dedup:" + ch.instanceID + ":" + messageID
}

// ─── Presence ─────────────────────────────────────────────────────────

func (ch *CrossHub) IsUserOnlineAnywhere(userID string, userType auth.RealmID) bool {
	if ch == nil {
		return false
	}

	if ch.rdb == nil {
		return ch.local.isUserConnected(userID, userType)
	}
	key := ch.userSetKey(userType, userID)
	count, err := ch.rdb.SCard(ch.ctx, key).Result()
	if err != nil {
		return false
	}
	return count > 0
}

func (ch *CrossHub) IsUserOnlineLocally(userID string, userType auth.RealmID) bool {
	return ch.local.isUserConnected(userID, userType)
}

func (ch *CrossHub) broadcastPresence(userID string, userType auth.RealmID, online bool) {
	msg := Message{
		Type: MsgPresence,
		Payload: PresencePayload{
			UserID:   userID,
			UserType: string(userType),
			Online:   online,
		},
	}
	ch.local.BroadcastAll(msg)
}

func (ch *CrossHub) TrackConnection(userID string, userType auth.RealmID) {
	if ch.rdb == nil {
		return
	}
	key := ch.userSetKey(userType, userID)
	countKey := ch.userCountKey(userType, userID)
	if err := ch.rdb.SAdd(ch.ctx, key, ch.instanceID).Err(); err != nil {
		log.Printf("[CrossHub] TrackConnection SAdd error: %v", err)
	}
	if err := ch.rdb.HIncrBy(ch.ctx, countKey, ch.instanceID, 1).Err(); err != nil {
		log.Printf("[CrossHub] TrackConnection HIncrBy error: %v", err)
	}
	ch.rdb.Expire(ch.ctx, key, instTTL()+30*time.Second)
	ch.rdb.Expire(ch.ctx, countKey, instTTL()+30*time.Second)
}

func (ch *CrossHub) UntrackConnection(userID string, userType auth.RealmID) {
	if ch.rdb == nil {
		return
	}
	key := ch.userSetKey(userType, userID)
	countKey := ch.userCountKey(userType, userID)
	count, err := ch.rdb.HIncrBy(ch.ctx, countKey, ch.instanceID, -1).Result()
	if err != nil {
		log.Printf("[CrossHub] UntrackConnection HIncrBy error: %v", err)
		count = 0
	}
	if count <= 0 {
		_ = ch.rdb.HDel(ch.ctx, countKey, ch.instanceID).Err()
		if err := ch.rdb.SRem(ch.ctx, key, ch.instanceID).Err(); err != nil {
			log.Printf("[CrossHub] UntrackConnection SRem error: %v", err)
		}
	}
	ch.rdb.Expire(ch.ctx, key, instTTL()+30*time.Second)
	ch.rdb.Expire(ch.ctx, countKey, instTTL()+30*time.Second)
}

func userCountKeyFromSetKey(key string) string {
	const prefix = "ws:user:"
	if len(key) <= len(prefix) || key[:len(prefix)] != prefix {
		return ""
	}
	return "ws:usercnt:" + key[len(prefix):]
}

func (ch *CrossHub) getTargetInstances(userID string, userType auth.RealmID) []string {
	if ch.rdb == nil {
		return nil
	}
	key := ch.userSetKey(userType, userID)
	instances, err := ch.rdb.SMembers(ch.ctx, key).Result()
	if err != nil {
		log.Printf("[CrossHub] SMembers error: %v", err)
		return nil
	}
	return instances
}

// ─── Rate Limiting ────────────────────────────────────────────────────

func (ch *CrossHub) AllowMessage(userID string, userType auth.RealmID) bool {
	if ch == nil {
		return false
	}

	if ch.rdb == nil {
		return true
	}
	key := ch.rateLimitKey(userID, userType)
	count, err := ch.rdb.Incr(ch.ctx, key).Result()
	if err != nil {
		return true
	}
	if count == 1 {
		ch.rdb.Expire(ch.ctx, key, rlWindow())
	}
	return count <= rlMax()
}

// ─── Message Deduplication ────────────────────────────────────────────

func (ch *CrossHub) markDelivered(messageID string) bool {
	if ch.rdb == nil || messageID == "" {
		return true
	}
	key := ch.dedupKey(messageID)
	ok, err := ch.rdb.SetNX(ch.ctx, key, ch.instanceID, dedupTTL()).Result()
	if err != nil {
		return true
	}
	return ok
}

// ─── Polling Loop ─────────────────────────────────────────────────────

func (ch *CrossHub) pollLoop() {
	defer ch.wg.Done()
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[CrossHub] pollLoop panicked: %v", r)
		}
	}()

	key := ch.msgListKey()
	for {
		result, err := ch.rdb.BRPop(ch.ctx, pollTO(), key).Result()
		if err != nil {
			if err == context.Canceled || err == redis.Nil {
				select {
				case <-ch.ctx.Done():
					return
				default:
					continue
				}
			}
			log.Printf("[CrossHub] BRPop error: %v", err)
			continue
		}
		if len(result) < 2 {
			continue
		}
		ch.handleMessage(result[1])
	}
}

func (ch *CrossHub) handleMessage(payload string) {
	// Panic recovery per message to avoid killing the poll loop
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[CrossHub] handleMessage panicked: %v, payload=%.120s", r, payload)
		}
	}()

	var xMsg crossInstanceMessage
	if err := json.Unmarshal([]byte(payload), &xMsg); err != nil {
		log.Printf("[CrossHub] Failed to unmarshal message: %v", err)
		return
	}

	if xMsg.MessageID != "" {
		if !ch.markDelivered(xMsg.MessageID) {
			return
		}
	}

	var msg Message
	if err := json.Unmarshal(xMsg.Message, &msg); err != nil {
		log.Printf("[CrossHub] Failed to unmarshal inner message: %v", err)
		return
	}

	switch auth.RealmID(xMsg.ToUserType) {
	case auth.BusinessID:
		ch.local.SendToUser(xMsg.ToUserID, msg)
	case auth.ConsumerID:
		ch.local.SendToConsumer(xMsg.ToUserID, msg)
	}
}

// ─── Heartbeat ────────────────────────────────────────────────────────

func (ch *CrossHub) heartbeatLoop() {
	defer ch.wg.Done()
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[CrossHub] heartbeatLoop panicked: %v", r)
		}
	}()

	ticker := time.NewTicker(hbInterval())
	defer ticker.Stop()

	ch.sendHeartbeat()

	for {
		select {
		case <-ch.ctx.Done():
			return
		case <-ticker.C:
			ch.sendHeartbeat()
		}
	}
}

func (ch *CrossHub) sendHeartbeat() {
	if ch.rdb == nil {
		return
	}
	key := ch.instanceKey()
	if err := ch.rdb.SetEx(ch.ctx, key, time.Now().UnixMilli(), instTTL()).Err(); err != nil {
		log.Printf("[CrossHub] Heartbeat error: %v", err)
	}
}

// ─── Stale Instance Cleanup ───────────────────────────────────────────

func (ch *CrossHub) staleCleanupLoop() {
	defer ch.wg.Done()
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[CrossHub] staleCleanupLoop panicked: %v", r)
		}
	}()

	ticker := time.NewTicker(staleClean())
	defer ticker.Stop()

	for {
		select {
		case <-ch.ctx.Done():
			return
		case <-ticker.C:
			ch.cleanStaleInstances()
		}
	}
}

func (ch *CrossHub) cleanStaleInstances() {
	if ch.rdb == nil {
		return
	}

	iter := ch.rdb.Scan(ch.ctx, 0, "ws:user:*", 1000).Iterator()
	cleaned := 0

	for iter.Next(ch.ctx) {
		key := iter.Val()
		members, err := ch.rdb.SMembers(ch.ctx, key).Result()
		if err != nil {
			continue
		}
		for _, instID := range members {
			if instID == ch.instanceID {
				continue
			}
			instKey := "ws:instance:" + instID
			exists, err := ch.rdb.Exists(ch.ctx, instKey).Result()
			if err != nil {
				continue
			}
			if exists == 0 {
				ch.rdb.SRem(ch.ctx, key, instID)
				if countKey := userCountKeyFromSetKey(key); countKey != "" {
					ch.rdb.HDel(ch.ctx, countKey, instID)
				}
				cleaned++
			}
		}
	}

	if cleaned > 0 {
		log.Printf("[CrossHub] Cleaned %d stale instance references", cleaned)
	}
}

// ─── Message List Cleanup ─────────────────────────────────────────────

// msgListCleanupLoop periodically trims stale messages from this instance's list.
// Prevents memory growth in Redis when the poll loop is slow.
func (ch *CrossHub) msgListCleanupLoop() {
	defer ch.wg.Done()
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[CrossHub] msgListCleanupLoop panicked: %v", r)
		}
	}()

	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ch.ctx.Done():
			return
		case <-ticker.C:
			key := ch.msgListKey()
			// Trim to last 1000 messages to prevent unbounded growth
			if err := ch.rdb.LTrim(ch.ctx, key, -1000, -1).Err(); err != nil {
				log.Printf("[CrossHub] LTrim error: %v", err)
			}
			// Set TTL so stale lists auto-expire
			ch.rdb.Expire(ch.ctx, key, msgListTTL)
		}
	}
}

// ─── Public API ───────────────────────────────────────────────────────

func (ch *CrossHub) SendToUsers(userIDs []string, msg Message) {
	if ch == nil {
		return
	}

	ch.local.SendToUsers(userIDs, msg)
	if ch.rdb != nil {
		for _, uid := range userIDs {
			ch.publishToRemote(uid, auth.BusinessID, msg, "")
		}
	}
}

func (ch *CrossHub) SendToConsumers(userIDs []string, msg Message) {
	if ch == nil {
		return
	}

	ch.local.SendToConsumers(userIDs, msg)
	if ch.rdb != nil {
		for _, uid := range userIDs {
			ch.publishToRemote(uid, auth.ConsumerID, msg, "")
		}
	}
}

func (ch *CrossHub) SendMessagesToUsers(messages map[string]Message, messageIDs map[string]string) {
	if ch == nil {
		return
	}

	ch.local.SendMessagesToUsers(messages)
	if ch.rdb == nil {
		return
	}
	for uid, msg := range messages {
		ch.publishToRemote(uid, auth.BusinessID, msg, messageIDs[uid])
	}
}

func (ch *CrossHub) SendMessagesToConsumers(messages map[string]Message, messageIDs map[string]string) {
	if ch == nil {
		return
	}

	ch.local.SendMessagesToConsumers(messages)
	if ch.rdb == nil {
		return
	}
	for uid, msg := range messages {
		ch.publishToRemote(uid, auth.ConsumerID, msg, messageIDs[uid])
	}
}

func (ch *CrossHub) SendToUser(userID string, msg Message, messageID ...string) {
	if ch == nil {
		return
	}
	ch.local.SendToUser(userID, msg)
	if ch.rdb != nil {
		mid := ""
		if len(messageID) > 0 {
			mid = messageID[0]
		}
		ch.publishToRemote(userID, auth.BusinessID, msg, mid)
	}
}

func (ch *CrossHub) SendToConsumer(userID string, msg Message, messageID ...string) {
	if ch == nil {
		return
	}
	ch.local.SendToConsumer(userID, msg)
	if ch.rdb != nil {
		mid := ""
		if len(messageID) > 0 {
			mid = messageID[0]
		}
		ch.publishToRemote(userID, auth.ConsumerID, msg, mid)
	}
}

// publishToRemote pushes a message to remote instances where the user is connected.
func (ch *CrossHub) publishToRemote(userID string, userType auth.RealmID, msg Message, messageID string) {
	instances := ch.getTargetInstances(userID, userType)
	if len(instances) == 0 {
		return
	}

	rawMsg, err := json.Marshal(msg)
	if err != nil {
		return
	}

	xMsg := crossInstanceMessage{
		ToUserID:   userID,
		ToUserType: string(userType),
		Message:    rawMsg,
		MessageID:  messageID,
		Timestamp:  time.Now().UnixMilli(),
	}

	data, err := json.Marshal(xMsg)
	if err != nil {
		return
	}

	for _, instID := range instances {
		if instID == ch.instanceID {
			continue
		}
		listKey := "ws:messages:" + instID
		if err := ch.rdb.LPush(ch.ctx, listKey, data).Err(); err != nil {
			log.Printf("[CrossHub] LPush to %s error: %v", instID, err)
		}
	}
}

func (ch *CrossHub) HandleWebSocket(w http.ResponseWriter, r *http.Request, userID string, userType auth.RealmID) {
	if ch == nil {
		return
	}

	ch.local.HandleWebSocket(w, r, userID, userType)
}

func (ch *CrossHub) OnlineCount() int {
	if ch == nil {
		return 0
	}

	return ch.local.OnlineCount()
}

func (ch *CrossHub) BroadcastAll(msg Message) {
	if ch == nil {
		return
	}

	ch.local.BroadcastAll(msg)
}

func (ch *CrossHub) BroadcastBusiness(msg Message) {
	if ch == nil {
		return
	}

	ch.local.BroadcastBusiness(msg)
}

func (ch *CrossHub) BroadcastConsumers(msg Message) {
	if ch == nil {
		return
	}

	ch.local.BroadcastConsumers(msg)
}

func (ch *CrossHub) Close() {
	if ch == nil {
		return
	}
	ch.closeOnce.Do(func() {
		if ch.rdb != nil {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			ch.removeInstancePresence(ctx, ch.instanceID)
			_ = ch.rdb.Del(ctx, ch.instanceKey()).Err()
			_ = ch.rdb.Del(ctx, ch.msgListKey()).Err()
			cancel()
		}
		ch.cancel()
		ch.wg.Wait()
	})
}

func (ch *CrossHub) removeInstancePresence(ctx context.Context, instanceID string) {
	iter := ch.rdb.Scan(ctx, 0, "ws:user:*", 1000).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		_ = ch.rdb.SRem(ctx, key, instanceID).Err()
		if countKey := userCountKeyFromSetKey(key); countKey != "" {
			_ = ch.rdb.HDel(ctx, countKey, instanceID).Err()
		}
	}
}
