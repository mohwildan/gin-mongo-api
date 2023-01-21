package main

import (
	"gin-mongo-api/configs"
	"gin-mongo-api/routes"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	configs.ConnectDB()

	routes.UserRoute(router)
	router.Run(":" + os.Getenv("PORT"))
}
