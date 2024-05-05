package main

import (
	customwebsocket "chatapplication/websocket"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func serverWebSocket(pool *customwebsocket.Pool, w http.ResponseWriter, r *http.Request) {
	conn, err := customwebsocket.Upgrade(w, r)
	if err != nil {
		log.Println("web socket upgrade error", err)
		return
	}
	client := &customwebsocket.Client{
		Conn: conn,
		Pool: pool,
	}
	pool.Register <- client
	client.Read()
}

func setupRoutes() {
	pool := customwebsocket.NewPool()

	log.Println("serever connected")

	go pool.Start()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serverWebSocket(pool, w, r)
	})
}

func main() {
	setupRoutes()
	server := &http.Server{Addr: ":9000"}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}
	log.Println("Server gracefully stopped")
}
