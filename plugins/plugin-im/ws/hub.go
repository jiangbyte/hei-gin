package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"strings"
	"time"

	"hei-gin/sdk/enums"
	"hei-gin/sdk/middleware"

	"github.com/gorilla/websocket"
)

// cfg is the package-level ws configuration, loaded once from config.C.Raw["ws"].
var cfg = loadConfig()

// maxClientsPerIP limits the number of concurrent WebSocket connections from a single IP.
const maxClientsPerIP = 10

// maxClientsPerUser limits the number of concurrent WebSocket connections for a single user.
const maxClientsPerUser = 3

var upgrader = websocket.Upgrader{
	ReadBufferSize:  cfg.ReadBufferSize,
	WriteBufferSize: cfg.WriteBufferSize,
	CheckOrigin: func(r *http.Request) bool {
		// In production, validate against config.C.CORS.AllowOrigins
		// For development, allow all origins (consistent with CORS config)
		return true
	},
}



// Hub maintains the set of active clients and broadcasts online counts.
type Hub struct {
	mu      sync.RWMutex
	clients map[*Client]bool

	// ipCount tracks connections per IP for rate limiting
	ipCount map[string]int

	// Lifecycle hooks for CrossHub integration.
	OnClientRegistered   func(c *Client)
	OnClientUnregistered func(c *Client)
}

// NewHub creates a new Hub.
func NewHub() *Hub {
	return &Hub{
		clients: make(map[*Client]bool),
		ipCount: make(map[string]int),
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
			log.Printf("[WS] IP %s exceeded max connections (%d)", ip, maxClientsPerIP)
			return false
		}
	}

	// Per-user connection limit
	userCount := 0
	for existing := range h.clients {
		if existing.UserID == client.UserID && existing.UserType == client.UserType {
			userCount++
			if userCount >= maxClientsPerUser {
				h.mu.Unlock()
				log.Printf("[WS] User %s/%s exceeded max connections (%d)",
					client.UserType, client.UserID, maxClientsPerUser)
				return false
			}
		}
	}

	h.clients[client] = true
	if ip != "" {
		h.ipCount[ip]++
	}
	count := len(h.clients)
	onRegistered := h.OnClientRegistered
	h.mu.Unlock()

	if onRegistered != nil {
		onRegistered(client)
	}

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
	}
	count := len(h.clients)
	h.mu.Unlock()

	if h.OnClientUnregistered != nil {
		h.OnClientUnregistered(client)
	}

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
	h.mu.RLock()
	defer h.mu.RUnlock()
	userSet := make(map[string]struct{}, len(userIDs))
	for _, uid := range userIDs {
		userSet[uid] = struct{}{}
	}
	for client := range h.clients {
		if client.UserType == enums.LoginTypeBusiness {
			if _, ok := userSet[client.UserID]; ok {
				client.SendJSON(msg)
			}
		}
	}
}

// SendToConsumers sends a message to multiple consumer users in a single lock acquisition.
func (h *Hub) SendToConsumers(userIDs []string, msg Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	userSet := make(map[string]struct{}, len(userIDs))
	for _, uid := range userIDs {
		userSet[uid] = struct{}{}
	}
	for client := range h.clients {
		if client.UserType == enums.LoginTypeConsumer {
			if _, ok := userSet[client.UserID]; ok {
				client.SendJSON(msg)
			}
		}
	}
}

func (h *Hub) SendToUser(userID string, msg Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for client := range h.clients {
		if client.UserType == enums.LoginTypeBusiness && client.UserID == userID {
			client.SendJSON(msg)
		}
	}
}

// SendToConsumer sends a message to a specific consumer (client) user.
func (h *Hub) SendToConsumer(userID string, msg Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for client := range h.clients {
		if client.UserType == enums.LoginTypeConsumer && client.UserID == userID {
			client.SendJSON(msg)
		}
	}
}

// BroadcastAll sends a message to all connected clients.
func (h *Hub) BroadcastAll(msg Message) {
	data, _ := json.Marshal(msg)
	h.mu.RLock()
	defer h.mu.RUnlock()
	for client := range h.clients {
		select {
		case client.Send <- data:
		default:
		}
	}
}

// BroadcastBusiness sends a message to all connected business (admin) clients.
func (h *Hub) BroadcastBusiness(msg Message) {
	data, _ := json.Marshal(msg)
	h.mu.RLock()
	defer h.mu.RUnlock()
	for client := range h.clients {
		if client.UserType == enums.LoginTypeBusiness {
			select {
			case client.Send <- data:
			default:
			}
		}
	}
}

func (h *Hub) BroadcastConsumers(msg Message) {
	data, _ := json.Marshal(msg)
	h.mu.RLock()
	defer h.mu.RUnlock()
	for client := range h.clients {
		if client.UserType == enums.LoginTypeConsumer {
			select {
			case client.Send <- data:
			default:
			}
		}
	}
}

// StartOnlineBroadcast periodically broadcasts the online count to all clients.
func (h *Hub) StartOnlineBroadcast() {
	interval := time.Duration(cfg.OnlineBroadcastInterval) * time.Second
	if interval <= 0 {
		interval = 60 * time.Second
	}
	ticker := time.NewTicker(interval)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[WS] Online broadcast panicked: %v", r)
			}
		}()
		for range ticker.C {
			middleware.GoSafe(func() {
				count := h.OnlineCount()
				h.BroadcastAll(Message{
					Type: MsgOnlineCount,
					Payload: OnlineCountPayload{
						Count: count,
					},
				})
			})
		}
	}()
}

// HandleWebSocket upgrades an HTTP connection to WebSocket and registers the client.
func (h *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request, userID string, userType enums.LoginTypeEnum) {
	conn, err := upgrader.Upgrade(w, r, nil)
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

// getClientIP extracts the client IP from the request, respecting X-Forwarded-For.
func getClientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		if idx := strings.Index(xff, ","); idx >= 0 {
			return xff[:idx]
		}
		return xff
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	// Strip port from RemoteAddr
	addr := r.RemoteAddr
	for i := len(addr) - 1; i >= 0; i-- {
		if addr[i] == ':' {
			return addr[:i]
		}
	}
	return addr
}

// GlobalHub is the singleton hub instance used by the application.
var GlobalHub = NewHub()

// GlobalCrossHub is the cross-instance hub. Initialized by app.go after Redis setup.
var GlobalCrossHub *CrossHub
