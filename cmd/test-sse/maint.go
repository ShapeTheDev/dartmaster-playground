package main

import (
	"encoding/json"
	"fmt"

	uricaller "github.com/One-Hundred-Eighty/Circle/pkg/uri-caller"
)

// Auto-Struktur
type Auto struct {
	Marke   string `json:"marke"`
	Baujahr int    `json:"baujahr"`
}

func main() {
	dartcounterUriCaller := uricaller.NewDartcounterUriCaller()

	eventStream, err := dartcounterUriCaller.DartcounterSSE()
	if err != nil {
		panic(err)
	}

	for sseEvent := range eventStream {
		// encode json
		var auto Auto
		err := json.Unmarshal([]byte(sseEvent.Data), &auto)
		if err != nil {
			panic(err)
		}

		fmt.Println("---------------------------")
		fmt.Printf("id: %s\n", sseEvent.Id)
		fmt.Printf("type: %s\n", sseEvent.Event)
		fmt.Printf("data: %s", sseEvent.Data)

		fmt.Println()

		fmt.Printf("id: %s\n", sseEvent.Id)
		fmt.Printf("type: %s\n", sseEvent.Event)
		fmt.Printf("data:\tMarke: %s\n", auto.Marke)
		fmt.Printf("\tBaujahr: %d\n", auto.Baujahr)
	}
}
