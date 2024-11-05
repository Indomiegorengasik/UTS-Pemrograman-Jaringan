package main

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Izinkan semua origin (ganti dengan kebijakan yang lebih ketat di produksi)
	},
}

var connections = make(map[*websocket.Conn]bool)
var connLock sync.Mutex

func main() {
	router := gin.Default()

	// Menyajikan HTML untuk client
	router.LoadHTMLFiles("client1.html", "client2.html")
	router.Static("/assets", "./assets")

	router.GET("/client1", func(c *gin.Context) {
		c.HTML(http.StatusOK, "client1.html", nil)
	})

	router.GET("/client2", func(c *gin.Context) {
		c.HTML(http.StatusOK, "client2.html", nil)
	})

	router.GET("/ws", handleWebSocket)
	router.POST("/donate", handleDonation)

	router.Run(":8080") // Jalankan server di port 8080
}

func handleWebSocket(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()

	connLock.Lock()
	connections[ws] = true
	connLock.Unlock()

	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			connLock.Lock()
			delete(connections, ws)
			connLock.Unlock()
			break
		}
	}
}

func handleDonation(c *gin.Context) {
	amountStr := c.PostForm("amount")
	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid donation amount"})
		return
	}

	message := "Donasi sebesar " + strconv.Itoa(amount) + " berhasil diterima!"
	broadcastMessage(message)

	c.JSON(http.StatusOK, gin.H{"message": message})
}

func broadcastMessage(message string) {
	connLock.Lock()
	defer connLock.Unlock()

	for conn := range connections {
		err := conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			conn.Close()
			delete(connections, conn)
		}
	}
}
