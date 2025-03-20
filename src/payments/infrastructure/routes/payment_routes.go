package routes

import (
	"app.payment/src/core/container"
	"app.payment/src/payments/infrastructure/controllers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupPaymentRoutes(api *gin.RouterGroup, db *gorm.DB, rabbitMQURL string) {
	// Inicializar contenedor de dependencias
	c := container.NewContainer(db)

	// Inicializar consumidor de eventos
	if err := c.ConfigureEventConsumer(rabbitMQURL); err != nil {
		panic(err)
	}

	// Inicializar controladores
	processPaymentController := controllers.NewProcessPaymentController(
		c.GetProcessPaymentUseCase(rabbitMQURL),
	)
	getPaymentController := controllers.NewGetPaymentController(
		c.GetGetPaymentUseCase(),
	)

	// Configurar rutas
	payments := api.Group("/payments")
	{
		payments.POST("", processPaymentController.Handle)
		payments.GET("/:id", getPaymentController.HandleGetByID)
		payments.GET("/order/:orderID", getPaymentController.HandleGetByOrderID)
	}
}
