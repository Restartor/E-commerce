package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"ecommerce/config"
	"ecommerce/routes"
)

func main(){
	gin.SetMode(gin.ReleaseMode) // setmode ada 3 mode : debug, release, test, default debug
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading the .env file!", err)
	}
	config.ConnectDataBase()

	// initialize gin default router
	router := gin.Default()
	
	// setup api routes
	routes.SetupRoutes(router)
	// run the server on port 3000

	router.Run(":3000")
	
}