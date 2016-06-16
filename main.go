package main

// Test json decoder: see if it can handle json logs
import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func main() {
	const jsonStream = `
	{"Name": "Ed", "Text": "Knock knock."}
	{"Name": "Sam", "Text": "Who's there?"}
	{"Name": "Ed", "Text": "Go fmt."}
	{"Name": "Sam", "Text": "Go fmt who?"}
	{"Name": "Ed", "Text": "Go fmt yourself!"}
`
	type Message struct {
		Name, Text string
	}
	dec := json.NewDecoder(strings.NewReader(jsonStream))
	var m Message
	for dec.More() {
		// decode an array value (Message)
		err := dec.Decode(&m)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error decoding: %s", err)
			continue
		}
		fmt.Printf("%v: %v\n", m.Name, m.Text)
	}
}
