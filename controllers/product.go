package controllers

import (
	"din/gopos/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductController struct{}

func NewProductController() *ProductController {
	return &ProductController{}
}

func (p *ProductController) GetAll(ctx *gin.Context) {
	db := ctx.MustGet("db").(*gorm.DB)
	var products []models.Product
	db.Find(&products)
	ctx.JSON(http.StatusOK, gin.H{"data": products})
}

func (p *ProductController) Create(ctx *gin.Context) {
	db := ctx.MustGet("db").(*gorm.DB)
	var input models.Product

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Name == "" || input.Price <= 0 || input.Stock < 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid product  data"})
		return
	}

	db.Create(&input)
	ctx.JSON(http.StatusCreated, gin.H{"data": input})
}

func (p *ProductController) Update(ctx *gin.Context) {
	db := ctx.MustGet("db").(*gorm.DB)
	id, _ := strconv.Atoi(ctx.Param("id"))
	var product models.Product

	if err := db.First(&product, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	var input models.Product
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Name == "" || input.Price <= 0 || input.Stock < 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid product data"})
		return
	}

	db.Model(&product).Updates(input)
	ctx.JSON(http.StatusOK, gin.H{"data": product})
}

func (p *ProductController) Delete(ctx *gin.Context) {
	db := ctx.MustGet("db").(*gorm.DB)
	id, _ := strconv.Atoi(ctx.Param("id"))
	var product models.Product

	if err := db.First(&product, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	db.Delete(&product)
	ctx.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
