package main

import (
	"din/gopos/config"
	"din/gopos/models"
	"din/gopos/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()
	db := config.InitDB()
	db.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Transaction{},
		&models.TransactionItem{},
	)
	r := gin.Default()
	routes.SetupRoutes(r, db)
	r.Run()
}
