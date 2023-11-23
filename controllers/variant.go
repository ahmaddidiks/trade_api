package controllers

import (
	"net/http"
	"trade-api/database"
	"trade-api/helpers"
	"trade-api/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type variantRequest struct {
	VariantName string `form:"variant_name" json:"variant_name" binding:"required"`
	Quantity    uint   `form:"quantity" json:"quantity" binding:"required"`
	ProductID   string `form:"product_id" json:"product_id"`
}

func CreateVariant(ctx *gin.Context) {
	db := database.GetDB()
	var variant models.Variant
	var reqVariant variantRequest
	var product models.Product

	contentType := helpers.GetContentType(ctx)

	if contentType == appJSON {
		if err = ctx.ShouldBindJSON(&reqVariant); err != nil {
			ctx.JSON(http.StatusNotAcceptable, gin.H{
				"error":   "Bad request",
				"message": err.Error()})
			return
		}
	} else {
		if err = ctx.ShouldBind(&reqVariant); err != nil {
			ctx.JSON(http.StatusNotAcceptable, gin.H{
				"error":   "Bad request",
				"message": err.Error()})
			return
		}
	}

	// // search product id by uuid
	err = db.First(&product).Where("uuid = ?", reqVariant.ProductID).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	variant.UUID = uuid.New().String()
	variant.VariantName = reqVariant.VariantName
	variant.Quantity = reqVariant.Quantity
	variant.ProductID = product.ID

	err = db.Debug().Create(&variant).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    variant,
	})
}

func GetAllVariants(ctx *gin.Context) {
	db := database.GetDB()

	results := []models.Variant{}

	err = db.Debug().Preload("Product").Find(&results).Error
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

func UpdateVariantByUUID(ctx *gin.Context) {
	db := database.GetDB()
	var variant models.Variant
	var reqVariant variantRequest

	productUUID := ctx.Param("productUUID")

	// adminData := ctx.MustGet("adminData").(jwt.MapClaims) // only need uuid
	contentType := helpers.GetContentType(ctx)

	if contentType == appJSON {
		if err = ctx.ShouldBindJSON(&reqVariant); err != nil {
			ctx.JSON(http.StatusNotAcceptable, gin.H{
				"error":   "Bad request",
				"message": err.Error()})
			return
		}

	} else {
		if err = ctx.ShouldBind(&reqVariant); err != nil {
			ctx.JSON(http.StatusNotAcceptable, gin.H{
				"error":   "Bad request",
				"message": err.Error()})
			return
		}
	}

	// retrieve product details from db
	err = db.First(&variant).Where("uuid = ?", productUUID).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	variant.VariantName = reqVariant.VariantName
	variant.Quantity = reqVariant.Quantity
	// product.AdminID = uint(adminData["id"].(float64)) // doesnt need to be updated with same value

	err = db.Save(&variant).Error
	if err != nil {
		ctx.JSON(http.StatusNotAcceptable, gin.H{
			"error":   "Bad request",
			"message": err.Error()})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    variant,
	})
}

func DeleteVariantByUUID(ctx *gin.Context) {
	db := database.GetDB()
	var result models.Variant

	productUUID := ctx.Param("productUUID")
	err = db.Where("uuid = ?", productUUID).Delete(&result).Error

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "record not found",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

func GetVariantByUUID(ctx *gin.Context) {
	db := database.GetDB()
	var result models.Variant

	productUUID := ctx.Param("productUUID")
	err = db.First(&result).Where("uuid = ?", productUUID).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}
