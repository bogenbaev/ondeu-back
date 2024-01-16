package keycloak

import (
	"context"
	"github.com/Nerzal/gocloak/v8"
)

type UserClaim struct {
	Id              int    `json:"id"`
	Lastname        string `json:"lastname"`
	Firstname       string `json:"firstname"`
	Middlename      string `json:"middlename"`
	Firstsignature  bool   `json:"firstsignature"`
	Secondsignature bool   `json:"secondsignature"`
	Iin             string `json:"iin"`
	Username        string `json:"username"`
	Organization    struct {
		CustomerName string `json:"customerName"`
		CustomerId   int    `json:"customerId"`
		Idn          string `json:"idn"`
		ApprovalType string `json:"approvalType"`
		CliCode      string `json:"cliCode"`
		ClientType   string `json:"clientType"`
	} `json:"organization"`
}

//go:generate mockgen -source=keycloak.go -destination=mocks/mock.go

type IClientAuth interface {
	SetClient(clientId string, clientSecret string) IClientAuth
	Auth(ctx context.Context)
	Close()
	GetAccessToken(clientId string) string
}

type IKeycloak interface {
	CheckAccessToken(ctx context.Context, headers map[string][]string, realms []string, resources map[string][]string) (bool, error)
	CheckRoles(ctx context.Context, claim map[string]interface{}, realms []string, resources map[string][]string) (bool, error)
	ValidateToken(ctx context.Context, headers map[string][]string) (bool, map[string]interface{}, error)
	GetUserInfoToken(ctx context.Context, accessToken string) (UserClaim, error)
	GetRoles(ctx context.Context, accessToken, clientID string) ([]*gocloak.Role, error)
}
