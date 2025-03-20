package container

import (
	"log"

	"app.payment/src/payments/application"
	"app.payment/src/payments/domain/repositories"
	"app.payment/src/payments/infrastructure/gateways"
	infraRepo "app.payment/src/payments/infrastructure/repositories"
	"gorm.io/gorm"
)

// Container centraliza la creación y gestión de dependencias
type Container struct {
	db                    *gorm.DB
	paymentRepo           repositories.PaymentRepository
	eventPublisher        repositories.EventPublisher
	paymentGateway        application.PaymentGateway
	processPaymentUseCase *application.ProcessPaymentUseCase
	getPaymentUseCase     *application.GetPaymentUseCase
}

// NewContainer crea un nuevo contenedor
func NewContainer(db *gorm.DB) *Container {
	return &Container{
		db: db,
	}
}

// GetPaymentRepository devuelve el repositorio de pagos
func (c *Container) GetPaymentRepository() repositories.PaymentRepository {
	if c.paymentRepo == nil {
		c.paymentRepo = infraRepo.NewMySQLPaymentRepository(c.db)
	}
	return c.paymentRepo
}

// GetEventPublisher devuelve el publicador de eventos
func (c *Container) GetEventPublisher(rabbitMQURL string) repositories.EventPublisher {
	if c.eventPublisher == nil {
		publisher, err := infraRepo.NewRabbitMQEventPublisher(rabbitMQURL, "payments_exchange")
		if err != nil {
			log.Fatalf("Error creating event publisher: %v", err)
		}
		c.eventPublisher = publisher
	}
	return c.eventPublisher
}

// GetPaymentGateway devuelve el procesador de pagos
func (c *Container) GetPaymentGateway() application.PaymentGateway {
	if c.paymentGateway == nil {
		c.paymentGateway = gateways.NewMockPaymentGateway(0.1)
	}
	return c.paymentGateway
}

// GetProcessPaymentUseCase devuelve el caso de uso de procesar pagos
func (c *Container) GetProcessPaymentUseCase(rabbitMQURL string) *application.ProcessPaymentUseCase {
	if c.processPaymentUseCase == nil {
		c.processPaymentUseCase = application.NewProcessPaymentUseCase(
			c.GetPaymentRepository(),
			c.GetEventPublisher(rabbitMQURL),
			c.GetPaymentGateway(),
		)
	}
	return c.processPaymentUseCase
}

// GetGetPaymentUseCase devuelve el caso de uso de obtener pagos
func (c *Container) GetGetPaymentUseCase() *application.GetPaymentUseCase {
	if c.getPaymentUseCase == nil {
		c.getPaymentUseCase = application.NewGetPaymentUseCase(
			c.GetPaymentRepository(),
		)
	}
	return c.getPaymentUseCase
}

// ConfigureEventConsumer configura el consumidor de eventos
func (c *Container) ConfigureEventConsumer(rabbitMQURL string) error {
	consumer, err := infraRepo.NewRabbitMQEventConsumer(
		rabbitMQURL,
		c.GetProcessPaymentUseCase(rabbitMQURL),
	)
	if err != nil {
		return err
	}

	return consumer.StartConsumingOrderEvents(
		"order_events_queue",
		"orders_exchange",
		"order.created",
	)
}
