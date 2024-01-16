package implementation

import (
	"context"
	"github.com/Nerzal/gocloak/v8"
	"github.com/sirupsen/logrus"
	keycloak2 "gitlab.com/a5805/ondeu/ondeu-back/pkg/gocloak"
	"strings"
	"time"
)

type TClientJWT struct {
	clientId     string
	clientSecret string
	jwt          *gocloak.JWT
	ticker       *time.Ticker
	duration     time.Duration
}

type TClientAuth struct {
	kc      gocloak.GoCloak
	realm   string
	clients map[string]*TClientJWT
	log     *logrus.Logger
}

func (T *TClientAuth) Auth(ctx context.Context) {
	for k, _ := range T.clients {
		go func(ctx context.Context, client *TClientJWT, key string) {
			client.ticker = time.NewTicker(time.Nanosecond)
			var (
				err         error
				timeExecute time.Time
				logParams   = map[string]interface{}{"event": "clientAuth", "class": "auth", "client": key}
			)
			for {
				timeExecute = time.Now()
				<-client.ticker.C
				if client.jwt, err = T.kc.LoginClient(ctx, client.clientId, client.clientSecret, T.realm); err != nil {
					client.duration = 10 * time.Second
					client.ticker = time.NewTicker(client.duration)
					logParams["duration"] = client.duration.String()
				} else {
					client.duration = time.Duration(client.jwt.ExpiresIn-15) * time.Second
					client.ticker = time.NewTicker(client.duration)
					logParams["duration"] = client.duration.String()
				}
				logParams["timeExecute"] = time.Since(timeExecute).String()
				T.log.WithFields(logParams).Info(err)
			}
		}(ctx, T.clients[k], k)
	}
}

func (T *TClientAuth) Close() {
	for k, _ := range T.clients {
		T.clients[k].ticker.Stop()
	}
}

func (T *TClientAuth) GetAccessToken(clientId string) string {
	clientId_ := strings.ToLower(clientId)
	if T.clients[clientId_] != nil {
		if T.clients[clientId_].jwt != nil {
			return T.clients[clientId_].jwt.AccessToken
		}
	}
	return ""
}

func (T *TClientAuth) SetClient(clientId string, clientSecret string) keycloak2.IClientAuth {
	T.clients[strings.ToLower(clientId)] = &TClientJWT{
		clientId:     clientId,
		clientSecret: clientSecret,
		jwt:          nil,
	}
	return T
}

func ClientAuth(log *logrus.Logger, url string, realm string) keycloak2.IClientAuth {
	return &TClientAuth{log: log, kc: gocloak.NewClient(url), realm: realm, clients: map[string]*TClientJWT{}}
}
