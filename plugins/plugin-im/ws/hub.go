package ws

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"hei-gin/sdk/enums"
	"hei-gin/sdk/observability"

	"github.com/gorilla/websocket"
)

var (
	_once        sync.Once
	_cfg         WSConfig
	_upgrader    *websocket.Upgrader
	_initHubOnce sync.Once
)

func getConfig() WSConfig {
	_once.Do(func() {
		_cfg = loadConfig()
		_upgrader = &websocket.Upgrader{
			ReadBufferSize:  _cfg.ReadBufferSize,
			WriteBufferSize: _cfg.WriteBufferSize,
			CheckOrigin:     checkOrigin,
		}
	})
	return _cfg
}

func getUpgrader() *websocket.Upgrader {
	getConfig()
	return _upgrader
}

// maxClientsPerIP limits the number of concurrent WebSocket connections from a single IP.
const maxClientsPerIP = 10

// maxClientsPerUser limits the number of concurrent WebSocket connections for a single user.
const maxClientsPerUser = 3

// Hub maintains the set of active clients and broadcasts online counts.
type Hub struct {
	mu          sync.RWMutex
	clients     map[*Client]bool
	lifecycleMu sync.Mutex
	onlineStop  chan struct{}
	onlineDone  chan struct{}

	// ipCount tracks connections per IP for rate limiting
	ipCount map[string]int

	// userCount tracks connections per login-type/user pair for O(1) per-user limiting.
	userCount map[string]int

	// Lifecycle hooks for CrossHub integration.
	OnClientRegistered   func(c *Client)
	OnClientUnregistered func(c *Client)
}

// NewHub creates a new Hub.
func NewHub() *Hub {
	return &Hub{
		clients:   make(map[*Client]bool),
		ipCount:   make(map[string]int),
		userCount: make(map[string]int),
	}
}

// Register adds a client to the hub with IP and per-user rate limiting.
func (h *Hub) Register(client *Client) bool {
	h.mu.Lock()

	// Per-IP connection limit
	ip := client.IP
	if ip != "" {
		if h.ipCount[ip] >= maxClientsPerIP {
			h.mu.Unlock()
			observability.IncWSRejected()
			log.Printf("[WS] IP %s exceeded max connections (%d)", ip, maxClientsPerIP)
			return false
		}
	}

	userKey := client.userKey()
	if h.userCount[userKey] >= maxClientsPerUser {
		h.mu.Unlock()
		observability.IncWSRejected()
		log.Printf("[WS] User %s/%s exceeded max connections (%d)",
			client.UserType, client.UserID, maxClientsPerUser)
		return false
	}

	h.clients[client] = true
	if ip != "" {
		h.ipCount[ip]++
	}
	h.userCount[userKey]++
	count := len(h.clients)
	onRegistered := h.OnClientRegistered
	h.mu.Unlock()

	if onRegistered != nil {
		onRegistered(client)
	}
	observability.IncWSConnection()

	log.Printf("[WS] Client connected: %s/%s from %s (online: %d)", client.UserType, client.UserID, ip, count)
	return true
}

// Unregister removes a client from the hub.
func (h *Hub) Unregister(client *Client) {
	h.mu.Lock()
	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		close(client.Send)

		ip := client.IP
		if ip != "" {
			h.ipCount[ip]--
			if h.ipCount[ip] <= 0 {
				delete(h.ipCount, ip)
			}
		}
		userKey := client.userKey()
		h.userCount[userKey]--
		if h.userCount[userKey] <= 0 {
			delete(h.userCount, userKey)
		}
	}
	count := len(h.clients)
	h.mu.Unlock()

	if h.OnClientUnregistered != nil {
		h.OnClientUnregistered(client)
	}
	observability.DecWSConnection()

	log.Printf("[WS] Client disconnected: %s/%s (online: %d)", client.UserType, client.UserID, count)
}

// OnlineCount returns the number of currently connected clients.
func (h *Hub) OnlineCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// isUserConnected checks if a specific user is connected to this hub.
func (h *Hub) isUserConnected(userID string, userType enums.LoginTypeEnum) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for client := range h.clients {
		if client.UserID == userID && client.UserType == userType {
			return true
		}
	}
	return false
}

// SendToUser sends a message to a specific business (admin) user.

// SendToUsers sends a message to multiple business users in a single lock acquisition.
func (h *Hub) SendToUsers(userIDs []string, msg Message) {
	userSet := make(map[string]struct{}, len(userIDs))
	for _, uid := range userIDs {
		userSet[uid] = struct{}{}
	}
	clients := make([]*Client, 0, len(userIDs))
	h.mu.RLock()
	for client := range h.clients {
		if client.UserType == enums.LoginTypeBusiness {
			if _, ok := userSet[client.UserID]; ok {
				clients = append(clients, client)
			}
		}
	}
	h.mu.RUnlock()
	for _, client := range clients {
		client.SendJSON(msg)
	}
	observability.ObserveWSMessage("business_users", len(clients))
}

