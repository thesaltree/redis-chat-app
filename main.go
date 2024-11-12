package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
)

var (
	client *redis.Client
	Users  = make(map[string]*websocket.Conn)
	mu     sync.Mutex // sync to avoid concurrent map access
	sub    *redis.PubSub
)

const (
	chatChannel = "chats"
	serverAddr  = ":8081"
	redisAddr   = "localhost:6379"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func initRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{Addr: redisAddr})
	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	}
	return rdb
}

func main() {
	client = initRedisClient()
	defer client.Close()

	// Start WebSocket broadcaster
	go startChatBroadcaster()

	// Start HTTP server
	http.HandleFunc("/chat/", handleChat)
	server := &http.Server{Addr: serverAddr, Handler: nil}

	go func() {
		log.Printf("Server started on %s\n", serverAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown handling
	shutdown(server)
}

func handleChat(w http.ResponseWriter, r *http.Request) {
	user := strings.TrimPrefix(r.URL.Path, "/chat/")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	mu.Lock()
	Users[user] = conn
	mu.Unlock()
	log.Printf("User %s joined chat", user)

	defer func() {
		mu.Lock()
		delete(Users, user)
		mu.Unlock()
		conn.Close()
		log.Printf("User %s disconnected", user)
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Connection closed by user %s", user)
			} else {
				log.Printf("Error reading message from user %s: %v", user, err)
			}
			break
		}
		publishMessage(user, message)
	}
}

func publishMessage(user string, message []byte) {
	err := client.Publish(context.Background(), chatChannel, fmt.Sprintf("%s:%s", user, string(message))).Err()
	if err != nil {
		log.Printf("Failed to publish message: %v", err)
	}
}

func startChatBroadcaster() {
	log.Println("Listening to Redis messages")
	sub = client.Subscribe(context.Background(), chatChannel)
	defer sub.Close()

	for msg := range sub.Channel() {
		broadcastMessage(msg.Payload)
	}
}

func broadcastMessage(payload string) {
	parts := strings.SplitN(payload, ":", 2)
	if len(parts) != 2 {
		log.Println("Invalid message format")
		return
	}
	fromUser, message := parts[0], parts[1]

	mu.Lock()
	defer mu.Unlock()
	for user, conn := range Users {
		if user != fromUser {
			if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
				log.Printf("Failed to send message to user %s: %v", user, err)
			}
		}
	}
}

func shutdown(server *http.Server) {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGTERM, syscall.SIGINT)
	<-exit

	log.Println("Shutdown initiated")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Close all user connections
	mu.Lock()
	for user, conn := range Users {
		conn.Close()
		log.Printf("Closed connection for user %s", user)
	}
	mu.Unlock()

	// Unsubscribe and shutdown server
	if err := sub.Unsubscribe(context.Background(), chatChannel); err != nil {
		log.Printf("Failed to unsubscribe: %v", err)
	}

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Failed to shutdown server gracefully: %v", err)
	} else {
		log.Println("Server shut down successfully")
	}
}
