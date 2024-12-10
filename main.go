package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Dejannnn/Restaurant.git/middleware"
	"github.com/Dejannnn/Restaurant.git/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	// initialization code
	println("Call this")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

}
func main() {
	fmt.Println("Start develop restorant app")
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8000"
	}

	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRouter(router)
	router.Use(middleware.Authentification())

	routes.FoodRoutes(router)
	routes.MenuRoutes(router)
	router.Run(port)

}
