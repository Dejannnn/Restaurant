package routes

import (
	"github.com/Dejannnn/Restaurant.git/controllers"
	"github.com/gin-gonic/gin"
)

func MenuRoutes(router *gin.Engine) {
	router.GET("/menus", controllers.GetMenus())
	router.GET("/menus/:menuId", controllers.GetMenu())
	router.POST("/menus", controllers.CreateMenu())
	router.PATCH("/menus/:menuId", controllers.UpdateMenu())
	router.DELETE("/menus/:menuId", controllers.DeleteMenu())
}
