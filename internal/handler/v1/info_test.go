package v1

import (
	"encoding/json"
	"github.com/Nerzal/gocloak/v8"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/a5805/ondeu/ondeu-back/internal/service"
	servicemocks "gitlab.com/a5805/ondeu/ondeu-back/internal/service/mocks"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_getRoles(t *testing.T) {
	type mockBehavior func(*servicemocks.MockInformationService)

	str := `[ { "id": "89c42dfb-d026-46ed-8e5d-0abea9639f80", "name": "student", "composite": false, "clientRole": true, "containerId": "6f150b81-33b4-47a2-a10a-1fb655cc0cab" }, { "id": "0636d9d3-d2ad-4b47-8cce-b7a8a8e1a6eb", "name": "manager", "composite": false, "clientRole": true, "containerId": "6f150b81-33b4-47a2-a10a-1fb655cc0cab" }, { "id": "70d3b502-538e-4ee0-98b8-d67748675a86", "name": "admin", "composite": false, "clientRole": true, "containerId": "6f150b81-33b4-47a2-a10a-1fb655cc0cab" } ]`

	var roles []*gocloak.Role
	err := json.Unmarshal([]byte(str), &roles)
	require.NoError(t, err)

	tests := []struct {
		name                 string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Success.",
			mockBehavior: func(r *servicemocks.MockInformationService) {
				r.EXPECT().
					GetRoles(gomock.Any()).
					Return(roles, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `[ { "id": "89c42dfb-d026-46ed-8e5d-0abea9639f80", "name": "student", "composite": false, "clientRole": true, "containerId": "6f150b81-33b4-47a2-a10a-1fb655cc0cab" }, { "id": "0636d9d3-d2ad-4b47-8cce-b7a8a8e1a6eb", "name": "manager", "composite": false, "clientRole": true, "containerId": "6f150b81-33b4-47a2-a10a-1fb655cc0cab" }, { "id": "70d3b502-538e-4ee0-98b8-d67748675a86", "name": "admin", "composite": false, "clientRole": true, "containerId": "6f150b81-33b4-47a2-a10a-1fb655cc0cab" } ]`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := servicemocks.NewMockInformationService(c)
			tt.mockBehavior(repo)

			services := &service.Services{InformationService: repo}
			handler := Handler{services, nil, nil}

			// Init Endpoint
			r := gin.New()
			r.GET("/api/v1/info/roles", handler.getRoles)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(
				http.MethodGet,
				"/api/v1/info/roles",
				nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t,
				strings.ReplaceAll(tt.expectedResponseBody, " ", ""),
				strings.ReplaceAll(w.Body.String(), " ", ""),
			)
		})
	}
}
