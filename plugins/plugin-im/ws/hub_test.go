package ws

import (
	"encoding/json"
	"testing"

	"hei-gin/sdk/enums"
)

func TestHubOnlineBroadcastStartStopIdempotent(t *testing.T) {
	h := NewHub()
	h.StartOnlineBroadcast()
	h.StartOnlineBroadcast()
	h.StopOnlineBroadcast()
	h.StopOnlineBroadcast()
}

func TestClientSendBytesToleratesClosedChannel(t *testing.T) {
	c := &Client{Send: make(chan []byte)}
	close(c.Send)

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("sendBytes panicked on closed channel: %v", r)
		}
	}()
	c.sendBytes([]byte("payload"))
}

func TestHubRegisterUsesPerUserCount(t *testing.T) {
	h := NewHub()
	clients := make([]*Client, 0, maxClientsPerUser)
	for i := 0; i < maxClientsPerUser; i++ {
		c := &Client{Hub: h, Send: make(chan []byte), UserID: "u1", UserType: "BUSINESS", IP: "127.0.0.1"}
		if !h.Register(c) {
			t.Fatalf("Register rejected client %d", i)
		}
		clients = append(clients, c)
	}

	extra := &Client{Hub: h, Send: make(chan []byte), UserID: "u1", UserType: "BUSINESS", IP: "127.0.0.1"}
	if h.Register(extra) {
		t.Fatal("Register accepted a client over per-user limit")
	}

	h.Unregister(clients[0])
	if !h.Register(extra) {
		t.Fatal("Register should accept after unregister decrements user count")
	}
}

func TestHubSendMessagesToUsersRoutesDistinctMessages(t *testing.T) {
	h := NewHub()
	c1 := &Client{Hub: h, Send: make(chan []byte, 1), UserID: "u1", UserType: enums.LoginTypeBusiness}
	c2 := &Client{Hub: h, Send: make(chan []byte, 1), UserID: "u2", UserType: enums.LoginTypeBusiness}
	c3 := &Client{Hub: h, Send: make(chan []byte, 1), UserID: "u3", UserType: enums.LoginTypeBusiness}
	if !h.Register(c1) || !h.Register(c2) || !h.Register(c3) {
		t.Fatal("failed to register test clients")
	}

	h.SendMessagesToUsers(map[string]Message{
		"u1": {Type: "m1"},
		"u2": {Type: "m2"},
	})

	assertNextMessageType(t, c1.Send, "m1")
	assertNextMessageType(t, c2.Send, "m2")
	select {
	case data := <-c3.Send:
		t.Fatalf("unexpected message for u3: %s", data)
	default:
	}
}

func assertNextMessageType(t *testing.T, ch <-chan []byte, want MessageType) {
	t.Helper()
	select {
	case data := <-ch:
		var msg Message
		if err := json.Unmarshal(data, &msg); err != nil {
			t.Fatalf("invalid message json: %v", err)
		}
		if msg.Type != want {
			t.Fatalf("message type = %q, want %q", msg.Type, want)
		}
	default:
		t.Fatalf("missing message %q", want)
	}
}
