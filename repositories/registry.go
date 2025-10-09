package repositories

import (
	paymentRepositories "payment-service/repositories/payment"
	paymentHistoryRepositories "payment-service/repositories/paymentHistory"

	"gorm.io/gorm"
)

type Registry struct {
	db *gorm.DB
}

type IRegistryRepository interface {
	GetPayment() paymentRepositories.IPaymentRepository
	GetPaymentHistory() paymentHistoryRepositories.IPaymentHistoryRepository
}

func NewRegistryRepository(db *gorm.DB) IRegistryRepository {
	return &Registry{db: db}
}

func (r *Registry) GetPayment() paymentRepositories.IPaymentRepository {
	return paymentRepositories.NewPaymentRepository(r.db)
}

func (r *Registry) GetPaymentHistory() paymentHistoryRepositories.IPaymentHistoryRepository {
	return paymentHistoryRepositories.NewPaymentHistory(r.db)
}
