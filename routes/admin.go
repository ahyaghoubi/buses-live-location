package routes

import (
	"fmt"
	"net/http"

	"github.com/ahyaghoubi/buses-live-location/models"
	"github.com/ahyaghoubi/buses-live-location/utils"
	"github.com/gin-gonic/gin"
)

func adminLogin(context *gin.Context) {
	var admin models.Admin

	err := context.ShouldBindJSON(&admin)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "The data could not be Parsed!",
		})
		return
	}

	err = admin.ValidateCredentials()
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{
			"message": "Email or Password was incorrect!",
			"error":   fmt.Sprintf("%v", err),
		})
		return
	}

	token, err := utils.GenerateToken(admin.Email, admin.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not authenticate!",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Login successful!",
		"token":   token,
	})
}

// func adminChangePassword(context *gin.Context) {
// 	var admin models.Admin
// 	err := context.ShouldBindJSON(&admin)
// 	if err != nil {
// 		context.JSON(http.StatusInternalServerError, gin.H{})
// 	}

// 	err = admin.ValidateCredentials()
// 	if err != nil {
// 		context.JSON(http.StatusUnauthorized, gin.H{
// 			"message": "Current password in incorrect",
// 			"error":   fmt.Sprintf("%v", err),
// 		})
// 		return
// 	}
// }

// func adminSignup(context *gin.Context) {
// 	var admin models.Admin
// 	err := context.ShouldBindJSON(&admin)
// 	if err != nil {
// 		context.JSON(http.StatusInternalServerError, gin.H{
// 			"message": "The data could not be parsed!",
// 		})
// 		return
// 	}

// 	if admin.Email == "" || admin.Name == "" || admin.Password == "" {
// 		context.JSON(http.StatusBadRequest, gin.H{
// 			"message": "Name or email or password",
// 		})
// 	}

// 	err = admin.Create()
// 	if err != nil {
// 		context.JSON(http.StatusInternalServerError, gin.H{
// 			"message": "Something went wrong!",
// 		})
// 		return
// 	}

// 	context.JSON(http.StatusCreated, gin.H{
// 		"message": "Admin created.",
// 	})
// }
