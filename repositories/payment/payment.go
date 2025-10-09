package repositories

import (
	"context"
	"errors"
	"fmt"
	errWrap "payment-service/common/error"
	"payment-service/constants"
	errConstants "payment-service/constants/error"
	errPayment "payment-service/constants/error/payment"
	"payment-service/domain/dto"
	"payment-service/domain/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentRepository struct {
	db *gorm.DB
}

type IPaymentRepository interface {
	FindAllWithPagination(context.Context, *dto.PaymentRequestParam) ([]models.Payment, int64, error)
	FindByUUID(context.Context, string) (*models.Payment, error)
	FindByOrderID(context.Context, string) (*models.Payment, error)
	Create(context.Context, *gorm.DB, *dto.PaymentRequest) (*models.Payment, error)
	Update(context.Context, *gorm.DB, string, *dto.UpdatePaymentRequest) (*models.Payment, error)
}

func NewPaymentRepository(db *gorm.DB) IPaymentRepository {
	return &PaymentRepository{db: db}
}

func (pr *PaymentRepository) FindAllWithPagination(ctx context.Context, param *dto.PaymentRequestParam) ([]models.Payment, int64, error) {
	var (
		payments []models.Payment
		sort     string
		total    int64
	)
	if param.SortColumn != nil {
		sort = fmt.Sprintf("%s %s", *param.SortColumn, *param.SortOrder)
	} else {
		sort = "created_at desc"
	}

	limit := param.Limit
	offset := (param.Page - 1) * limit
	err := pr.db.
		WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order(sort).
		Find(&payments).
		Error
	if err != nil {
		return nil, 0, errWrap.WrapError(errConstants.ErrSQLError)
	}

	err = pr.db.
		WithContext(ctx).
		Model(&payments).
		Count(&total).
		Error
	if err != nil {
		return nil, 0, errWrap.WrapError(errConstants.ErrSQLError)
	}

	return payments, total, nil
}

func (pr *PaymentRepository) FindByUUID(ctx context.Context, uuid string) (*models.Payment, error) {
	var payment models.Payment
	err := pr.db.WithContext(ctx).Where("uuid = ?", uuid).First(&payment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errWrap.WrapError(errPayment.ErrPaymentNotFound)
		}
		return nil, errWrap.WrapError(errConstants.ErrSQLError)
	}

	return &payment, nil
}

func (pr *PaymentRepository) FindByOrderID(ctx context.Context, orderID string) (*models.Payment, error) {
	var payment models.Payment
	err := pr.db.WithContext(ctx).Where("order_id = ?", orderID).First(&payment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errWrap.WrapError(errPayment.ErrPaymentNotFound)
		}
		return nil, errWrap.WrapError(errConstants.ErrSQLError)
	}

	return &payment, nil
}

func (pr *PaymentRepository) Create(ctx context.Context, tx *gorm.DB, req *dto.PaymentRequest) (*models.Payment, error) {
	status := constants.Initial
	orderID := uuid.MustParse(req.OrderID)
	payment := models.Payment{
		UUID:        uuid.New(),
		OrderID:     orderID,
		Amount:      req.Amount,
		PaymentLink: req.PaymentLink,
		ExpiredAt:   req.ExpiredAt,
		Description: req.Description,
		Status:      &status,
	}

	err := tx.WithContext(ctx).Create(&payment).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstants.ErrSQLError)
	}

	return &payment, nil
}

func (pr *PaymentRepository) Update(ctx context.Context, tx *gorm.DB, order_id string, req *dto.UpdatePaymentRequest) (*models.Payment, error) {
	payment := models.Payment{
		Status:        &req.Status,
		TransactionID: req.TransactionID,
		InvoiceLink:   req.InvoiceLink,
		PaidAt:        req.PaidAt,
		VANumber:      req.VANumber,
		Bank:          req.Bank,
		Acquirer:      req.Acquirer,
	}

	err := tx.WithContext(ctx).Where("order_id", order_id).Updates(&payment).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstants.ErrSQLError)
	}

	return &payment, nil
}
