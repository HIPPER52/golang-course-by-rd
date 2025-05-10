package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"lesson_13/internal/commands"
	"lesson_13/internal/document_store"
	"log"
	"net"
)

var store = documentstore.NewStore()

func handleConnection(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	writer := bufio.NewWriter(conn)

	for scanner.Scan() {
		line := scanner.Bytes()

		var inMsg commands.CommandMessage
		if err := json.Unmarshal(line, &inMsg); err != nil {
			sendError(writer, "invalid command format")
			continue
		}

		outMsg, err := processCommand(inMsg)
		if err != nil {
			outMsg = commands.CommandMessage{
				Command: inMsg.Command + "_response",
				Payload: mustMarshal(map[string]string{"error": err.Error()}),
			}
		}

		sendMessage(writer, outMsg)
	}
	if err := scanner.Err(); err != nil {
		log.Printf("connection error: %v", err)
	}
}

func processCommand(in commands.CommandMessage) (commands.CommandMessage, error) {
	cmd := in.Command + "_response"
	var resp commands.CommandMessage
	resp.Command = cmd

	switch in.Command {
	case commands.CreateCollectionCommand:
		var req commands.CreateCollectionRequest
		if err := json.Unmarshal(in.Payload, &req); err != nil {
			return resp, err
		}

		_, err := store.CreateCollection(req.Name, &documentstore.CollectionConfig{PrimaryKey: req.Name})
		if err != nil {
			return resp, err
		}
		resp.Payload = mustMarshal(commands.CreateCollectionResponse{})

	case commands.DeleteCollectionCommand:
		var req commands.DeleteCollectionRequest
		if err := json.Unmarshal(in.Payload, &req); err != nil {
			return resp, err
		}
		if err := store.DeleteCollection(req.Name); err != nil {
			return resp, err
		}
		resp.Payload = mustMarshal(commands.DeleteCollectionResponse{})

	case commands.ListCollectionsCommand:
		names := make([]string, 0, len(store.Collections))
		for name := range store.Collections {
			names = append(names, name)
		}
		resp.Payload = mustMarshal(commands.ListCollectionsResponse{Collections: names})

	case commands.PutDocumentCommand:
		var req commands.PutDocumentRequest
		if err := json.Unmarshal(in.Payload, &req); err != nil {
			return resp, err
		}
		coll, err := store.GetCollection(req.Collection)
		if err != nil {
			return resp, err
		}

		doc, err := documentstore.MarshalDocument(req.Document)
		if err != nil {
			return resp, err
		}
		if err := coll.Put(*doc); err != nil {
			return resp, err
		}
		resp.Payload = mustMarshal(commands.PutDocumentResponse{})

	case commands.GetDocumentCommand:
		var req commands.GetDocumentRequest
		if err := json.Unmarshal(in.Payload, &req); err != nil {
			return resp, err
		}
		coll, err := store.GetCollection(req.Collection)
		if err != nil {
			return resp, err
		}
		doc, err := coll.Get(req.Key)
		if err != nil {
			r := commands.GetDocumentResponse{Found: false}
			resp.Payload = mustMarshal(r)
			return resp, nil
		}

		m := make(map[string]interface{})
		if err := documentstore.UnmarshalDocument(doc, &m); err != nil {
			return resp, err
		}
		r := commands.GetDocumentResponse{
			Document: m,
			Found:    true,
		}
		resp.Payload = mustMarshal(r)

	case commands.DeleteDocumentCommand:
		var req commands.DeleteDocumentRequest
		if err := json.Unmarshal(in.Payload, &req); err != nil {
			return resp, err
		}
		coll, err := store.GetCollection(req.Collection)
		if err != nil {
			return resp, err
		}
		success := coll.Delete(req.Key) == nil
		r := commands.DeleteDocumentResponse{Success: success}
		resp.Payload = mustMarshal(r)

	case commands.ListDocumentsCommand:
		var req commands.ListDocumentsRequest
		if err := json.Unmarshal(in.Payload, &req); err != nil {
			return resp, err
		}
		coll, err := store.GetCollection(req.Collection)
		if err != nil {
			return resp, err
		}
		docs := coll.List()

		outDocs := make([]map[string]interface{}, 0, len(docs))
		for _, d := range docs {
			m := make(map[string]interface{})
			if err := documentstore.UnmarshalDocument(&d, &m); err != nil {
				return resp, err
			}
			outDocs = append(outDocs, m)
		}
		r := commands.ListDocumentsResponse{Documents: outDocs}
		resp.Payload = mustMarshal(r)

	default:
		resp.Payload, _ = json.Marshal(map[string]string{"error": "unknown command"})
	}

	return resp, nil
}

func sendMessage(w *bufio.Writer, msg commands.CommandMessage) {
	data, _ := json.Marshal(msg)
	w.Write(data)
	w.WriteByte('\n')
	w.Flush()
}

func sendError(w *bufio.Writer, text string) {
	fmt.Fprintf(w, `{"error":%q}`+"\n", text)
	w.Flush()
}

func mustMarshal(v any) json.RawMessage {
	b, _ := json.Marshal(v)
	return b
}

func main() {
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	fmt.Println("Server listening on :8080")

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			continue
		}
		go handleConnection(conn)
	}
}
