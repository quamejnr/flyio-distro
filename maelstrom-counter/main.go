package main

import (
	"context"
	"encoding/json"
	"log"
	// "sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {	
	
	// mu := &sync.Mutex{}

	n := maelstrom.NewNode()
	kv := maelstrom.NewSeqKV(n)
	ctx := context.Background()	

	// kv.CompareAndSwap(ctx, "globalCounter", 0, globalCounter, true)

	// Handler for add
	n.Handle("add", func(msg maelstrom.Message) error {
		// Unmarshal the message body as an loosely-typed map.
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}
	
		// Update globalCounter
		delta := int(body["delta"].(float64))
		// mu.Lock()
		key := n.ID()
		counter, _ := kv.ReadInt(ctx, key)
		counter += delta
		kv.Write(ctx, key, counter)
		// mu.Unlock()
		
		// Update the message type to return back.
		body["type"] = "add_ok"
		delete(body, "delta")
	
		// Echo the original message back with the updated message type.
		return n.Reply(msg, body)
	})

	// Handler for read 
	n.Handle("read", func(msg maelstrom.Message) error {
		// Unmarshal the message body as an loosely-typed map.
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}
	
		// Update the message type to return back.
		// mu.Lock()
		value, _ := kv.ReadInt(ctx, n.ID())
		// if err != nil {
		// 	return err
		// }
		// mu.Unlock()
		
		body["type"] = "read_ok"
		body["value"] = value
	
		// Echo the original message back with the updated message type.
		return n.Reply(msg, body)
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}	
}