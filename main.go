package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	db       *gorm.DB
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all connections
		},
	}

	rooms         = make(map[string]*Room)
	connectionsMu sync.Mutex
)

type Message struct {
	SenderID   int64     `json:"sender_id"`
	ReceiverID int64     `json:"receiver_id"`
	Message    string    `json:"message"`
	Time       time.Time `json:"time"`
}

type User struct {
	UserID int64  `json:"user_id" gorm:"primary_key"`
	Name   string `json:"name"`
}

func main() {
	// Connect to PostgreSQL
	var err error
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		"localhost", "5432", "postgres", "8596", "demo_svc")
	db, err = gorm.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&User{}, &Message{})

	// Initialize Gin
	r := gin.Default()

	// WebSocket endpoint
	r.GET("/ws", HandleWebSocket)

	// Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

type Room struct {
	GroupID     int64
	Connections []*websocket.Conn
	Ch          chan *Message
}

// HandleWebSocket function to handle WebSocket connections
func HandleWebSocket(c *gin.Context) {
	// Upgrade initial GET request to a WebSocket
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Error upgrading to WebSocket: %v", err)
		return
	}
	defer ws.Close()

	// // Get user ID and group ID from query parameters
	// userID, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	// if err != nil {
	// 	log.Printf("Invalid user ID: %v", err)
	// 	return
	// }

	groupID, err := strconv.ParseInt(c.Query("group_id"), 10, 64)
	if err != nil {
		log.Printf("Invalid group ID: %v", err)
		return
	}

	roomID := strconv.FormatInt(groupID, 10) // Use group ID as room ID for simplicity

	// Create a new room if it doesn't exist
	connectionsMu.Lock()
	if _, ok := rooms[roomID]; !ok {
		rooms[roomID] = &Room{
			GroupID:     groupID,
			Connections: []*websocket.Conn{ws},
			Ch:          make(chan *Message),
		}
		go broadcastMessages(roomID)
	} else {
		rooms[roomID].Connections = append(rooms[roomID].Connections, ws)
	}
	connectionsMu.Unlock()

	// Fetch previous messages from the database (if needed)

	// Remove connection when this function returns
	defer func() {
		connectionsMu.Lock()
		for i, conn := range rooms[roomID].Connections {
			if conn == ws {
				rooms[roomID].Connections = append(rooms[roomID].Connections[:i], rooms[roomID].Connections[i+1:]...)
				break
			}
		}
		connectionsMu.Unlock()
	}()

	// Listen for incoming messages
	for {
		var msg Message
		if err := ws.ReadJSON(&msg); err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}
		// Save message to database
		// msg.Time = time.Now()
		// if err := saveMessage(&msg); err != nil {
		// 	log.Printf("Error saving message: %v", err)
		// 	continue
		// }

		// Broadcast message to all members in the room
		rooms[roomID].Ch <- &msg
	}
}

// broadcastMessages function to broadcast messages to all members in the room
func broadcastMessages(roomID string) {
	room := rooms[roomID]
	for msg := range room.Ch {
		for _, conn := range room.Connections {
			if err := conn.WriteJSON(msg); err != nil {
				log.Printf("Error sending message to connection: %v", err)
			}
		}
	}
}

// saveMessage function to save message to the database
// func saveMessage(msg *Message) error {
// 	return db.Create(msg).Error
// }
