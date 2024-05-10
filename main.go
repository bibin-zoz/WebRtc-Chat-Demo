package main

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	db            *gorm.DB
	upgrader      = websocket.Upgrader{}
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

type Room struct {
	User1       int64
	User2       int64
	Connections []*websocket.Conn
	Ch          chan *Message
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
	r.GET("/ws", handleWebSocket)

	// Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func handleWebSocket(c *gin.Context) {
	// Upgrade initial GET request to a WebSocket
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Error upgrading to WebSocket: %v", err)
		return
	}
	defer ws.Close()

	// Get user ID from query parameters
	userID, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		log.Printf("Invalid user ID: %v", err)
		return
	}

	// Get receiver ID from query parameters
	receiverID, err := strconv.ParseInt(c.Query("receiver_id"), 10, 64)
	if err != nil {
		log.Printf("Invalid receiver ID: %v", err)
		return
	}

	roomID := generateRoomID(userID, receiverID)

	// Create a new room if it doesn't exist
	connectionsMu.Lock()
	if _, ok := rooms[roomID]; !ok {
		rooms[roomID] = &Room{
			User1:       userID,
			User2:       receiverID,
			Connections: []*websocket.Conn{ws},
			Ch:          make(chan *Message),
		}
		go broadcastMessages(roomID)
	} else {
		rooms[roomID].Connections = append(rooms[roomID].Connections, ws)
	}
	connectionsMu.Unlock()

	// Fetch previous messages from the database
	prevMessages, err := getPreviousMessages(userID, receiverID)
	if err != nil {
		log.Printf("Error fetching previous messages: %v", err)
	} else {
		// Send previous messages to the client
		for _, msg := range prevMessages {
			if err := ws.WriteJSON(msg); err != nil {
				log.Printf("Error sending previous message to connection: %v", err)
			}
		}
	}

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
		msg.Time = time.Now()
		if err := saveMessage(&msg); err != nil {
			log.Printf("Error saving message: %v", err)
			continue
		}

		// Send message to other user in the room
		rooms[roomID].Ch <- &msg
	}
}

func getPreviousMessages(userID, receiverID int64) ([]*Message, error) {
	var messages []*Message
	if err := db.Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)", userID, receiverID, receiverID, userID).Order("time").Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func generateRoomID(user1, user2 int64) string {
	// Sort user IDs to ensure consistency
	if user1 > user2 {
		user1, user2 = user2, user1
	}
	return strconv.FormatInt(user1, 10) + "-" + strconv.FormatInt(user2, 10)
}

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

func saveMessage(msg *Message) error {
	return db.Create(msg).Error
}
