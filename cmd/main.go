package main

import (
	"din/gopos/config"
	"din/gopos/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()
	db := config.InitDB()
	r := gin.Default()
	routes.SetupRoutes(r, db)
	r.Run()
}
