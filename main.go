package app_payment

import (
	"app.payment/src/core/database"
	"app.payment/src/core/router"
	"app.payment/src/core/server"
	"app.payment/src/payments/domain/entities"
	paymentRoutes "app.payment/src/payments/infrastructure/routes"
	"log"
)

func main() {
	// Database connection
	db := database.NewMySQLConnection()

	// Auto-migrate payment entities
	err := db.AutoMigrate(&entities.Payment{})
	if err != nil {
		log.Fatal("Error migrating payment entities:", err)
	}

	// RabbitMQ URL
	rabbitMQURL := "amqp://guest:guest@localhost:5672/"

	// Router setup
	r := router.NewRouter(db)

	// Setup payment routes
	api := r.GetEngine().Group("/api")
	paymentRoutes.SetupPaymentRoutes(api, db, rabbitMQURL)

	// Create server
	srv := server.NewServer("8081", r.GetEngine())

	// Start server
	log.Fatal(srv.Start())
}
