package controllers

import (
	controllers "payment-service/controllers/http/payment"
	"payment-service/services"
)

type Registry struct {
	service services.IRegistryService
}

type IRegistryController interface {
	GetPayment() controllers.IPaymentController
}

func NewRegistryController(service services.IRegistryService) IRegistryController {
	return &Registry{service: service}
}

func (r *Registry) GetPayment() controllers.IPaymentController {
	return controllers.NewPaymentController(r.service)
}
