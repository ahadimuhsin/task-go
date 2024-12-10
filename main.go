package main

import (
	"net/http"
	"tusk-bwa/config"
	"tusk-bwa/controllers"
	"tusk-bwa/models"

	"github.com/gin-gonic/gin"
)

func main() {
	// config db
	db := config.DatabaseConnection()
	// migrasi
	db.AutoMigrate(&models.User{}, &models.Task{})

	// create new owner (seeder)
	config.CreateOwnerAccount(db)

	// list controller
	userController := controllers.UserController{DB: db}
	taskController := controllers.TaskController{DB: db}

	// router
	router := gin.Default()
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "API is running")
	})

	// users
	router.POST("users/login", userController.Login)
	router.POST("users", userController.Create)
	router.DELETE("users/:id", userController.Delete)
	router.GET("users/Employee", userController.GetEmployee)

	// task
	router.POST("tasks", taskController.Create)
	router.DELETE("tasks/:id", taskController.Delete)
	router.PATCH("tasks/:id/submit", taskController.Submit)
	router.PATCH("tasks/:id/reject", taskController.Reject)
	router.PATCH("tasks/:id/fix", taskController.Fix)
	router.PATCH("tasks/:id/approve", taskController.Approve)
	router.GET("tasks/:id", taskController.FindById)
	router.GET("tasks/review", taskController.Review)
	router.GET("tasks/stats/:userId", taskController.Statistic)
	router.GET("tasks/user/:userId/:status", taskController.FindByUserAndStatus)

	// routing static untuk assets
	router.Static("/assets", "./assets")

	router.Run("localhost:8080")
	// router.Run("172.16.12.136:8080")
}
