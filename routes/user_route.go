package routes

import (
	"gin-mongo-api/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoute(router *gin.Engine) {
	router.POST("/", controllers.CreateUser())
	router.GET("/", controllers.GetAllUsers())
	router.GET("/user/:userId", controllers.GetaUser())
	router.PUT("/user/:userId", controllers.UpdateUser())
	router.DELETE("/user/:userId", controllers.DeleteUser())
}
