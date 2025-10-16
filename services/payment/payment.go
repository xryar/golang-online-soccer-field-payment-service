package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	clients "payment-service/clients/midtrans"
	"payment-service/config"
	"payment-service/constants"
	"strings"

	"payment-service/common/gcs"
	"payment-service/common/util"
	errPayment "payment-service/constants/error/payment"
	"payment-service/controllers/kafka"
	"payment-service/domain/dto"
	"payment-service/domain/models"
	"payment-service/repositories"
	"time"

	"gorm.io/gorm"
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
	var (
		txErr, err error
		payment    *models.Payment
		response   *dto.PaymentResponse
		midtrans   *clients.MidtransData
	)

	err = ps.repository.GetTx().Transaction(func(tx *gorm.DB) error {
		if !req.ExpiredAt.After(time.Now()) {
			return errPayment.ErrExpireAtInvalid
		}

		midtrans, txErr = ps.midtrans.CreatePaymentLink(req)
		if txErr != nil {
			return txErr
		}

		paymentRequest := &dto.PaymentRequest{
			OrderID:     req.OrderID,
			Amount:      req.Amount,
			Description: req.Description,
			ExpiredAt:   req.ExpiredAt,
			PaymentLink: midtrans.RedirectURL,
		}
		payment, txErr = ps.repository.GetPayment().Create(ctx, tx, paymentRequest)
		if txErr != nil {
			return txErr
		}

		ps.repository.GetPaymentHistory().Create(ctx, tx, &dto.PaymentHistoryRequest{
			PaymentID: payment.ID,
			Status:    payment.Status.GetStatusString(),
		})

		return nil
	})
	if err != nil {
		return nil, err
	}

	response = &dto.PaymentResponse{
		UUID:        payment.UUID,
		OrderID:     payment.OrderID,
		Amount:      payment.Amount,
		Status:      payment.Status.GetStatusString(),
		PaymentLink: payment.PaymentLink,
		Description: payment.Description,
	}

	return response, nil
}

