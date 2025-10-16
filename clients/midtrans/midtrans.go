package clients

import (
	errPayment "payment-service/constants/error/payment"
	"payment-service/domain/dto"
	"time"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/sirupsen/logrus"
)

type MidtransClient struct {
	ServerKey    string
	IsProduction bool
}

type IMidtransClient interface {
	CreatePaymentLink(request *dto.PaymentRequest) (*MidtransData, error)
}

func NewMidtransClient(serverKey string, isProduction bool) IMidtransClient {
	return &MidtransClient{
		ServerKey:    serverKey,
		IsProduction: isProduction,
	}
}

func (mc *MidtransClient) CreatePaymentLink(request *dto.PaymentRequest) (*MidtransData, error) {
	var (
		snapClient   snap.Client
		isProduction = midtrans.Sandbox
	)

	expiryDateTime := request.ExpiredAt
	currentTime := time.Now()
	duration := expiryDateTime.Sub(currentTime)
	if duration <= 0 {
		logrus.Errorf("Expired at is invalid")
		return nil, errPayment.ErrExpireAtInvalid
	}

	expiryUnit := "minute"
	expiryDuration := int64(duration.Minutes())

	if duration.Hours() >= 1 {
		expiryUnit = "hour"
		expiryDuration = int64(duration.Hours())
	} else if duration.Hours() >= 24 {
		expiryUnit = "day"
		expiryDuration = int64(duration.Hours() / 24)
	}

	if mc.IsProduction {
		isProduction = midtrans.Production
	}

	snapClient.New(mc.ServerKey, isProduction)
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  request.OrderID,
			GrossAmt: int64(request.Amount),
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: request.CustomerDetail.Name,
			Email: request.CustomerDetail.Email,
			Phone: request.CustomerDetail.Phone,
		},
		Items: &[]midtrans.ItemDetails{
			{
				ID:    request.ItemDetails[0].ID,
				Price: int64(request.ItemDetails[0].Amount),
				Qty:   int32(request.ItemDetails[0].Quantity),
				Name:  request.ItemDetails[0].Name,
			},
		},
		Expiry: &snap.ExpiryDetails{
			Unit:     expiryUnit,
			Duration: expiryDuration,
		},
	}

	response, err := snapClient.CreateTransaction(req)
	if err != nil {
		logrus.Errorf("Error create transaction: %v", err)
		return nil, err
	}

	return &MidtransData{
		RedirectURL: response.RedirectURL,
		Token:       response.Token,
	}, nil
}
