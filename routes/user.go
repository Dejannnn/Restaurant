package routes

import (
	"github.com/Dejannnn/Restaurant.git/controllers"
	"github.com/gin-gonic/gin"
)

func UserRouter(router *gin.Engine) {

	router.GET("/users", controllers.GetUsers())
	router.GET("/users/:userId", controllers.GetUser())
	router.POST("/users/signup", controllers.SignUp())
	router.POST("/users/signin", controllers.SignIn())

}
