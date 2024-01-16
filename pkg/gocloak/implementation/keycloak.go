package implementation

import (
	"context"
	"fmt"
	"github.com/Nerzal/gocloak/v8"
	keycloak2 "gitlab.com/a5805/ondeu/ondeu-back/pkg/gocloak"
	"strings"

	"github.com/mitchellh/mapstructure"
)

type tKeyCloak struct {
	kc    gocloak.GoCloak
	realm string
}

func Keycloak(url string, realm string) keycloak2.IKeycloak {
	return &tKeyCloak{gocloak.NewClient(url), realm}
}

func (k *tKeyCloak) CheckAccessToken(ctx context.Context, headers map[string][]string, realms []string, resources map[string][]string) (bool, error) {
	var (
		result bool
		claims map[string]interface{}
		err    error
	)
	if result, claims, err = k.ValidateToken(ctx, headers); !result {
		return false, err
	}

	if result, err = k.CheckRoles(ctx, claims, realms, resources); !result {
		return false, err
	}

	return true, nil
}

func (k *tKeyCloak) CheckRoles(ctx context.Context, claim map[string]interface{}, realms []string, resources map[string][]string) (bool, error) {
	var (
		err        error
		claims     map[string]interface{}
		resources_ = make(map[string][]string, 0)
		realms_    = make([]string, 0)
	)

	if err = mapstructure.Decode(claim, &claims); err != nil {
		return false, err
	}

	if resources != nil && len(resources) != 0 && claims["resource_access"] != nil {
		for k_, resource := range claims["resource_access"].(map[string]interface{}) {
			for _, v := range resource.(map[string]interface{}) {
				for _, role := range v.([]interface{}) {
					resources_[k_] = append(resources_[k_], role.(string))
				}
			}
		}
	}

	if realms != nil && len(realms) != 0 && claims["realm_access"] != nil {
		for _, v := range claims["realm_access"].(map[string]interface{}) {
			for _, role := range v.([]interface{}) {
				realms_ = append(realms_, role.(string))
			}
		}
	}

	if resources_ != nil && len(resources_) != 0 {
		for client, roles := range resources {
			if resources_[client] != nil && len(resources_[client]) != 0 {
				for _, v := range roles {
					for _, v1 := range resources_[client] {
						if v == v1 {
							return true, nil
						}
					}
				}
			}
		}
	}

	if realms_ != nil && len(realms_) != 0 {
		for _, v := range realms {
			for _, v1 := range realms_ {
				if v == v1 {
					return true, nil
				}
			}
		}
	}

	return false, fmt.Errorf("access denied")
}

func (k *tKeyCloak) ValidateToken(ctx context.Context, headers map[string][]string) (bool, map[string]interface{}, error) {
	var (
		err error
	)

	if len(headers["Authorization"]) == 0 {
		return false, nil, fmt.Errorf("authorization header is not present")
	}

	_, claims, err := k.kc.DecodeAccessToken(ctx, headers["Authorization"][0], k.realm, "")
	if err != nil {
		if strings.Contains(err.Error(), "token is expired") {
			return false, nil, err
		} else {
			return false, nil, err
		}
	}

	return true, *claims, nil
}

func (k *tKeyCloak) GetUserInfoToken(ctx context.Context, accessToken string) (keycloak2.UserClaim, error) {
	var (
		err  error
		user keycloak2.UserClaim
	)

	_, claims, err := k.kc.DecodeAccessToken(ctx, accessToken, k.realm, "")
	if err != nil {
		if strings.Contains(err.Error(), "token is expired") {
			return user, err
		} else {
			return user, err
		}
	}
	c := *claims
	resp := c["provider_response"].(map[string]interface{})
	err = mapstructure.Decode(resp["user_info"], &user)
	return user, err
}

func (k *tKeyCloak) GetRoles(ctx context.Context, accessToken, clientID string) ([]*gocloak.Role, error) {
	roles, err := k.kc.GetClientRoles(ctx, accessToken, k.realm, clientID)
	if err != nil {
		return roles, err
	}

	return roles, err
}
