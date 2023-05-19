package electrum

import (
	"context"
	"crypto/tls"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type WebSocketTransport struct {
	conn      *websocket.Conn
	responses chan []byte
	errors    chan error
}

// NewWebSocketTransport initializes new WebSocket transport.
func NewWebSocketTransport(
	ctx context.Context,
	url string,
	tlsConfig *tls.Config,
) (*WebSocketTransport, error) {
	dialer := websocket.Dialer{
		TLSClientConfig: tlsConfig,
	}

	conn, response, err := dialer.DialContext(ctx, url, nil)
	if err != nil {
		if DebugMode {
			log.Printf(
				"%s [debug] connect -> status: %v, error: %v",
				time.Now().Format("2006-01-02 15:04:05"),
				response.Status,
				err,
			)
		}
		return nil, err
	}

	ws := &WebSocketTransport{
		conn:      conn,
		responses: make(chan []byte),
		errors:    make(chan error),
	}

	go ws.listen()

	return ws, nil
}

func (t *WebSocketTransport) listen() {
	defer t.conn.Close()

	for {
		_, msg, err := t.conn.ReadMessage()
		if DebugMode {
			log.Printf(
				"%s [debug] %s -> msg: %s, err: %v",
				time.Now().Format("2006-01-02 15:04:05"),
				t.conn.RemoteAddr(),
				msg,
				err,
			)
		}
		if err != nil {
			t.errors <- err
			break
		}

		t.responses <- msg
	}
}

// SendMessage sends a message to the remote server through the WebSocket transport.
func (t *WebSocketTransport) SendMessage(body []byte) error {
	if DebugMode {
		log.Printf("%s [debug] %s <- %s", time.Now().Format("2006-01-02 15:04:05"), t.conn.RemoteAddr(), body)
	}

	return t.conn.WriteMessage(websocket.TextMessage, body)
}

// Responses returns chan to WebSocket transport responses.
func (t *WebSocketTransport) Responses() <-chan []byte {
	return t.responses
}

// Errors returns chan to WebSocket transport errors.
func (t *WebSocketTransport) Errors() <-chan error {
	return t.errors
}

// Close closes WebSocket transport.
func (t *WebSocketTransport) Close() error {
	// Cleanly close the connection by sending a close message and then
	// waiting (with timeout) for the server to close the connection.
	err := t.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Printf("%s [error] %s -> close error: %s", time.Now().Format("2006-01-02 15:04:05"), t.conn.RemoteAddr(), err)
	}

	return t.conn.Close()
}
