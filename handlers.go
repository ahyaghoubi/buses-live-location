package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	upgrader     = websocket.Upgrader{}
	clients      = make(map[*websocket.Conn]bool)
	clientsMutex = &sync.Mutex{}
	buses        = make(map[*websocket.Conn]bool)
	busesMutex   = &sync.Mutex{}
	broadcast    = make(chan []byte)
	everyBus     []Bus
)

type Bus struct {
	Latitude  string `json:"latitude" binding:"required"`
	Longitude string `json:"longitude" binding:"required"`
	BusId     int64  `json:"bus_id" binding:"required"`
}

func ClientsWSHandler(context *gin.Context) {
	conn, err := upgrader.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		context.String(http.StatusBadRequest, "Could not upgrade to websocket: %v", err)
		return
	}

	clientsMutex.Lock()
	clients[conn] = true
	clientsMutex.Unlock()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}

	clientsMutex.Lock()
	delete(clients, conn)
	clientsMutex.Unlock()
	conn.Close()
}

func HandleBroadcast() {
	for {
		msg := <-broadcast

		clientsMutex.Lock()
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				fmt.Printf("Broadcast error: %v\n", err)
				client.Close()
				delete(clients, client)
			}
		}
		clientsMutex.Unlock()
	}
}

func BusesWSHandler(context *gin.Context) {
	conn, err := upgrader.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		context.String(http.StatusBadRequest, "Could not upgrade to websocket: %v", err)
		return
	}
	busesMutex.Lock()
	buses[conn] = true
	busesMutex.Unlock()

	defer func() {
		busesMutex.Lock()
		delete(buses, conn)
		busesMutex.Unlock()
		conn.Close()
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Read error:", err)
			break
		}

		var bus Bus
		err = json.Unmarshal(msg, &bus)
		if err != nil {
			context.String(http.StatusBadRequest, "Could not upgrade to websocket: %v", err)
			return
		}

		for i, b := range everyBus {
			if b.BusId == bus.BusId {
				everyBus = slices.Delete(everyBus, i, i+1)
				break
			}
		}

		everyBus = append(everyBus, bus)

		busJSON, err := json.Marshal(everyBus)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{
				"message": "Bad Request!",
			})
			return
		}

		broadcast <- busJSON
	}
}

func BusesHandler(context *gin.Context) {
	var bus Bus
	err := context.ShouldBindJSON(&bus)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad Request!",
		})
		return
	}

	for i, b := range everyBus {
		if b.BusId == bus.BusId {
			everyBus = slices.Delete(everyBus, i, i+1)
			break
		}
	}

	everyBus = append(everyBus, bus)

	busJSON, err := json.Marshal(everyBus)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad Request!",
		})
		return
	}

	broadcast <- busJSON

	context.JSON(http.StatusOK, gin.H{
		"message": "Success",
	})
}

func BusGETHandler(context *gin.Context) {
	context.JSON(http.StatusOK, everyBus)
}
