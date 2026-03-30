package main

import (
	"net/http"

	"nhooyr.io/websocket"
)

var clients map[*Client]bool

func main() {
	mux := http.NewServeMux()
	hub := NewHub()
	go hub.Run()

	mux.HandleFunc("GET /connect", func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{InsecureSkipVerify: true})
		if err != nil {
			return
		}

		client := NewClient(hub, conn)
		hub.register <- client

		go client.WritePump(r.Context())
		client.ReadPump(r.Context())
	})
	mux.Handle("/", http.FileServer(http.Dir("html")))
	http.ListenAndServe(":8080", mux)
}
