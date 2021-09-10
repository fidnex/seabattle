// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"net/http"

	"seabattle/internal/chat"
	"seabattle/internal/game"
	"seabattle/internal/repo"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", ":38080", "http service address")

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

	//err := http.ListenAndServeTLS(*addr, "ca-cert.pem", "ca-key.pem", nil)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// serveWs handles websocket requests from the peer.
func serveWs(hub *chat.Hub, w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	userID := q.Get("userId")
	chatID := q.Get("chatId")

	if userID == "" || chatID == "" {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`want userId and chatId params`))
		return
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
