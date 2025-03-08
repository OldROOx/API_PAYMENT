package repositories

import (
	"app.payment/src/payments/domain/entities"
	"gorm.io/gorm"
)

type MySQLPaymentRepository struct {
	db *gorm.DB
}

func NewMySQLPaymentRepository(db *gorm.DB) *MySQLPaymentRepository {
	return &MySQLPaymentRepository{db: db}
}

func (r *MySQLPaymentRepository) Save(payment *entities.Payment) error {
	return r.db.Create(payment).Error
}

func (r *MySQLPaymentRepository) FindByID(id string) (*entities.Payment, error) {
	var payment entities.Payment
	result := r.db.First(&payment, "id = ?", id)
	return &payment, result.Error
}

func (r *MySQLPaymentRepository) FindByOrderID(orderID uint) (*entities.Payment, error) {
	var payment entities.Payment
	result := r.db.First(&payment, "order_id = ?", orderID)
	return &payment, result.Error
}

func (r *MySQLPaymentRepository) UpdateStatus(id string, status entities.PaymentStatus, transactionID string) error {
	return r.db.Model(&entities.Payment{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":         status,
			"transaction_id": transactionID,
		}).Error
}

func (r *MySQLPaymentRepository) UpdatePaymentFailed(id string, errorMessage string) error {
	return r.db.Model(&entities.Payment{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":        entities.PaymentStatusFailed,
			"error_message": errorMessage,
		}).Error
}
