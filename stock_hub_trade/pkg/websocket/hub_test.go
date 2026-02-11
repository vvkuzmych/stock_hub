package websocket

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

// Mock message logger for testing
type mockLogger struct {
	messages []string
}

func (m *mockLogger) SaveMessage(content string, clientID string) error {
	m.messages = append(m.messages, content)
	return nil
}

func TestWritePump_SendsSeparateFrames(t *testing.T) {
	// Create test server
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	hub := NewHub(&mockLogger{})
	go hub.Run()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatal(err)
		}
		defer conn.Close()

		client := NewClient(hub, conn, "test-client")
		hub.Register(client)
		go client.WritePump()

		// Send multiple messages quickly to test batching behavior
		msg1 := map[string]string{"type": "message1", "content": "first"}
		msg2 := map[string]string{"type": "message2", "content": "second"}
		msg3 := map[string]string{"type": "message3", "content": "third"}

		json1, _ := json.Marshal(msg1)
		json2, _ := json.Marshal(msg2)
		json3, _ := json.Marshal(msg3)

		client.Send(json1)
		client.Send(json2)
		client.Send(json3)

		// Give WritePump time to process
		time.Sleep(100 * time.Millisecond)
	}))
	defer server.Close()

	// Connect as client
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	// Read messages and verify each is valid JSON
	receivedMessages := 0
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))

	for receivedMessages < 3 {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		// Verify it's valid JSON
		var parsed map[string]interface{}
		if err := json.Unmarshal(message, &parsed); err != nil {
			t.Errorf("Received invalid JSON: %s, error: %v", string(message), err)
		}

		// Verify no newlines (which would indicate batching)
		if strings.Contains(string(message), "\n") {
			t.Errorf("Message contains newline (batched): %s", string(message))
		}

		receivedMessages++
	}

	if receivedMessages != 3 {
		t.Errorf("Expected 3 messages, got %d", receivedMessages)
	}
}

func TestWritePump_ClosesOnChannelClose(t *testing.T) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	hub := NewHub(&mockLogger{})
	go hub.Run()

	closed := make(chan bool)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatal(err)
		}

		client := NewClient(hub, conn, "test-client")
		hub.Register(client)
		
		go func() {
			client.WritePump()
			closed <- true
		}()

		// Close the channel to trigger WritePump exit
		close(client.send)
		
		// Wait for WritePump to close
		select {
		case <-closed:
			// Success
		case <-time.After(1 * time.Second):
			t.Error("WritePump did not close within timeout")
		}
	}))
	defer server.Close()

	// Connect as client
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	// Wait for close message
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, _, err = conn.ReadMessage()
	
	// Should receive close error
	if err == nil {
		t.Error("Expected close error, got nil")
	}
}

func TestClient_Send_NonBlocking(t *testing.T) {
	hub := NewHub(&mockLogger{})
	
	// Create a mock connection (we won't actually use it)
	upgrader := websocket.Upgrader{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil)
		defer conn.Close()
		time.Sleep(5 * time.Second) // Keep server alive
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client := NewClient(hub, conn, "test")

	// Fill the buffer
	for i := 0; i < 256; i++ {
		client.Send([]byte("message"))
	}

	// This should not block (channel is full, but Send handles it)
	done := make(chan bool)
	go func() {
		client.Send([]byte("overflow message"))
		done <- true
	}()

	select {
	case <-done:
		// Success - Send did not block
	case <-time.After(100 * time.Millisecond):
		t.Error("Send blocked when buffer was full")
	}
}
