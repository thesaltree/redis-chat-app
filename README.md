# Redis Pub/Sub Chat Application

This is a simple real-time chat application built with Go, WebSockets, and Redis Pub/Sub. Users can join the chat by connecting to the WebSocket server and will receive messages from other users in real-time.

## Features

- **Real-Time Communication**: Messages are sent and received instantly using WebSockets.
- **Redis Pub/Sub**: Used to broadcast messages to all connected users efficiently.
- **Graceful Shutdown**: Ensures all connections are closed properly when the server shuts down.

## Prerequisites

- [Go](https://golang.org/dl/) (version 1.16 or later)
- [Redis](https://redis.io/download) server (running on `localhost:6379` by default)

## Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/redis-chat-app.git
cd redis-chat-app
```
### 2. Install Dependencies
This application uses the github.com/go-redis/redis/v8 and github.com/gorilla/websocket packages. You can install them with:

```bash
go get github.com/go-redis/redis/v8
go get github.com/gorilla/websocket
```

### 3. Configure Redis
Ensure your Redis server is running on `localhost:6379`. If it's hosted elsewhere, update the `redisAddr` constant in the code to reflect the correct address.

### 4. Run the Application
Start the server by running:

```bash
go run main.go
```

The server will start on http://localhost:8081.

### 5. Connect to the Chat
To join the chat, open multiple browser tabs (or WebSocket clients) and connect to:

```bash
ws://localhost:8081/chat/username
```

Replace username with a unique name for each user. Messages sent by one user will be received by all other connected users.

## Reference Screenshots

Initiate chat

<img width="940" alt="Initiate chat" src="https://github.com/user-attachments/assets/d1c87c4c-2ed3-4147-9770-60707484a3bf">


User 1

<img width="957" alt="user1" src="https://github.com/user-attachments/assets/175fdf68-6df5-4b1d-885f-ab4b1f170db8">


User 2

<img width="954" alt="user2" src="https://github.com/user-attachments/assets/e9aafc9e-0e89-4600-9bbf-650f82c9b5a9">


User 3

<img width="956" alt="user3" src="https://github.com/user-attachments/assets/9c7a05ea-6af5-46b0-b2c9-d1e783a26923">




