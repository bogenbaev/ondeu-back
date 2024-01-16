package information

import (
	"context"
	"fmt"
	"github.com/Nerzal/gocloak/v8"
	"github.com/sirupsen/logrus"
	keycloak "gitlab.com/a5805/ondeu/ondeu-back/pkg/gocloak"
	"gitlab.com/a5805/ondeu/ondeu-back/pkg/modules"
)

type Service struct {
	cfg *modules.Keycloak
	kc  keycloak.IKeycloak
}

func NewService(cfg *modules.Keycloak, kc keycloak.IKeycloak) *Service {
	return &Service{
		cfg: cfg,
		kc:  kc,
	}
}

func (s *Service) GetRoles(ctx context.Context) ([]*gocloak.Role, error) {
	client := gocloak.NewClient(s.cfg.Host)
	token, err := client.LoginClient(ctx, s.cfg.AdminClientID, s.cfg.AdminClientSecret, s.cfg.Realm)
	if err != nil {
		logrus.Errorf("[internal service error] - %+v", err)
		return nil, err
	}

	clientID, ok := ctx.Value(modules.ClientID).(string)
	if !ok {
		return nil, fmt.Errorf("clientID not found in context")
	}

	clients, err := client.GetClients(ctx, token.AccessToken, s.cfg.Realm, gocloak.GetClientsParams{ClientID: &clientID})
	if err != nil {
		logrus.Errorf("[internal service error] - %+v", err)
		return nil, err
	}

	roles, err := s.kc.GetRoles(ctx, token.AccessToken, *clients[0].ID)
	if err != nil {
		logrus.Errorf("[internal service error] - %+v", err)
		return nil, err
	}

	return roles, nil
}
