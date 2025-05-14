package routes

import (
	"din/gopos/controllers"
	"din/gopos/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	auth := controllers.NewAuthController(db)
	r.POST("/signup", auth.SignUp)
	r.POST("login", auth.Login)

	authGroup := r.Group("/")
	authGroup.Use(middlewares.JWTAuthMiddleware())
	{
		authGroup.GET("/me", auth.Profile)
	}

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "pong"})
	})
}
