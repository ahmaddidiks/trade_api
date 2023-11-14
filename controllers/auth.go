package controllers

import (
	"net/http"
	"trade-api/database"
	"trade-api/helpers"
	"trade-api/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	appJSON = "application/json"
)

func RegisterAdmin(ctx *gin.Context) {
	db := database.GetDB()
	contextType := helpers.GetContentType(ctx)
	admin := models.Admin{}

	if contextType == appJSON {
		ctx.ShouldBindJSON(&admin)
	} else {
		ctx.ShouldBind(&admin)
	}

	// uuid
	admin.UUID = uuid.New().String()

	err := db.Debug().Create(&admin).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    admin,
	})
}

func LoginAdmin(ctx *gin.Context) {
	db := database.GetDB()
	contentType := helpers.GetContentType(ctx)
	admin := models.Admin{}

	if contentType == appJSON {
		ctx.ShouldBindJSON(&admin)
	} else {
		ctx.ShouldBind(&admin)
	}

	password := admin.Password

	err := db.Debug().Where("email = ?", admin.Email).Take(&admin).Error
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Invalid emali",
		})
		return
	}

	comparePass := helpers.ComparePass([]byte(admin.Password), []byte(password))

	if !comparePass {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Invalid password",
		})
		return
	}

	token := helpers.GenerateToken(admin.ID, admin.Email)

	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