// SendToConsumers sends a message to multiple consumer users in a single lock acquisition.
func (h *Hub) SendToConsumers(userIDs []string, msg Message) {
	userSet := make(map[string]struct{}, len(userIDs))
	for _, uid := range userIDs {
		userSet[uid] = struct{}{}
	}
	clients := make([]*Client, 0, len(userIDs))
	h.mu.RLock()
	for client := range h.clients {
		if client.UserType == enums.LoginTypeConsumer {
			if _, ok := userSet[client.UserID]; ok {
				clients = append(clients, client)
			}
		}
	}
	h.mu.RUnlock()
	for _, client := range clients {
		client.SendJSON(msg)
	}
	observability.ObserveWSMessage("consumer_users", len(clients))
}

func (h *Hub) SendMessagesToUsers(messages map[string]Message) {
	h.sendMessagesByUser(enums.LoginTypeBusiness, messages)
}

func (h *Hub) SendMessagesToConsumers(messages map[string]Message) {
	h.sendMessagesByUser(enums.LoginTypeConsumer, messages)
}

func (h *Hub) sendMessagesByUser(userType enums.LoginTypeEnum, messages map[string]Message) {
	if len(messages) == 0 {
		return
	}
	clients := make([]*Client, 0, len(messages))
	clientMsgs := make([]Message, 0, len(messages))
	h.mu.RLock()
	for client := range h.clients {
		if client.UserType != userType {
			continue
		}
		msg, ok := messages[client.UserID]
		if !ok {
			continue
		}
		clients = append(clients, client)
		clientMsgs = append(clientMsgs, msg)
	}
	h.mu.RUnlock()
	for i, client := range clients {
		client.SendJSON(clientMsgs[i])
	}
	observability.ObserveWSMessage("by_user", len(clients))
}

func (h *Hub) SendToUser(userID string, msg Message) {
	clients := make([]*Client, 0, 1)
	h.mu.RLock()
	for client := range h.clients {
		if client.UserType == enums.LoginTypeBusiness && client.UserID == userID {
			clients = append(clients, client)
		}
	}
	h.mu.RUnlock()
	for _, client := range clients {
		client.SendJSON(msg)
	}
	observability.ObserveWSMessage("business_single", len(clients))
}

// SendToConsumer sends a message to a specific consumer (client) user.
func (h *Hub) SendToConsumer(userID string, msg Message) {
	clients := make([]*Client, 0, 1)
	h.mu.RLock()
	for client := range h.clients {
		if client.UserType == enums.LoginTypeConsumer && client.UserID == userID {
			clients = append(clients, client)
		}
	}
	h.mu.RUnlock()
	for _, client := range clients {
		client.SendJSON(msg)
	}
	observability.ObserveWSMessage("consumer_single", len(clients))
}

// BroadcastAll sends a message to all connected clients.
func (h *Hub) BroadcastAll(msg Message) {
	data, _ := json.Marshal(msg)
	clients := h.snapshotClients(func(*Client) bool { return true })
	for _, client := range clients {
		client.sendBytes(data)
	}
}

func (h *Hub) snapshotClients(match func(*Client) bool) []*Client {
	h.mu.RLock()
	defer h.mu.RUnlock()
	clients := make([]*Client, 0, len(h.clients))
	for client := range h.clients {
		if match(client) {
			clients = append(clients, client)
		}
	}
	return clients
}

// BroadcastBusiness sends a message to all connected business (admin) clients.
func (h *Hub) BroadcastBusiness(msg Message) {
	data, _ := json.Marshal(msg)
	clients := h.snapshotClients(func(client *Client) bool {
		return client.UserType == enums.LoginTypeBusiness
	})
	for _, client := range clients {
		client.sendBytes(data)
	}
}

func (h *Hub) BroadcastConsumers(msg Message) {
	data, _ := json.Marshal(msg)
	clients := h.snapshotClients(func(client *Client) bool {
		return client.UserType == enums.LoginTypeConsumer
	})
	for _, client := range clients {
		client.sendBytes(data)
	}
}

