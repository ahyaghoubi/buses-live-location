package routes

import (
	"github.com/ahyaghoubi/buses-live-location/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	server.POST("/admin/login", adminLogin)

	adminAuthenticated := server.Group("/busManagement")
	adminAuthenticated.Use(middleware.Authentication)
	adminAuthenticated.GET("/", getAllBuses)
	adminAuthenticated.POST("/", createBus)
	adminAuthenticated.PATCH("/:id", updateBus)
	adminAuthenticated.DELETE("/:id", deleteBus)
	adminAuthenticated.POST("/login/:id", busLogin)

	server.GET("/bus", getAllActiveBuses)
	server.GET("/clientws", ClientsWSHandler)

	busAuthenticated := server.Group("/bus")
	busAuthenticated.Use(middleware.BusAuthentication)
	busAuthenticated.POST("/", busStatusUpdate)
	busAuthenticated.GET("/ws", BusesWSHandler)

	server.GET("/", func(context *gin.Context) {
		context.File("./public/index.html")
	})
	server.Static("/assets", "./public/")

	server.GET("/admin/login", func(context *gin.Context) {
		context.File("./public/adminLogin.html")
	})

	server.GET("/admin", func(context *gin.Context) {
		context.File("./public/adminDashboard.html")
	})
}
