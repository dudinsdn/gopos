package routes

import (
	"din/gopos/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	auth := controllers.NewAuthController(db)
	r.POST("/signup", auth.SignUp)
	r.POST("login", auth.Login)

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "pong"})
	})
}
