package ws

import (
	"encoding/json"
	"log"
	"time"

	"hei-gin/sdk/enums"

	"github.com/gorilla/websocket"
)

// Client represents a single WebSocket connection.
type Client struct {
	Hub      *Hub
	Conn     *websocket.Conn
	Send     chan []byte
	UserID   string
	UserType enums.LoginTypeEnum
	IP       string
}

// ReadPump pumps messages from the WebSocket connection to the hub.
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister(c)
		c.Conn.Close()
	}()

	pongTimeout := time.Duration(getConfig().PongTimeout) * time.Second
	c.Conn.SetReadDeadline(time.Now().Add(pongTimeout))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongTimeout))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Printf("[WS] Read error: %v", err)
			}
			break
		}

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}

		switch msg.Type {
		case MsgHeartbeat:
			c.SendPong()
		}
	}
}

// WritePump pumps messages from the hub to the WebSocket connection.
func (c *Client) WritePump() {
	heartbeatInterval := time.Duration(getConfig().HeartbeatInterval) * time.Second
	if heartbeatInterval <= 0 {
		heartbeatInterval = 30 * time.Second
	}
	writeTimeout := time.Duration(getConfig().WriteTimeout) * time.Second

	ticker := time.NewTicker(heartbeatInterval)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeTimeout))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeTimeout))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// SendPong replies to a heartbeat with an online count.
func (c *Client) SendPong() {
	data, _ := json.Marshal(Message{
		Type: MsgOnlineCount,
		Payload: OnlineCountPayload{
			Count: c.Hub.OnlineCount(),
		},
	})
	select {
	case c.Send <- data:
	default:
	}
}

// SendJSON sends a JSON message to this client.
func (c *Client) SendJSON(msg Message) {
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}
	select {
	case c.Send <- data:
	default:
	}
}
