package routes

import (
	"app.payment/src/payments/application"
	"app.payment/src/payments/infrastructure/controllers"
	"app.payment/src/payments/infrastructure/gateways"
	"app.payment/src/payments/infrastructure/repositories"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupPaymentRoutes(api *gin.RouterGroup, db *gorm.DB, rabbitMQURL string) {
	// Initialize repositories
	paymentRepo := repositories.NewMySQLPaymentRepository(db)
	eventPublisher, err := repositories.NewRabbitMQEventPublisher(rabbitMQURL, "payments_exchange")
	if err != nil {
		panic(err)
	}

	// Initialize payment gateway
	// For production, you would use a real payment gateway
	// For now, we use a mock with 10% failure rate
	paymentGateway := gateways.NewMockPaymentGateway(0.1)

	// Initialize use cases
	processPaymentUseCase := application.NewProcessPaymentUseCase(paymentRepo, eventPublisher, paymentGateway)
	getPaymentUseCase := application.NewGetPaymentUseCase(paymentRepo)

	// Initialize event consumer
	eventConsumer, err := repositories.NewRabbitMQEventConsumer(rabbitMQURL, processPaymentUseCase)
	if err != nil {
		panic(err)
	}

	// Start consuming order events
	err = eventConsumer.StartConsumingOrderEvents(
		"order_events_queue",
		"orders_exchange",
		"order.created",
	)
	if err != nil {
		panic(err)
	}

	// Initialize controllers
	processPaymentController := controllers.NewProcessPaymentController(processPaymentUseCase)
	getPaymentController := controllers.NewGetPaymentController(getPaymentUseCase)

	// Setup routes
	payments := api.Group("/payments")
	{
		payments.POST("", processPaymentController.Handle)
		payments.GET("/:id", getPaymentController.HandleGetByID)
		payments.GET("/order/:orderID", getPaymentController.HandleGetByOrderID)
	}
}