func (ps *PaymentService) Webhook(ctx context.Context, req *dto.Webhook) error {
	var (
		txErr, err         error
		paymentAfterUpdate *models.Payment
		paidAt             *time.Time
		invoiceLink        string
		pdf                []byte
	)

	err = ps.repository.GetTx().Transaction(func(tx *gorm.DB) error {
		_, txErr = ps.repository.GetPayment().FindByOrderID(ctx, req.OrderID.String())
		if txErr != nil {
			return txErr
		}

		if req.TransactionStatus == constants.SettlementString {
			now := time.Now()
			paidAt = &now
		}

		status := req.TransactionStatus.GetStatusInt()
		vaNumber := req.VANumbers[0].VaNumber
		bank := req.VANumbers[0].Bank
		_, txErr = ps.repository.GetPayment().Update(ctx, tx, req.OrderID.String(), &dto.UpdatePaymentRequest{
			TransactionID: &req.TransactionID,
			Status:        &status,
			PaidAt:        paidAt,
			VANumber:      &vaNumber,
			Bank:          &bank,
			Acquirer:      req.Acquirer,
		})
		if txErr != nil {
			return txErr
		}

		paymentAfterUpdate, txErr = ps.repository.GetPayment().FindByOrderID(ctx, req.OrderID.String())
		if txErr != nil {
			return txErr
		}

		txErr = ps.repository.GetPaymentHistory().Create(ctx, tx, &dto.PaymentHistoryRequest{
			PaymentID: paymentAfterUpdate.ID,
			Status:    paymentAfterUpdate.Status.GetStatusString(),
		})

		if req.TransactionStatus == constants.SettlementString {
			paidDay := paidAt.Format("02")
			paidMonth := ps.convertToIndonesiaMonth(paidAt.Format("January"))
			paidYear := paidAt.Format("2006")
			invoiceNumber := fmt.Sprintf("INV/%s/ORD/%d", time.Now().Format(time.DateOnly), ps.randomNumber())
			total := util.RupiahFormat(&paymentAfterUpdate.Amount)
			invoiceRequest := &dto.InvoiceRequest{
				InvoiceNumber: invoiceNumber,
				Data: dto.InvoiceData{
					PaymentDetail: dto.InvoicePaymentDetail{
						PaymentMethod: req.PaymentType,
						BankName:      strings.ToUpper(*paymentAfterUpdate.Bank),
						VANumber:      *paymentAfterUpdate.VANumber,
						Date:          fmt.Sprintf("%s %s %s", paidDay, paidMonth, paidYear),
						IsPaid:        true,
					},
					Items: []dto.InvoiceItem{
						{
							Description: *paymentAfterUpdate.Description,
							Price:       total,
						},
					},
					Total: total,
				},
			}
			pdf, txErr = ps.generatePDF(invoiceRequest)
			if txErr != nil {
				return txErr
			}

			invoiceLink, txErr = ps.UploadToGCS(ctx, invoiceNumber, pdf)
			if txErr != nil {
				return txErr
			}

			_, txErr = ps.repository.GetPayment().Update(ctx, tx, req.OrderID.String(), &dto.UpdatePaymentRequest{
				InvoiceLink: &invoiceLink,
			})
			if txErr != nil {
				return txErr
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	err = ps.produceToKafka(req, paymentAfterUpdate, paidAt)
	if err != nil {
		return err
	}

	return nil
}

func (ps *PaymentService) convertToIndonesiaMonth(englishMonth string) string {
	monthMap := map[string]string{
		"January":   "Januari",
		"February":  "Februari",
		"March":     "Maret",
		"April":     "April",
		"May":       "Mei",
		"June":      "Juni",
		"July":      "Juli",
		"August":    "Agustus",
		"September": "September",
		"October":   "Oktober",
		"November":  "November",
		"Descember": "Desember",
	}
	indonesianMonth, ok := monthMap[englishMonth]
	if !ok {
		return errors.New("month not found").Error()
	}

	return indonesianMonth
}

func (ps *PaymentService) generatePDF(req *dto.InvoiceRequest) ([]byte, error) {
	htmlTemplatePath := "template/invoice.html"
	htmlTemplate, err := os.ReadFile(htmlTemplatePath)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	jsonData, _ := json.Marshal(req)
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		return nil, err
	}

	pdf, err := util.GeneratePDFFromHTML(string(htmlTemplate), data)
	if err != nil {
		return nil, err
	}

	return pdf, nil
}

func (ps *PaymentService) UploadToGCS(ctx context.Context, invoiceNumber string, pdf []byte) (string, error) {
	invoiceNumberReplace := strings.ToLower(strings.ReplaceAll(invoiceNumber, "/", "-"))
	filename := fmt.Sprintf("%s.pdf", invoiceNumberReplace)
	url, err := ps.gcs.UploadFile(ctx, filename, pdf)
	if err != nil {
		return "", err
	}

	return url, nil
}

func (ps *PaymentService) randomNumber() int {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	number := random.Intn(900000 + 100000)
	return number
}

func (ps *PaymentService) mapTransactionStatusToEvent(status constants.PaymentStatusString) string {
	var paymentStatus string
	switch status {
	case constants.PendingString:
		paymentStatus = strings.ToUpper(constants.PendingString.String())
	case constants.SettlementString:
		paymentStatus = strings.ToUpper(constants.SettlementString.String())
	case constants.ExpireString:
		paymentStatus = strings.ToUpper(constants.ExpireString.String())
	}

	return paymentStatus
}

func (ps *PaymentService) produceToKafka(req *dto.Webhook, payment *models.Payment, paidAt *time.Time) error {
	event := dto.KafkaEvent{
		Name: ps.mapTransactionStatusToEvent(req.TransactionStatus),
	}

	metadata := dto.KafkaMetaData{
		Sender:    "payment-service",
		SendingAt: time.Now().Format(time.RFC3339),
	}

	body := dto.KafkaBody{
		Type: "JSON",
		Data: &dto.KafkaData{
			OrderID:   payment.OrderID,
			PaymentID: payment.UUID,
			Status:    req.TransactionStatus.String(),
			PaidAt:    paidAt,
		},
	}

	kafkaMessage := dto.KafkaMessage{
		Event:    event,
		Metadata: metadata,
		Body:     body,
	}

	topic := config.Config.Kafka.Topic
	kafkaMessageJSON, _ := json.Marshal(kafkaMessage)
	err := ps.kafka.GetKafkaProducer().ProduceMessage(topic, kafkaMessageJSON)
	if err != nil {
		return err
	}

	return nil
}
