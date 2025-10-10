package services

import (
	clients "payment-service/clients/midtrans"
	"payment-service/common/gcs"
	"payment-service/controllers/kafka"
	"payment-service/repositories"
	services "payment-service/services/payment"
)

type Registry struct {
	repository repositories.IRegistryRepository
	gcs        gcs.IGCSClient
	kafka      kafka.IKafkaRegistry
	midtrans   clients.IMidtransClient
}

type IRegistryService interface {
	GetPayment() services.IPaymentService
}

func NewRegistry(repository repositories.IRegistryRepository, gcs gcs.IGCSClient, kafka kafka.IKafkaRegistry, midtrans clients.IMidtransClient) IRegistryService {
	return &Registry{
		repository: repository,
		gcs:        gcs,
		kafka:      kafka,
		midtrans:   midtrans,
	}
}

func (r *Registry) GetPayment() services.IPaymentService {
	return services.NewPaymentService(r.repository, r.gcs, r.kafka, r.midtrans)
}
