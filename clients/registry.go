package clients

import (
	fieldConfig "payment-service/clients/config"
	userConfig "payment-service/clients/user"
	"payment-service/config"
)

type RegistryClient struct{}

type IRegistryClient interface {
	GetUser() userConfig.IUserClient
}

func NewRegistryClient() IRegistryClient {
	return &RegistryClient{}
}

func (rc *RegistryClient) GetUser() userConfig.IUserClient {
	return userConfig.NewUserClient(
		fieldConfig.NewClientConfig(
			fieldConfig.WithBaseURL(config.Config.InternalService.User.Host),
			fieldConfig.WithSignatureKey(config.Config.InternalService.User.SignatureKey),
		),
	)
}
