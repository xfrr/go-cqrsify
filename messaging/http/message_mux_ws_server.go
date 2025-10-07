package messaginghttp

import (
	"bytes"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	defaultReadBufferSize  = 1024
	defaultWriteBufferSize = 1024
	// defaultWriteWait is the time allowed to write a message to the peer.
	defaultWriteWait = 10 * time.Second
	// pongWait is the time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
	// pingPeriod is the period at which pings are sent to the peer. It must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
	// maxMessageSize is the maximum message size allowed from peer.
	maxMessageSize = 5120 // 5KB
)

// MessageMUXWebsocketServer allows receiving Messages over WebSocket connections
// and dispatching them to a MessagePublisher (e.g., CommandBus or EventBus).
type MessageMUXWebsocketServer struct {
	handler      *MessageHandler
	upgrader     websocket.Upgrader
	errorHandler func(error)
}

// NewMessageWebsocketServer creates a new MessageWebsocketServer
// with the given MessagePublisher and options.
func NewMessageWebsocketServer(handler *MessageHandler) *MessageMUXWebsocketServer {
	return &MessageMUXWebsocketServer{
		handler: handler,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  defaultReadBufferSize,
			WriteBufferSize: defaultWriteBufferSize,
			CheckOrigin: func(_ *http.Request) bool {
				return true // Allow all origins, customize as needed
			},
		},
	}
}

// ServeHTTP implements http.Handler.
func (s *MessageMUXWebsocketServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Handle write and read in separate goroutines to free up the main goroutine
	go s.readPump(conn, w, r)
	go s.writePump(conn)
}

func (s *MessageMUXWebsocketServer) Close() error {
	// No persistent resources to close in this implementation
	return nil
}

func (s *MessageMUXWebsocketServer) readPump(conn *websocket.Conn, w http.ResponseWriter, r *http.Request) {
	defer func() {
		closeErr := conn.Close()
		if closeErr != nil {
			s.errorHandler(closeErr)
		}
	}()

	conn.SetReadLimit(maxMessageSize)
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error { _ = conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				s.errorHandler(err)
			}
			break
		}

		req, reqErr := http.NewRequestWithContext(r.Context(), http.MethodPost, "/", bytes.NewReader(msg))
		if reqErr != nil {
			s.errorHandler(reqErr)
			continue
		}

		req.Header.Set("Content-Type", "application/vnd.api+json")
		s.handler.ServeHTTP(w, req)
	}
}

func (s *MessageMUXWebsocketServer) writePump(conn *websocket.Conn) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		closeErr := conn.Close()
		if closeErr != nil {
			s.errorHandler(closeErr)
		}
	}()

	for range ticker.C {
		_ = conn.SetWriteDeadline(time.Now().Add(defaultWriteWait))
		if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
			return
		}
	}
}
