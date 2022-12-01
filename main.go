package main

import (
	"fmt"
	"log"
	"tcp-server/server"
)

func main() {
	server := server.NewServer(":8000")

	go func() {
		for msg := range server.Msgch {
			fmt.Printf("received message from connection (%s): %s\n", msg.From, string(msg.Payload))
		}
	}()

	log.Fatal(server.Start())
}
