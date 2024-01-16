package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gitlab.com/a5805/ondeu/ondeu-back/internal/service"
	servicemocks "gitlab.com/a5805/ondeu/ondeu-back/internal/service/mocks"
	"gitlab.com/a5805/ondeu/ondeu-back/pkg/modules/dto"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestHandler_createTree(t *testing.T) {
	type mockBehavior func(*servicemocks.MockTreeService, dto.Tree)

	createdData := time.Now()
	userID := uuid.New().String()

	true_ := true

	tests := []struct {
		name                 string
		raw                  string
		input                dto.Tree
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:  "Failed. Database. Duplicate Key",
			raw:   `{"parentId":0,"name":"1 grade","role":"bachelor","template":true,"group":true}`,
			input: dto.Tree{},
			mockBehavior: func(r *servicemocks.MockTreeService, tree dto.Tree) {
				r.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(dto.Tree{}, gorm.ErrDuplicatedKey)
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"reason":"duplicated key not allowed"}`,
		},
		{
			name:  "Failed. Database. Invalid Value",
			raw:   `{"parentId":0,"name":"1 grade","role":"bachelor","template":true,"group":true}`,
			input: dto.Tree{},
			mockBehavior: func(r *servicemocks.MockTreeService, tree dto.Tree) {
				r.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(dto.Tree{}, gorm.ErrInvalidValue)
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"reason":"invalid value, should be pointer to struct or slice"}`,
		},
		{
			name: "Success.",
			raw:  `{"parentId":0,"name":"1 grade","role":"bachelor","template":true,"group":true}`,
			input: dto.Tree{
				UserID:    userID,
				CreatedAt: createdData,
				UpdatedAt: createdData,
				Name:      "1 grade",
				Role:      "bachelor",
				Template:  &true_,
				Group:     &true_,
			},
			mockBehavior: func(r *servicemocks.MockTreeService, tree dto.Tree) {
				r.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(dto.Tree{
						ID:        12,
						UserID:    userID,
						CreatedAt: createdData,
						UpdatedAt: createdData,
						Name:      "1 grade",
						Role:      "bachelor",
						Template:  &true_,
						Group:     &true_,
					}, nil)
			},
			expectedStatusCode: 200,
			expectedResponseBody: fmt.Sprintf(
				`{"id":12,"parentID":0,"createdAt":"%s","updatedAt":"%s","name":"1 grade","role":"bachelor","template":true,"group":true,"documents":null}`,
				createdData.Format(time.RFC3339Nano),
				createdData.Format(time.RFC3339Nano),
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := servicemocks.NewMockTreeService(c)
			tt.mockBehavior(repo, tt.input)

			services := &service.Services{TreeService: repo}
			handler := Handler{services, nil, nil}

			// Init Endpoint
			r := gin.New()
			r.POST("/api/v1/tree/:treeID", handler.createTree)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(
				http.MethodPost,
				fmt.Sprintf("/api/v1/tree/%d", 1),
				strings.NewReader(tt.raw))
			req.Header.Set("Content-Type", "application/json")

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_updateTree(t *testing.T) {
	type mockBehavior func(*servicemocks.MockTreeService, dto.Tree)

	createdData := time.Now()
	userID := uuid.New().String()

	false_ := false

	tests := []struct {
		name                 string
		raw                  string
		input                dto.Tree
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:  "Failed. Database. Duplicate Key",
			raw:   `{"parentId":0,"name":"1 grade","role":"bachelor","template":true,"group":true}`,
			input: dto.Tree{},
			mockBehavior: func(r *servicemocks.MockTreeService, tree dto.Tree) {
				r.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(dto.Tree{}, gorm.ErrDuplicatedKey)
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"reason":"duplicated key not allowed"}`,
		},
		{
			name:  "Failed. Database. Invalid Value",
			raw:   `{"parentId":0,"name":"1 grade","role":"bachelor","template":false,"group":false}`,
			input: dto.Tree{},
			mockBehavior: func(r *servicemocks.MockTreeService, tree dto.Tree) {
				r.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(dto.Tree{}, gorm.ErrInvalidValue)
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"reason":"invalid value, should be pointer to struct or slice"}`,
		},
		{
			name: "Success.",
			raw:  `{"parentId":0,"name":"1 grade","role":"bachelor","template":false,"group":false}`,
			input: dto.Tree{
				UserID:    userID,
				CreatedAt: createdData,
				UpdatedAt: createdData,
				Name:      "1 grade",
				Role:      "bachelor",
				Template:  &false_,
				Group:     &false_,
			},
			mockBehavior: func(r *servicemocks.MockTreeService, tree dto.Tree) {
				r.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(dto.Tree{
						ID:        12,
						UserID:    userID,
						CreatedAt: createdData,
						UpdatedAt: createdData,
						Name:      "1 grade",
						Role:      "bachelor",
						Template:  &false_,
						Group:     &false_,
					}, nil)
			},
			expectedStatusCode: 200,
			expectedResponseBody: fmt.Sprintf(
				`{"id":12,"parentID":0,"createdAt":"%s","updatedAt":"%s","name":"1 grade","role":"bachelor","template":false,"group":false,"documents":null}`,
				createdData.Format(time.RFC3339Nano),
				createdData.Format(time.RFC3339Nano),
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := servicemocks.NewMockTreeService(c)
			tt.mockBehavior(repo, tt.input)

			services := &service.Services{TreeService: repo}
			handler := Handler{services, nil, nil}

			// Init Endpoint
			r := gin.New()
			r.PUT("/api/v1/tree/:treeID", handler.updateTree)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(
				http.MethodPut,
				fmt.Sprintf("/api/v1/tree/%d", 1),
				strings.NewReader(tt.raw))
			req.Header.Set("Content-Type", "application/json")

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_deleteTree(t *testing.T) {
	type mockBehavior func(*servicemocks.MockTreeService, dto.Tree)

	createdData := time.Now()
	userID := uuid.New().String()

	false_ := false

	tests := []struct {
		name                 string
		input                dto.Tree
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Failed. Database. Record Not Found",
			input: dto.Tree{
				ID: 12,
			},
			mockBehavior: func(r *servicemocks.MockTreeService, tree dto.Tree) {
				r.EXPECT().
					Delete(gomock.Any(), gomock.Any()).
					Return(dto.Tree{}, gorm.ErrRecordNotFound)
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"reason":"record not found"}`,
		},
		{
			name: "Failed. Database. Invalid Value",
			input: dto.Tree{
				ID: 12,
			},
			mockBehavior: func(r *servicemocks.MockTreeService, tree dto.Tree) {
				r.EXPECT().
					Delete(gomock.Any(), gomock.Any()).
					Return(dto.Tree{}, gorm.ErrInvalidValue)
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"reason":"invalid value, should be pointer to struct or slice"}`,
		},
		{
			name: "Success.",
			input: dto.Tree{
				ID: 12,
			},
			mockBehavior: func(r *servicemocks.MockTreeService, tree dto.Tree) {
				r.EXPECT().
					Delete(gomock.Any(), gomock.Any()).
					Return(dto.Tree{
						ID:        12,
						UserID:    userID,
						CreatedAt: createdData,
						UpdatedAt: createdData,
						Name:      "1 grade",
						Role:      "bachelor",
						Template:  &false_,
						Group:     &false_,
					}, nil)
			},
			expectedStatusCode: 200,
			expectedResponseBody: fmt.Sprintf(
				`{"id":12,"parentID":0,"createdAt":"%s","updatedAt":"%s","name":"1 grade","role":"bachelor","template":false,"group":false,"documents":null}`,
				createdData.Format(time.RFC3339Nano),
				createdData.Format(time.RFC3339Nano),
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := servicemocks.NewMockTreeService(c)
			tt.mockBehavior(repo, tt.input)

			services := &service.Services{TreeService: repo}
			handler := Handler{services, nil, nil}

			// Init Endpoint
			r := gin.New()
			r.DELETE("/api/v1/tree/:treeID", handler.deleteTree)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(
				http.MethodDelete,
				fmt.Sprintf("/api/v1/tree/%d", 1),
				nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}
