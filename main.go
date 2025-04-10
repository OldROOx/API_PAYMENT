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

	db := database.NewMySQLConnection()

	err := db.AutoMigrate(&entities.Payment{})
	if err != nil {
		log.Fatal("Error migrating payment entities:", err)
	}

	rabbitMQURL := "amqp://guest:guest@54.84.215.25:5672/"

	r := router.NewRouter(db)

	api := r.GetEngine().Group("/api")
	paymentRoutes.SetupPaymentRoutes(api, db, rabbitMQURL)

	srv := server.NewServer("8081", r.GetEngine())

	log.Fatal(srv.Start())
}
