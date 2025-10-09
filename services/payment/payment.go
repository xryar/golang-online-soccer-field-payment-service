package services

import (
	"context"
	clients "payment-service/clients/midtrans"
	"payment-service/common/gcs"
	"payment-service/common/util"
	"payment-service/controllers/kafka"
	"payment-service/domain/dto"
	"payment-service/repositories"
)

type PaymentService struct {
	repository repositories.IRegistryRepository
	gcs        gcs.IGCSClient
	kafka      kafka.IKafkaRegistry
	midtrans   clients.IMidtransClient
}

type IPaymentService interface {
	GetAllWithPagination(context.Context, *dto.PaymentRequestParam) (*util.PaginationResult, error)
	GetByUUID(context.Context, string) (*dto.PaymentResponse, error)
	Create(context.Context, *dto.PaymentRequest) (*dto.PaymentResponse, error)
	Webhook(context.Context, *dto.Webhook) error
}

func NewPaymentService(repository repositories.IRegistryRepository, gcs gcs.IGCSClient, kafka kafka.IKafkaRegistry, midtrans clients.IMidtransClient) IPaymentService {
	return &PaymentService{
		repository: repository,
		gcs:        gcs,
		kafka:      kafka,
		midtrans:   midtrans,
	}
}

func (ps *PaymentService) GetAllWithPagination(ctx context.Context, param *dto.PaymentRequestParam) (*util.PaginationResult, error) {
	payments, total, err := ps.repository.GetPayment().FindAllWithPagination(ctx, param)
	if err != nil {
		return nil, err
	}

	paymentResult := make([]dto.PaymentResponse, 0, len(payments))
	for _, payment := range payments {
		paymentResult = append(paymentResult, dto.PaymentResponse{
			UUID:          payment.UUID,
			TransactionID: payment.TransactionID,
			OrderID:       payment.OrderID,
			Amount:        payment.Amount,
			Status:        payment.Status.GetStatusString(),
			PaymentLink:   payment.PaymentLink,
			InvoiceLink:   payment.InvoiceLink,
			VANumber:      payment.VANumber,
			Bank:          payment.Bank,
			Description:   payment.Description,
			ExpiredAt:     payment.ExpiredAt,
			CreatedAt:     payment.CreatedAt,
			UpdatedAt:     payment.UpdatedAt,
		})
	}

	paginationParam := util.PaginationParam{
		Page:  param.Page,
		Limit: param.Limit,
		Count: total,
		Data:  paymentResult,
	}

	response := util.GeneratePagination(paginationParam)

	return &response, nil
}

func (ps *PaymentService) GetByUUID(ctx context.Context, uuid string) (*dto.PaymentResponse, error) {
	payment, err := ps.repository.GetPayment().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	return &dto.PaymentResponse{
		UUID:          payment.UUID,
		TransactionID: payment.TransactionID,
		OrderID:       payment.OrderID,
		Amount:        payment.Amount,
		Status:        payment.Status.GetStatusString(),
		PaymentLink:   payment.PaymentLink,
		InvoiceLink:   payment.InvoiceLink,
		VANumber:      payment.VANumber,
		Bank:          payment.Bank,
		Description:   payment.Description,
		ExpiredAt:     payment.ExpiredAt,
		CreatedAt:     payment.CreatedAt,
		UpdatedAt:     payment.UpdatedAt,
	}, nil
}

func (ps *PaymentService) Create(ctx context.Context, req *dto.PaymentRequest) (*dto.PaymentResponse, error) {
}

func (ps *PaymentService) Webhook(ctx context.Context, webhook *dto.Webhook) error {}
