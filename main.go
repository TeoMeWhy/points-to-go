package main

import (
	"log"
	"points/controllers"
	"points/db"
	"points/models"
	"points/repositories"
	"points/services"

	"github.com/gin-gonic/gin"
)

func main() {

	db, err := db.OpenDBConnection()
	if err != nil {
		log.Fatal("failed to connect database")
	}

	log.Println("Executando migrations...")
	db.AutoMigrate(&models.Customer{}, &models.Transaction{}, &models.TransactionProduct{}, &models.Product{})
	log.Println("ok")

	r := gin.Default()
	r.Use(gin.Recovery())

	pointsRepo := repositories.NewPointsRepo(db)
	pointsService := services.NewPointsService(pointsRepo)
	controller := controllers.NewController(pointsService)

	r.GET("customers/:id", controller.GetCustomerByID)
	r.GET("customers/", controller.GetCustomers)
	r.POST("customers/", controller.PostCustomer)
	r.PUT("customers/:id", controller.PutCustomer)

	r.POST("/transactions", controller.PostTransaction)

	r.PUT("/migrate_customers", controller.MigrateCustomers)

	r.Run("0.0.0.0:8081")

}
