package controllers

import (
	"mime/multipart"
	"net/http"
	"trade-api/database"
	"trade-api/helpers"
	"trade-api/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type productRequest struct {
	Name  string                `form:"name" binding:"required"`
	Image *multipart.FileHeader `form:"file"`
}

var (
	err error
)

func CreateProduct(ctx *gin.Context) {
	db := database.GetDB()
	var product models.Product
	var reqProduct productRequest

	adminData := ctx.MustGet("adminData").(jwt.MapClaims)
	contentType := helpers.GetContentType(ctx)

	if contentType == appJSON {
		ctx.ShouldBindJSON(&reqProduct)
	} else {
		ctx.ShouldBind(&reqProduct)
	}

	if err = ctx.ShouldBind(&reqProduct); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract the filename without extension
	fileName := helpers.RemoveExtension(reqProduct.Image.Filename)

	uploadResult, err := helpers.UploadFile(reqProduct.Image, fileName)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product.Name = reqProduct.Name
	product.ImageURL = uploadResult
	product.AdminID = uint(adminData["id"].(float64))
	product.UUID = uuid.New().String()

	err = db.Debug().Create(&product).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    product,
	})
}

func GetAllProducts(ctx *gin.Context) {
	db := database.GetDB()

	results := []models.Product{}

	// err = db.Debug().Preload("Admin").Find(&results).Error
	err = db.Find(&results).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": results,
	})
}

func UpdateProductByUUID(ctx *gin.Context) {
	db := database.GetDB()
	var product models.Product
	var reqProduct productRequest

	productUUID := ctx.Param("productUUID")

	// adminData := ctx.MustGet("adminData").(jwt.MapClaims) // only need uuid
	contentType := helpers.GetContentType(ctx)

	if contentType == appJSON {
		if err = ctx.ShouldBindJSON(&reqProduct); err != nil {
			ctx.JSON(http.StatusNotAcceptable, gin.H{
				"error":   "Bad request",
				"message": err.Error()})
			return
		}

	} else {
		if err = ctx.ShouldBind(&reqProduct); err != nil {
			ctx.JSON(http.StatusNotAcceptable, gin.H{
				"error":   "Bad request",
				"message": err.Error()})
			return
		}
	}

	// retrieve product details from db
	err = db.First(&product).Where("uuid = ?", productUUID).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	// check if received a new image, then upload to cludinary and update product image url
	if reqProduct.Image != nil {
		// Extract the filename without extension
		fileName := helpers.RemoveExtension(reqProduct.Image.Filename)

		uploadResult, err := helpers.UploadFile(reqProduct.Image, fileName)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		product.ImageURL = uploadResult
	}

	product.Name = reqProduct.Name
	// product.AdminID = uint(adminData["id"].(float64)) // doesnt need to be updated with same value

	err = db.Save(&product).Error
	if err != nil {
		ctx.JSON(http.StatusNotAcceptable, gin.H{
			"error":   "Bad request",
			"message": err.Error()})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    product,
	})

}

func GetProductByUUID(ctx *gin.Context) {
	db := database.GetDB()
	var product models.Product

	productUUID := ctx.Param("productUUID")
	err = db.First(&product).Where("uuid = ?", productUUID).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    product,
	})
}

func DeteleProductByUUID(ctx *gin.Context) {

}
