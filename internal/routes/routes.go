package routes

import (
	"awesomeGo/internal/controller"
	"github.com/gin-gonic/gin"
)

// InitRoutes initializes the API routes
func InitRoutes() *gin.Engine {
	router := gin.Default()

	// Define user-related routes
	router.POST("/users/createUser", controller.CreateUser)
	router.GET("/users/getUsers", controller.GetUsers)
	router.GET("/users/getUser/:id", controller.GetUser)
	router.PUT("/users/updateUser/:id", controller.UpdateUser)
	router.DELETE("/users/deleteUser/:id", controller.DeleteUser)
	router.GET("/users/downloadCSV", controller.DownloadCSV)
	router.POST("/users/importCSV", controller.ImportCSV)

	return router
}
