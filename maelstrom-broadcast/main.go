package main

import (
	"encoding/json"
	"log"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type server struct {
	n *maelstrom.Node
	mu       sync.RWMutex
	messages []int
	
}

func (s *server) broadcastHandler(msg maelstrom.Message) error {
	// Unmarshal the message body as an loosely-typed map.
	var body map[string]any
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	// Store value in message body into messages
	if v, ok := body["message"]; ok {
		s.mu.Lock()
		s.messages = append(s.messages, int(v.(float64)))
		s.mu.Unlock()
	}
	// Update the message type to return back.

	// Echo the original message back with the updated message type.
	return s.n.Reply(msg, map[string]string{
		"type": "broadcast_ok",
	})
}

func (s *server) readHandler(msg maelstrom.Message) error {
	// Unmarshal the message body as an loosely-typed map.
	var body map[string]any
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}
	s.mu.RLock()
	messages := s.messages
	s.mu.RUnlock()

	// Echo the original message back with the updated message type.
	return s.n.Reply(msg, map[string]any{
		"type":     "read_ok",
		"messages": messages,
	})
}

func (s *server) topologyHandler(msg maelstrom.Message) error {
	// Unmarshal the message body as an loosely-typed map.
	var body map[string]any
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	// Echo the original message back with the updated message type.
	return s.n.Reply(msg, map[string]string{
		"type": "topology_ok",
	})

}

func main() {
	n := maelstrom.NewNode()
	s := &server{n: n}

	n.Handle("broadcast", s.broadcastHandler)
	n.Handle("read", s.readHandler)
	n.Handle("topology", s.topologyHandler)

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
