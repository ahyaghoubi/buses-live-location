package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"sync"
	"time"

	"github.com/ahyaghoubi/buses-live-location/models"
	"github.com/ahyaghoubi/buses-live-location/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ActiveBus struct {
	Latitude   string    `json:"latitude" binding:"required"`
	Longitude  string    `json:"longitude" binding:"required"`
	BusId      int64     `json:"bus_id"`
	Name       string    `json:"name"`
	LastUpdate time.Time `json:"last_update"`
}

var allActiveBuses []ActiveBus

var (
	upgrader     = websocket.Upgrader{}
	clients      = make(map[*websocket.Conn]bool)
	clientsMutex = &sync.Mutex{}
	buses        = make(map[*websocket.Conn]bool)
	busesMutex   = &sync.Mutex{}
	broadcast    = make(chan []byte)
)

func getAllBuses(context *gin.Context) {
	buses, err := models.GetAllBuses()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not fetch the buses!",
		})
		return
	}
	context.JSON(http.StatusOK, buses)
}

func createBus(context *gin.Context) {
	var bus models.Bus

	err := context.ShouldBindJSON(&bus)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not parse data!",
		})
		return
	}

	err = bus.Create()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong!",
		})
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": "Bus created!",
		"bus":     bus,
	})
}

func updateBus(context *gin.Context) {
	busId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "The ID was not valid!",
		})
		return
	}
	_, err = models.GetBusById(busId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not fetch the bus",
			"error":   fmt.Sprintf("%v", err),
		})
		return
	}

	var updatedBus models.Bus
	err = context.ShouldBindJSON(&updatedBus)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad Request!",
		})
		return
	}

	updatedBus.ID = busId

	err = updatedBus.Update()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not Update bus!",
		})
		return
	}

	for i, activeBus := range allActiveBuses {
		if activeBus.BusId == busId {
			allActiveBuses[i].Name = updatedBus.Name
			break
		}
	}

	allActiveBusesJSON, err := json.Marshal(allActiveBuses)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad Request!!",
			"error":   fmt.Sprintf("%v", err),
		})
		return
	}

	broadcast <- allActiveBusesJSON

	context.JSON(http.StatusOK, gin.H{
		"messge": "Bus Updated",
	})
}

func deleteBus(context *gin.Context) {
	busId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "The ID is Invalid",
		})
		return
	}

	bus, err := models.GetBusById(busId)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"message": "Could not find the bus",
		})
		return
	}

	err = bus.Delete()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not delete the bus",
		})
		return
	}

	for i, activeBus := range allActiveBuses {
		if activeBus.BusId == busId {
			allActiveBuses = slices.Delete(allActiveBuses, i, i+1)
			break
		}
	}

	allActiveBusesJSON, err := json.Marshal(allActiveBuses)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad Request!!",
			"error":   fmt.Sprintf("%v", err),
		})
		return
	}

	broadcast <- allActiveBusesJSON

	context.JSON(http.StatusOK, gin.H{
		"message": "Bus Deleted",
	})
}

func busLogin(context *gin.Context) {
	busId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "The ID is Invalid",
		})
		return
	}

	bus, err := models.GetBusById(busId)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"message": "Could not find the bus",
		})
		return
	}

	token, err := utils.GenerateBusToken(bus.ID, bus.Name)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not authenticate!",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Generate token for bus was successful!",
		"token":   token,
	})
}

func UpdateActiveBusesEvery5min() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		filtered := allActiveBuses[:0]

		for _, bus := range allActiveBuses {
			if bus.LastUpdate.After(now.Add(-5 * time.Minute)) {
				filtered = append(filtered, bus)
			}
		}

		allActiveBuses = filtered

		allActiveBusesJSON, err := json.Marshal(allActiveBuses)
		if err != nil {
			fmt.Println(err)
			return
		}

		broadcast <- allActiveBusesJSON
	}
}

func getAllActiveBuses(context *gin.Context) {
	context.JSON(http.StatusOK, allActiveBuses)
}

func busStatusUpdate(context *gin.Context) {
	var activeBus ActiveBus
	err := context.ShouldBindJSON(&activeBus)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad request!",
			"error":   fmt.Sprintf("%v", err),
		})
		return
	}

	busId := context.GetInt64("busId")
	bus, err := models.GetBusById(busId)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not find this bus",
		})
		return
	}

	for i, bus := range allActiveBuses {
		if bus.BusId == busId {
			allActiveBuses = slices.Delete(allActiveBuses, i, i+1)
			break
		}
	}

	activeBus.Name = bus.Name
	activeBus.BusId = bus.ID

	activeBus.LastUpdate = time.Now()

	allActiveBuses = append(allActiveBuses, activeBus)

	allActiveBusesJSON, err := json.Marshal(allActiveBuses)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad Request!!",
			"error":   fmt.Sprintf("%v", err),
		})
		return
	}

	broadcast <- allActiveBusesJSON

	context.JSON(http.StatusOK, gin.H{
		"message": "Success",
	})
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

		var activeBus ActiveBus
		err = json.Unmarshal(msg, &activeBus)
		if err != nil {
			context.String(http.StatusBadRequest, "Could not upgrade to websocket: %v", err)
			return
		}
		busId := context.GetInt64("userId")
		bus, err := models.GetBusById(busId)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{
				"message": "Could not find this bus",
			})
			return
		}

		for i, bus := range allActiveBuses {
			if bus.BusId == busId {
				allActiveBuses = slices.Delete(allActiveBuses, i, i+1)
				break
			}
		}

		activeBus.Name = bus.Name

		activeBus.LastUpdate = time.Now()

		allActiveBuses = append(allActiveBuses, activeBus)

		allActiveBusesJSON, err := json.Marshal(allActiveBuses)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{
				"message": "Bad Request!",
			})
			return
		}

		broadcast <- allActiveBusesJSON
	}
}
