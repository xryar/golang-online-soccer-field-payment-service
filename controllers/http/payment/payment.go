package controllers

import (
	services "payment-service/services"

	"github.com/gin-gonic/gin"
)

type PaymentController struct {
	service services.IRegistryService
}

type IPaymentController interface {
	GetAllWithPagination(*gin.Context)
	GetByUUID(*gin.Context)
	Create(*gin.Context)
	Webhook(*gin.Context)
}

func NewPaymentController(service services.IRegistryService) IPaymentController {
	return &PaymentController{service: service}
}

func (pc *PaymentController) GetAllWithPagination(*gin.Context) {}

func (pc *PaymentController) GetByUUID(*gin.Context) {}

func (pc *PaymentController) Create(*gin.Context) {}

func (pc *PaymentController) Webhook(*gin.Context) {}
