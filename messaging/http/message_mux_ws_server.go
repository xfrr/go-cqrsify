package messaginghttp

import (
	"bytes"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/xfrr/go-cqrsify/messaging"
	"github.com/xfrr/go-cqrsify/pkg/apix"
	"github.com/xfrr/go-cqrsify/pkg/multierror"
)

const (
	// defaultReadBufferSize is the size of the read buffer used when upgrading the HTTP connection to a WebSocket connection.
	defaultReadBufferSize = 1024
	// defaultWriteBufferSize is the size of the write buffer used when upgrading the HTTP connection to a WebSocket connection.
	defaultWriteBufferSize = 1024
	// defaultWriteWait is the time allowed to write a message to the peer.
	defaultWriteWait = 10 * time.Second
	// defaultCloseDeadline is the time allowed to close the connection gracefully.
	defaultCloseDeadline = 15 * time.Second
	// defaultCloseReadWait is the time allowed to read the close message from the peer.
	defaultCloseReadWait = 5 * time.Second
	// defaultPongWait is the time allowed to read the next pong message from the peer.
	defaultPongWait = 60 * time.Second
	// defaultPingPeriod is the period at which pings are sent to the peer. It must be less than pongWait.
	defaultPingPeriod = (defaultPongWait * 9) / 10
	// defaultMaxMessageSize is the maximum message size allowed from peer.
	defaultMaxMessageSize = 5120 // 5KB
)

// MessageWebsocketServer allows receiving Messages over WebSocket connections
// and dispatching them to a MessagePublisher (e.g., CommandBus or EventBus).
type MessageWebsocketServer struct {
	handler      *MessageHandler
	upgrader     websocket.Upgrader
	errorHandler func(error)
	clients      map[string]client
}

type client struct {
	conn *websocket.Conn
	send chan messaging.Message
}

// NewMessageWebsocketServer creates a new MessageMUXWebsocketServer with the given
// MessagePublisher.
func NewMessageWebsocketServer(publisher messaging.MessagePublisher) *MessageWebsocketServer {
	handler := NewMessageHandler(publisher)
	return newMessageWebsocketServer(handler)
}

func newMessageWebsocketServer(handler *MessageHandler) *MessageWebsocketServer {
	wsServer := &MessageWebsocketServer{
		handler:      handler,
		errorHandler: func(_ error) {},
		clients:      make(map[string]client),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  defaultReadBufferSize,
			WriteBufferSize: defaultWriteBufferSize,
			CheckOrigin: func(_ *http.Request) bool {
				return true // Allow all origins, customize as needed
			},
		},
	}

	return wsServer
}

// ServeHTTP implements http.Handler.
func (s *MessageWebsocketServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.errorHandler(err)
		return
	}

	// Register the new client
	clientID := conn.RemoteAddr().String()
	s.clients[clientID] = client{
		conn: conn,
		// TODO: implement message sending to client
		send: make(chan messaging.Message),
	}

	// Handle write and read in separate goroutines to free up the main goroutine
	go s.readPump(conn, w, r)
	go s.writePump(conn)
}

func (s *MessageWebsocketServer) Close() error {
	multierr := multierror.New()

	for _, client := range s.clients {
		// Send close message with deadline to allow graceful shutdown
		deadline := time.Now().Add(defaultCloseDeadline)
		err := client.conn.WriteControl(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(
				websocket.CloseNormalClosure,
				"Server is shutting down",
			),
			deadline,
		)
		if err != nil {
			return err
		}

		// Set deadline for reading the next message
		err = client.conn.SetReadDeadline(time.Now().Add(defaultCloseReadWait))
		if err != nil {
			return err
		}

		// Read messages until the close message is confirmed
		for {
			_, _, err = client.conn.NextReader()
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				break
			}
			if err != nil {
				break
			}
		}

		// Definitely close the connection
		err = client.conn.Close()
		if err != nil {
			multierr.Append(err)
		}

		// Remove client from the map
		delete(s.clients, client.conn.RemoteAddr().String())
	}

	return multierr.ErrorOrNil()
}

func (s *MessageWebsocketServer) readPump(conn *websocket.Conn, w http.ResponseWriter, r *http.Request) {
	// Ensure connection is closed when this function exits
	defer func() {
		_ = conn.Close()
		delete(s.clients, conn.RemoteAddr().String())
	}()

	conn.SetReadLimit(defaultMaxMessageSize)
	err := conn.SetReadDeadline(time.Now().Add(defaultPongWait))
	if err != nil {
		s.errorHandler(err)
		return
	}

	conn.SetPongHandler(func(string) error { _ = conn.SetReadDeadline(time.Now().Add(defaultPongWait)); return nil })

	for {
		_, msg, readErr := conn.ReadMessage()
		if readErr != nil {
			if websocket.IsUnexpectedCloseError(
				readErr,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure) {
				s.errorHandler(readErr)
			}
			break
		}

		req, reqErr := http.NewRequestWithContext(r.Context(), http.MethodPost, "/", bytes.NewReader(msg))
		if reqErr != nil {
			s.errorHandler(reqErr)
			continue
		}

		req.Header.Set("Content-Type", apix.ContentTypeJSONAPI.String())
		s.handler.ServeHTTP(w, req)
	}
}

func (s *MessageWebsocketServer) writePump(conn *websocket.Conn) {
	ticker := time.NewTicker(defaultPingPeriod)
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
