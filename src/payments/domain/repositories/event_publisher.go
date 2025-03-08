package repositories

import (
	"app.payment/src/payments/domain/entities"
)

type EventPublisher interface {
	PublishEvent(event entities.Event) error
	Close() error
}
