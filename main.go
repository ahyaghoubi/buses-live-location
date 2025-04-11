package main

import (
	"fmt"

	"github.com/ahyaghoubi/buses-live-location/db"
	"github.com/ahyaghoubi/buses-live-location/models"
	"github.com/ahyaghoubi/buses-live-location/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()
	err := models.CreateFirstAdminIfNotExist()
	if err != nil {
		fmt.Println(err)
	}

	go routes.HandleBroadcast()
	go routes.UpdateActiveBusesEvery5min()

	server := gin.Default()

	routes.RegisterRoutes(server)

	server.Run(":8080")
}
