package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	go HandleBroadcast()

	server := gin.Default()

	server.GET("/", func(context *gin.Context) {
		context.File("./public/index.html")
	})
	server.Static("/assets", "./public/")

	server.GET("/clientws", ClientsWSHandler)
	server.GET("/busws", BusesWSHandler)
	server.POST("/bus", BusesHandler)
	server.GET("/bus", BusGETHandler)

	server.Run(":8080")
}
