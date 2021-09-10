// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"seabattle/internal/chat"
	"seabattle/internal/game"
	"seabattle/internal/repo"
)

var addr = flag.String("addr", ":8080", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "./home.html")
}

func main() {
	flag.Parse()
	hub := chat.NewHub(game.New(repo.New()))

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// serveWs handles websocket requests from the peer.
func serveWs(hub *chat.Hub, w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-UserId")
	chatID := r.Header.Get("X-ChatId")

	if userID != "" || chatID != "" {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`want X-UserId and X-ChatId headers`))
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	room := hub.GetRoom(chatID)
	if room == nil {
		room = hub.CreateNew(chatID)
		go room.Run()
	}

	client := chat.NewClient(room, conn, userID)
	client.Run()
}
