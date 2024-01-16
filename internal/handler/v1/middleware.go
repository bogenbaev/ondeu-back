package v1

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	keycloak "gitlab.com/a5805/ondeu/ondeu-back/pkg/gocloak"
	"gitlab.com/a5805/ondeu/ondeu-back/pkg/modules"
	"io/ioutil"
	"math"
	"net/http"
)

var (
	ErrAccessDenied = "access denied"
	ErrInvalidToken = "token missing required parameters"
	ErrUnauthorized = "you can not perform this action"
)

func authorize(auth keycloak.IKeycloak, roles []string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		valid, claims, err := auth.ValidateToken(ctx, ctx.Request.Header)
		if !valid {
			ctx.JSON(http.StatusBadRequest, gin.H{"reason": err.Error()})
			ctx.Abort()
			return
		}

		logrus.Debugf("claims: %v", claims)

		if len(roles) == 0 {
			roles = []string{"default-roles-ondeu"}
		}

		access, err := auth.CheckAccessToken(ctx, ctx.Request.Header, nil, map[string][]string{"ondeu-front": roles})
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"reason": err.Error()})
			ctx.Abort()
			return
		}

		if !access {
			ctx.JSON(http.StatusUnauthorized, gin.H{"reason": ErrAccessDenied})
			ctx.Abort()
			return
		}

		userId, ok := claims["sub"].(string)
		if !ok && userId == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"reason": ErrInvalidToken})
			ctx.Abort()
			return
		}

		clientID, ok := claims["azp"].(string)
		if !ok && clientID == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"reason": ErrInvalidToken})
			ctx.Abort()
			return
		}

		ctx.Set(modules.ClientID, clientID)
		ctx.Set(modules.UserID, userId)

		ctx.Next()
	}
}

func getRole(c *gin.Context) (string, error) {
	role := c.GetString("role")
	if role == "" {
		return "", errors.New("empty role")
	}
	return role, nil
}

func (h *Handler) adminIdentity(c *gin.Context) {
	role, err := getRole(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"reason": err.Error()})
		return
	}
	if role != modules.Admin {
		c.JSON(http.StatusUnauthorized, gin.H{"reason": ErrUnauthorized})
		return
	}
}

func (h *Handler) managerIdentity(c *gin.Context) {
	role, err := getRole(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"reason": err.Error()})
		return
	}
	if role != modules.Manager {
		c.JSON(http.StatusUnauthorized, gin.H{"reason": ErrUnauthorized})
		return
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,PATCH,OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		}

		c.Next()
	}
}

func ReadRequestBody() gin.HandlerFunc {
	return func(c *gin.Context) {
		var body []byte
		if c.Request.Body != nil {
			body, _ = ioutil.ReadAll(c.Request.Body)
		}
		defer c.Request.Body.Close()
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		if len(body) > 0 {
			logrus.Debugf("Request body: %s", string(body)[:int(math.Min(500, float64(len(string(body)))))])
		}
		c.Next()
	}
}
