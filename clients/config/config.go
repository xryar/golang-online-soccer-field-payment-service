package config

import "github.com/parnurzeal/gorequest"

type ClientConfig struct {
	client       *gorequest.SuperAgent
	baseURL      string
	signatureKey string
}

type IClientConfig interface {
	Client() *gorequest.SuperAgent
	BaseURL() string
	SignatureKey() string
}

type Option func(*ClientConfig)

func NewClientConfig(options ...Option) IClientConfig {
	clientConfig := &ClientConfig{
		client: gorequest.New().
			Set("Content-Type", "application/json").
			Set("Accept", "application/json"),
	}

	for _, option := range options {
		option(clientConfig)
	}

	return clientConfig
}

func (cc *ClientConfig) Client() *gorequest.SuperAgent {
	return cc.client
}

func (cc *ClientConfig) BaseURL() string {
	return cc.baseURL
}

func (cc *ClientConfig) SignatureKey() string {
	return cc.signatureKey
}

func WithBaseURL(baseURL string) Option {
	return func(cc *ClientConfig) {
		cc.baseURL = baseURL
	}
}

func WithSignatureKey(signatureKey string) Option {
	return func(cc *ClientConfig) {
		cc.signatureKey = signatureKey
	}
}
