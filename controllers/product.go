package controllers

import (
	"din/gopos/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductController struct{}

// class function
func NewProductController() *ProductController {
	return &ProductController{}
}

// GET /products
func (p *ProductController) GetAll(ctx *gin.Context) {
	db := ctx.MustGet("db").(*gorm.DB)

	// Query Params
	search := ctx.Query("search")
	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	offset := (page - 1) * limit

	var products []models.Product
	query := db.Model(&models.Product{})

	// search pagination (limit default 10)
	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	// filter by stock min / max
	minStockStr := ctx.Query("min_stock")
	maxStockStr := ctx.Query("max_stock")
	if minStockStr != "" {
		if minStock, err := strconv.Atoi(minStockStr); err == nil {
			query = query.Where("stock >= ?", minStock)
		}
	}
	if maxStockStr != "" {
		if maxStock, err := strconv.Atoi(maxStockStr); err == nil {
			query = query.Where("stock <= ?", maxStock)
		}
	}

	// filter min_price, max_price
	priceMinStr := ctx.Query("min_price")
	priceMaxStr := ctx.Query("max_price")
	if priceMinStr != "" {
		if priceMin, err := strconv.ParseFloat(priceMinStr, 64); err == nil {
			query = query.Where("price >= ?", priceMin)
		}
	}
	if priceMaxStr != "" {
		if priceMax, err := strconv.ParseFloat(priceMaxStr, 64); err == nil {
			query = query.Where("price <= ?", priceMax)
		}
	}

	// sort name and price
	sort := ctx.DefaultQuery("sort", "id_asc")
	switch sort {
	case "id_desc":
		query = query.Order("id DESC")
	case "name_asc":
		query = query.Order("name ASC")
	case "name_desc":
		query = query.Order("name DESC")
	case "price_asc":
		query = query.Order("price ASC")
	case "price_desc":
		query = query.Order("price DESC")
	default:
		query = query.Order("id ASC")
	}

	query.Offset(offset).Limit(limit).Find(&products)

	// total products
	var total int64
	query.Count(&total)

	ctx.JSON(http.StatusOK, gin.H{
		"data":  products,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// POST /product
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

// PUT /product/:id
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

// DELETE /product/:id
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

func (p *ProductController) GetByID(ctx *gin.Context) {
	db := ctx.MustGet("db").(*gorm.DB)
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var product models.Product
	if err := db.First(&product, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": product})
}

func (p *ProductController) Patch(ctx *gin.Context) {
	db := ctx.MustGet("db").(*gorm.DB)
	id, _ := strconv.Atoi(ctx.Param("id"))
	var product models.Product
	if err := db.First(&product, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	var input map[string]interface{}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if name, ok := input["name"].(string); ok && strings.TrimSpace(name) == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "name cannot be empty"})
		return
	}
	if price, ok := input["price"].(float64); ok && price <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "price must be greater than 0"})
		return
	}
	if stock, ok := input["stock"].(float64); ok && int(stock) < 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "stock must be >= 0"})
		return
	}

	db.Model(&product).Updates(input)
	ctx.JSON(http.StatusOK, gin.H{"data": product})
}
