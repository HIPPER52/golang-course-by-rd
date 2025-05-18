package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"lesson_12/internal/commands"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Fprintln(os.Stderr, "cannot connect:", err)
		os.Exit(1)
	}
	defer conn.Close()

	in := bufio.NewScanner(os.Stdin)
	out := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	fmt.Println("Examples commands:")
	fmt.Println("  create_collection {\"name\":\"users\",\"primary_key\":\"id\"}")
	fmt.Println("  put_document      {\"collection\":\"users\",\"document\":{...}}")
	fmt.Println("  list_collections")
	fmt.Println("  get_document      {\"collection\":\"users\",\"key\":\"1\"}")
	fmt.Println()

	for in.Scan() {
		line := in.Text()

		parts := strings.SplitN(line, " ", 2)
		cmd := parts[0]

		var payload string
		if len(parts) > 1 {
			payload = parts[1]
		} else {
			payload = "{}"
		}

		msg := commands.CommandMessage{
			Command: cmd,
			Payload: json.RawMessage(payload),
		}
		data, _ := json.Marshal(msg)

		writer.Write(data)
		writer.WriteByte('\n')
		writer.Flush()

		respLine, err := out.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "read error:", err)
			break
		}
		fmt.Print("-> ", respLine)
	}

	if err := in.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "stdin error:", err)
	}
}