// StartOnlineBroadcast periodically broadcasts the online count to all clients.
func (h *Hub) StartOnlineBroadcast() {
	h.lifecycleMu.Lock()
	if h.onlineStop != nil {
		h.lifecycleMu.Unlock()
		return
	}
	stop := make(chan struct{})
	done := make(chan struct{})
	h.onlineStop = stop
	h.onlineDone = done
	h.lifecycleMu.Unlock()

	interval := time.Duration(getConfig().OnlineBroadcastInterval) * time.Second
	if interval <= 0 {
		interval = 60 * time.Second
	}
	ticker := time.NewTicker(interval)
	go func() {
		defer close(done)
		defer ticker.Stop()
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[WS] Online broadcast panicked: %v", r)
			}
		}()
		for {
			select {
			case <-stop:
				return
			case <-ticker.C:
				count := h.OnlineCount()
				h.BroadcastAll(Message{
					Type: MsgOnlineCount,
					Payload: OnlineCountPayload{
						Count: count,
					},
				})
			}
		}
	}()
}

func (h *Hub) StopOnlineBroadcast() {
	h.lifecycleMu.Lock()
	stop := h.onlineStop
	done := h.onlineDone
	if stop == nil {
		h.lifecycleMu.Unlock()
		return
	}
	h.onlineStop = nil
	h.onlineDone = nil
	close(stop)
	h.lifecycleMu.Unlock()
	<-done
}

// HandleWebSocket upgrades an HTTP connection to WebSocket and registers the client.
func (h *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request, userID string, userType enums.LoginTypeEnum) {
	conn, err := getUpgrader().Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[WS] Upgrade error: %v", err)
		return
	}

	client := &Client{
		Hub:      h,
		Conn:     conn,
		Send:     make(chan []byte, 256),
		UserID:   userID,
		UserType: userType,
		IP:       getClientIP(r),
	}

	if !h.Register(client) {
		conn.Close()
		return
	}

	go client.WritePump()
	go client.ReadPump()
}

// getClientIP extracts the client IP from the request, trusting proxy headers only for trusted proxies.
func getClientIP(r *http.Request) string {
	if isTrustedProxy(remoteIP(r)) {
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			if idx := strings.Index(xff, ","); idx >= 0 {
				return strings.TrimSpace(xff[:idx])
			}
			return strings.TrimSpace(xff)
		}
		if xri := r.Header.Get("X-Real-IP"); xri != "" {
			return strings.TrimSpace(xri)
		}
	}
	return remoteIP(r)
}

func remoteIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err == nil {
		return host
	}
	addr := strings.TrimSpace(r.RemoteAddr)
	if parsed := net.ParseIP(addr); parsed != nil {
		return parsed.String()
	}
	if strings.HasPrefix(addr, "[") && strings.HasSuffix(addr, "]") {
		return strings.Trim(addr, "[]")
	}
	return addr
}

func isTrustedProxy(ip string) bool {
	if ip == "" {
		return false
	}
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}
	for _, rule := range getConfig().TrustedProxies {
		rule = strings.TrimSpace(rule)
		if rule == "" {
			continue
		}
		if strings.Contains(rule, "/") {
			_, cidr, err := net.ParseCIDR(rule)
			if err == nil && cidr.Contains(parsedIP) {
				return true
			}
			continue
		}
		if proxyIP := net.ParseIP(rule); proxyIP != nil && proxyIP.Equal(parsedIP) {
			return true
		}
	}
	return false
}

func checkOrigin(r *http.Request) bool {
	origin := strings.TrimSpace(r.Header.Get("Origin"))
	if origin == "" {
		return true
	}
	parsedOrigin, err := url.Parse(origin)
	if err != nil || parsedOrigin.Host == "" {
		return false
	}
	for _, candidate := range getConfig().AllowedOrigins {
		candidate = strings.TrimSpace(candidate)
		if candidate == "" {
			continue
		}
		if candidate == "*" || candidate == origin || candidate == parsedOrigin.Host {
			return true
		}
		if strings.HasPrefix(candidate, "http://") || strings.HasPrefix(candidate, "https://") {
			parsedCandidate, err := url.Parse(candidate)
			if err == nil && strings.EqualFold(parsedCandidate.Scheme, parsedOrigin.Scheme) && strings.EqualFold(parsedCandidate.Host, parsedOrigin.Host) {
				return true
			}
		}
	}
	return false
}

// GlobalHub is the singleton hub instance used by the application.
// NOTE: NewHub() must NOT call getConfig() — it runs at package init time,
// before config.FindAndLoad(). Use lazy loading for any config-dependent fields.
var GlobalHub = NewHub()

var globalCrossHub *CrossHub

// InitCrossHub initializes the process-wide hub runtime exactly once.
func InitCrossHub(local *Hub, rdb *redis.Client) *CrossHub {
	_initHubOnce.Do(func() {
		globalCrossHub = NewCrossHub(local, rdb)
	})
	return globalCrossHub
}

func Runtime() *CrossHub {
	return globalCrossHub
}
