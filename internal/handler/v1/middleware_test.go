package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	servicemocks "gitlab.com/a5805/ondeu/ondeu-back/pkg/gocloak/mocks"
	"gitlab.com/a5805/ondeu/ondeu-back/pkg/modules"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_getRole(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
		err   error
	}{
		{
			name:  "Success.",
			input: "test",
			want:  "test",
			err:   nil,
		},
		{
			name:  "Failed. No role in context.",
			input: "",
			want:  "",
			err:   errors.New("empty role"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			testCtx, _ := gin.CreateTestContext(w)
			testCtx.Set("role", tt.input)

			got, err := getRole(testCtx)
			assert.Equalf(t, tt.want, got, "getRole()")
			assert.Equalf(t, tt.err, err, "getRole()")
		})
	}
}

func Test_authorize(t *testing.T) {
	type mockBehavior func(recorder *servicemocks.MockIKeycloak)

	userId := uuid.New().String()
	realm := "test-realm"

	tests := []struct {
		name         string
		mockBehavior mockBehavior
		roles        []string
		userId       string
		realm        string
		wantCode     int
		wantMessage  string
	}{
		{
			name: "Success.",
			mockBehavior: func(recorder *servicemocks.MockIKeycloak) {
				recorder.EXPECT().
					ValidateToken(gomock.Any(), gomock.Any()).
					Return(true, map[string]interface{}{"sub": userId, "azp": realm}, nil).
					AnyTimes()
				recorder.EXPECT().
					CheckAccessToken(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(true, nil).
					AnyTimes()
			},
			userId:      userId,
			realm:       realm,
			wantCode:    200,
			wantMessage: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			keycloak := servicemocks.NewMockIKeycloak(c)
			tt.mockBehavior(keycloak)

			w := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(w)

			engine.GET("/test", authorize(keycloak, tt.roles), func(ctx *gin.Context) {
				ctx.Status(200)
				ctx.Header(modules.UserID, ctx.Value(modules.UserID).(string))
				ctx.Header(modules.ClientID, ctx.Value(modules.ClientID).(string))
				return
			})

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			req.Header.Set("Authorization", "Bearer test")

			engine.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)
			assert.Equal(t, tt.wantMessage, "")
			assert.Equal(t, tt.userId, w.Header().Get(modules.UserID))
			assert.Equal(t, tt.realm, w.Header().Get(modules.ClientID))
		})
	}
}
