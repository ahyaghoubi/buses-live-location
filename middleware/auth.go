package middleware

import (
	"fmt"
	"net/http"

	"github.com/ahyaghoubi/buses-live-location/utils"
	"github.com/gin-gonic/gin"
)

func Authentication(context *gin.Context) {
	token := context.Request.Header.Get("Authorization")
	if token == "" {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Unauthorized because there is no token!",
		})
		return
	}

	adminId, err := utils.VerifyToken(token)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized!",
			"error":   fmt.Sprintf("%v", err),
		})
		return
	}

	context.Set("adminId", adminId)

	context.Next()
}

func BusAuthentication(context *gin.Context) {
	token := context.Request.Header.Get("Authorization")

	if token == "" {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Unauthorized because there is no token!",
		})
		return
	}

	BusId, err := utils.VerifyBusToken(token)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized!",
			"error":   fmt.Sprintf("%v", err),
		})
		return
	}

	context.Set("busId", BusId)

	context.Next()
}
