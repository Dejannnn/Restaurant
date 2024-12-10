package routes

import (
	"github.com/Dejannnn/Restaurant.git/controllers"
	"github.com/gin-gonic/gin"
)

func FoodRoutes(router *gin.Engine) {
	router.GET("/foods", controllers.GetFoods())
	router.GET("/food/:foodId", controllers.GetFood())
	router.POST("/food", controllers.CreateFood())
	router.PATCH("/food/:foodId", controllers.EditFood())
	router.DELETE("/food/:foodId", controllers.DeleteFood())
}
