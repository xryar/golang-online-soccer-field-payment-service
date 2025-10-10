package controllers

import (
	"net/http"
	errValidation "payment-service/common/error"
	"payment-service/common/response"
	"payment-service/domain/dto"
	"payment-service/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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

func (pc *PaymentController) GetAllWithPagination(ctx *gin.Context) {
	var param dto.PaymentRequestParam
	err := ctx.ShouldBindQuery(&param)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResponse{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})

		return
	}

	validate := validator.New()
	if err = validate.Struct(param); err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errValidation.ErrValidationResponse(err)
		response.HttpResponse(response.ParamHTTPResponse{
			Err:     err,
			Code:    http.StatusUnprocessableEntity,
			Message: &errMessage,
			Data:    errResponse,
			Gin:     ctx,
		})

		return
	}

	result, err := pc.service.GetPayment().GetAllWithPagination(ctx, &param)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResponse{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})

		return
	}

	response.HttpResponse(response.ParamHTTPResponse{
		Code: http.StatusOK,
		Data: result,
		Gin:  ctx,
	})

}

func (pc *PaymentController) GetByUUID(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	result, err := pc.service.GetPayment().GetByUUID(ctx, uuid)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResponse{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})

		return
	}

	response.HttpResponse(response.ParamHTTPResponse{
		Code: http.StatusOK,
		Data: result,
		Gin:  ctx,
	})
}

func (pc *PaymentController) Create(ctx *gin.Context) {
	var request dto.PaymentRequest
	err := ctx.ShouldBindJSON(request)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResponse{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})

		return
	}

	validate := validator.New()
	if err = validate.Struct(request); err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errValidation.ErrValidationResponse(err)
		response.HttpResponse(response.ParamHTTPResponse{
			Err:     err,
			Code:    http.StatusUnprocessableEntity,
			Message: &errMessage,
			Data:    errResponse,
			Gin:     ctx,
		})

		return
	}

	result, err := pc.service.GetPayment().Create(ctx, &request)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResponse{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})

		return
	}

	response.HttpResponse(response.ParamHTTPResponse{
		Code: http.StatusCreated,
		Data: result,
		Gin:  ctx,
	})
}

func (pc *PaymentController) Webhook(ctx *gin.Context) {
	var request dto.Webhook
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResponse{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})

		return
	}

	err = pc.service.GetPayment().Webhook(ctx, &request)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResponse{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})

		return
	}

	response.HttpResponse(response.ParamHTTPResponse{
		Code: http.StatusOK,
		Gin:  ctx,
	})
}
