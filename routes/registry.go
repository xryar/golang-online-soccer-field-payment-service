package routes

import (
	"payment-service/clients"
	controllers "payment-service/controllers/http"
	routes "payment-service/routes/payment"

	"github.com/gin-gonic/gin"
)

type Registry struct {
	controller controllers.IRegistryController
	group      *gin.RouterGroup
	client     clients.IRegistryClient
}

type IRegistryRoute interface {
	Serve()
}

func NewRouteRegistry(controller controllers.IRegistryController, group *gin.RouterGroup, client clients.IRegistryClient) IRegistryRoute {
	return &Registry{
		controller: controller,
		group:      group,
		client:     client,
	}
}

func (r *Registry) Serve() {
	r.paymentRoute().Run()
}

func (r *Registry) paymentRoute() routes.IPaymentRoute {
	return routes.NewPaymentRoute(r.group, r.controller, r.client)
}
