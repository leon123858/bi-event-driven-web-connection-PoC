package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/websocket"
)

func Echo(ws *websocket.Conn) {
	var err error

	for {
		var reply string

		if err = websocket.Message.Receive(ws, &reply); err != nil {
			fmt.Println("Can't receive")
			break
		}

		fmt.Println("Received back from client: " + reply)

		msg := "Received:  " + reply
		fmt.Println("Sending to client: " + msg)

		if err = websocket.Message.Send(ws, msg); err != nil {
			fmt.Println("Can't send")
			break
		}
	}
}

func main() {
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "1234"
	}
	http.Handle("/", websocket.Handler(Echo))

	println("Server started: on port " + PORT)
	if err := http.ListenAndServe(":"+PORT, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
