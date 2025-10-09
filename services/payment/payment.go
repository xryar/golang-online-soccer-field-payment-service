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
}

func (ps *PaymentService) GetByUUID(ctx context.Context, uuid string) (*dto.PaymentResponse, error) {}

func (ps *PaymentService) Create(ctx context.Context, req *dto.PaymentRequest) (*dto.PaymentResponse, error) {
}

func (ps *PaymentService) Webhook(ctx context.Context, webhook *dto.Webhook) error {}
