package main

import (
	"log"
	"net/http"
	"video-chatapp/server"
)

func main() {
	server.AllRooms.Init()

	http.HandleFunc("/create", server.CreateRoomRequestHandler)
	http.HandleFunc("/join/", server.JoinRoomRequestHandler)

	log.Println("Starting Server on Post 8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
