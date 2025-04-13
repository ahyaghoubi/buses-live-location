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

func adminChangePassword(context *gin.Context) {
	type passwords struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	var pass passwords
	err := context.ShouldBindJSON(&pass)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "The data could not be parsed!",
		})
		return
	}

	adminId := context.GetInt64("adminId")

	admin, err := models.GetAdminById(adminId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not find the admin!",
		})
		return
	}

	passwordIsValid := utils.ComparePassword(pass.OldPassword, admin.Password)

	if !passwordIsValid {
		context.JSON(http.StatusUnauthorized, gin.H{
			"message": "Old password is incorrect!",
		})
		return
	}

	admin.Password = pass.NewPassword
	err = admin.UpdatePassword()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not update password!",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Password updated successfully!",
	})
}
