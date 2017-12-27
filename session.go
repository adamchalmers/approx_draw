package main

import (
	"time"

	. "github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	conn      *websocket.Conn
	broadcast chan *Event
}

func newClient(conn *websocket.Conn) *Client {
	return &Client{
		conn:      conn,
		broadcast: make(chan *Event),
	}
}

// The Dispatcher; we spawn one for every client
func (c *Client) tx() {
	defer func() {
		c.conn.Close()
	}()

	for {
		event, ok := <-c.broadcast

		c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

		if !ok {
			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		c.conn.WriteJSON(event)
	}

}

type SessionManager struct {
	sessions map[UUID]*Client
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[UUID]*Client),
	}
}

func (sm *SessionManager) register(sessionID UUID, conn *websocket.Conn) {
	sm.sessions[sessionID] = &Client{
		conn:      conn,
		broadcast: make(chan *Event),
	}

	go sm.sessions[sessionID].tx()
}
