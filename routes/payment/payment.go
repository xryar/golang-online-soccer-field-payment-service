package routes

import (
	"payment-service/clients"
	"payment-service/constants"
	controllers "payment-service/controllers/http"
	"payment-service/middlewares"

	"github.com/gin-gonic/gin"
)

type PaymentRoute struct {
	controller controllers.IRegistryController
	client     clients.IRegistryClient
	group      *gin.RouterGroup
}

type IPaymentRoute interface {
	Run()
}

func NewPaymentRoute(group *gin.RouterGroup, controller controllers.IRegistryController, client clients.IRegistryClient) IPaymentRoute {
	return &PaymentRoute{
		group:      group,
		controller: controller,
		client:     client,
	}
}

func (pr *PaymentRoute) Run() {
	group := pr.group.Group("/payment")
	group.POST("/webhook", pr.controller.GetPayment().Webhook)
	group.Use(middlewares.Authenticate())
	group.GET("", middlewares.CheckRole([]string{
		constants.Admin,
		constants.Customer,
	}, pr.client), pr.controller.GetPayment().GetAllWithPagination)
	group.GET("/uuid", middlewares.CheckRole([]string{
		constants.Admin,
		constants.Customer,
	}, pr.client), pr.controller.GetPayment().GetByUUID)
	group.POST("", middlewares.CheckRole([]string{
		constants.Customer,
	}, pr.client), pr.controller.GetPayment().Create)
}
