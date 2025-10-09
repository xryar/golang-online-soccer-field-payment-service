package repositories

import (
	"context"
	"payment-service/domain/dto"
	"payment-service/domain/models"

	"gorm.io/gorm"
)

type PaymentRepository struct {
	db *gorm.DB
}

type IPaymentRepository interface {
	FindAllWithPagination(context.Context, *dto.PaymentRequestParam) ([]models.Payment, error)
	FindByUUID(context.Context, string) (*models.Payment, error)
	FindByOrderID(context.Context, string) (*models.Payment, error)
	Create(context.Context, *gorm.DB, *dto.PaymentRequest) (*models.Payment, error)
	Update(context.Context, *gorm.DB, string, *dto.UpdatePaymentRequest) (*models.Payment, error)
}

func NewPaymentRepository(db *gorm.DB) IPaymentRepository {
	return &PaymentRepository{db: db}
}

func (pr *PaymentRepository) FindAllWithPagination(context.Context, *dto.PaymentRequestParam) ([]models.Payment, error) {
}

func (pr *PaymentRepository) FindByUUID(context.Context, string) (*models.Payment, error) {
}

func (pr *PaymentRepository) FindByOrderID(context.Context, string) (*models.Payment, error) {
}

func (pr *PaymentRepository) Create(context.Context, *gorm.DB, *dto.PaymentRequest) (*models.Payment, error) {
}

func (pr *PaymentRepository) Update(context.Context, *gorm.DB, string, *dto.UpdatePaymentRequest) (*models.Payment, error) {
}
