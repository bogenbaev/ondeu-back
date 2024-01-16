package v1

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/a5805/ondeu/ondeu-back/internal/service"
	servicemocks "gitlab.com/a5805/ondeu/ondeu-back/internal/service/mocks"
	"gitlab.com/a5805/ondeu/ondeu-back/pkg/modules/dto"
	"gorm.io/gorm"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestHandler_createDocument(t *testing.T) {
	type mockBehavior func(r *servicemocks.MockDocumentService, document dto.Document)

	createdData := time.Now()
	userID := uuid.New().String()
	createdFileUUID := uuid.New()

	tests := []struct {
		name                 string
		fileExists           bool
		inputDocument        dto.Document
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:          "Failed. Service. File Invalid",
			fileExists:    true,
			inputDocument: dto.Document{},
			mockBehavior: func(r *servicemocks.MockDocumentService, document dto.Document) {
				r.EXPECT().
					Create(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(dto.Document{}, os.ErrInvalid)
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"reason":"invalid argument"}`,
		},
		{
			name:          "Failed. Database. Duplicate Key",
			fileExists:    true,
			inputDocument: dto.Document{},
			mockBehavior: func(r *servicemocks.MockDocumentService, document dto.Document) {
				r.EXPECT().
					Create(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(dto.Document{}, gorm.ErrDuplicatedKey)
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"reason":"duplicated key not allowed"}`,
		},
		{
			name:          "Failed. Database. Invalid Value",
			fileExists:    true,
			inputDocument: dto.Document{},
			mockBehavior: func(r *servicemocks.MockDocumentService, document dto.Document) {
				r.EXPECT().
					Create(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(dto.Document{}, gorm.ErrInvalidValue)
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"reason":"invalid value, should be pointer to struct or slice"}`,
		},
		{
			name:       "Success.",
			fileExists: true,
			inputDocument: dto.Document{
				ID:        123,
				UserID:    userID,
				TreeID:    1,
				CreatedAt: createdData,
				UpdatedAt: createdData,
				Name:      "test",
				Path:      createdFileUUID,
			},
			mockBehavior: func(r *servicemocks.MockDocumentService, document dto.Document) {
				r.EXPECT().
					Create(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(dto.Document{
						ID:        123,
						UserID:    userID,
						TreeID:    1,
						CreatedAt: createdData,
						UpdatedAt: createdData,
						Name:      "test",
						Path:      createdFileUUID,
					}, nil)
			},
			expectedStatusCode: 200,
			expectedResponseBody: fmt.Sprintf(`{"id":123,"createdAt":"%s","updatedAt":"%s","name":"test","path":"%s","template":null}`,
				createdData.Format(time.RFC3339Nano),
				createdData.Format(time.RFC3339Nano), createdFileUUID.String()),
		},
	}
	for _, tt := range tests {
		body := new(bytes.Buffer)
		m := multipart.NewWriter(body)

		if tt.fileExists {
			writer, err := m.CreateFormFile("file", "./documents.go")
			require.NoError(t, err)

			file, err := os.Open("./documents.go")
			require.NoError(t, err)

			_, err = io.Copy(writer, file)
			require.NoError(t, err)

			m.WriteField("name", "test")

			m.Close()
		}

		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := servicemocks.NewMockDocumentService(c)
			tt.mockBehavior(repo, tt.inputDocument)

			services := &service.Services{DocumentService: repo}
			handler := Handler{services, nil, nil}

			// Init Endpoint
			r := gin.New()
			r.POST("/api/v1/tree/:treeID/document", handler.createDocument)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", fmt.Sprintf("/api/v1/tree/%d/document", 1), body)
			req.Header.Add("Content-Type", m.FormDataContentType())

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, tt.expectedStatusCode)
			assert.Equal(t, w.Body.String(), tt.expectedResponseBody)
		})
	}
}

func TestHandler_updateDocument(t *testing.T) {
	type mockBehavior func(r *servicemocks.MockDocumentService, document dto.Document)

	updated := time.Now()
	created := updated.Add(-time.Hour)
	isTemplate := true
	userID := uuid.New().String()
	createdFileUUID := uuid.New()

	tests := []struct {
		name                 string
		inputBody            string
		inputDocument        dto.Document
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Failed. Database. Duplicate Key",
			inputBody: `{"name":"This is a template document","template":true}`,
			inputDocument: dto.Document{
				ID:     123,
				TreeID: 1,
				Name:   "This is a template document",
			},
			mockBehavior: func(r *servicemocks.MockDocumentService, document dto.Document) {
				r.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(dto.Document{
						ID:     123,
						TreeID: 1,
						Name:   "This is a template document",
					}, gorm.ErrDuplicatedKey)
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"reason":"duplicated key not allowed"}`,
		},
		{
			name:      "Failed. Database. Invalid Value",
			inputBody: `{"name":"This is a template document","template":true}`,
			inputDocument: dto.Document{
				ID:     123,
				TreeID: 1,
				Name:   "This is a template document",
			},
			mockBehavior: func(r *servicemocks.MockDocumentService, document dto.Document) {
				r.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(dto.Document{
						ID:     123,
						TreeID: 1,
						Name:   "This is a template document",
					}, gorm.ErrInvalidValue)
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"reason":"invalid value, should be pointer to struct or slice"}`,
		},
		{
			name:      "Success. Update document",
			inputBody: `{"name":"This is a template document","template":true}`,
			inputDocument: dto.Document{
				ID:        123,
				UserID:    userID,
				TreeID:    1,
				UpdatedAt: updated,
				Name:      "This is a template document",
				Path:      createdFileUUID,
				Template:  &isTemplate,
			},
			mockBehavior: func(r *servicemocks.MockDocumentService, document dto.Document) {
				r.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(dto.Document{
						ID:        123,
						UserID:    userID,
						TreeID:    1,
						CreatedAt: created,
						UpdatedAt: updated,
						Name:      "This is a template document",
						Path:      createdFileUUID,
						Template:  &isTemplate,
					}, nil)
			},
			expectedStatusCode: 200,
			expectedResponseBody: fmt.Sprintf(`{"id":123,"createdAt":"%s","updatedAt":"%s","name":"This is a template document","path":"%s","template":true}`,
				created.Format(time.RFC3339Nano),
				updated.Format(time.RFC3339Nano),
				createdFileUUID.String()),
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := servicemocks.NewMockDocumentService(c)
			tt.mockBehavior(repo, tt.inputDocument)

			services := &service.Services{DocumentService: repo}
			handler := Handler{services, nil, nil}

			// Init Endpoint
			r := gin.New()
			r.PUT("/api/v1/tree/:treeID/document/:docID", handler.updateDocument)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/tree/%d/document/%d", 1, 1), strings.NewReader(tt.inputBody))
			req.Header.Add("Content-Type", "application/json")

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_readDocument(t *testing.T) {
	type mockBehavior func(r *servicemocks.MockDocumentService, document dto.Document)

	updated := time.Now()
	created := updated.Add(-time.Hour)
	isTemplate := true
	userID := uuid.New().String()
	createdFileUUID := uuid.New()

	tests := []struct {
		name                 string
		inputDocument        dto.Document
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Failed. Database. Duplicate Key",
			inputDocument: dto.Document{
				ID:     123,
				TreeID: 1,
				Name:   "This is a template document",
			},
			mockBehavior: func(r *servicemocks.MockDocumentService, document dto.Document) {
				r.EXPECT().
					Get(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(dto.Document{
						ID:     123,
						TreeID: 1,
						Name:   "This is a template document",
					}, gorm.ErrDuplicatedKey)
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"reason":"duplicated key not allowed"}`,
		},
		{
			name: "Failed. Database. Invalid Value",
			inputDocument: dto.Document{
				ID:     123,
				TreeID: 1,
				Name:   "This is a template document",
			},
			mockBehavior: func(r *servicemocks.MockDocumentService, document dto.Document) {
				r.EXPECT().
					Get(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(dto.Document{
						ID:     123,
						TreeID: 1,
						Name:   "This is a template document",
					}, gorm.ErrInvalidValue)
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"reason":"invalid value, should be pointer to struct or slice"}`,
		},
		{
			name: "Success. Update document",
			inputDocument: dto.Document{
				ID:        123,
				UserID:    userID,
				TreeID:    1,
				UpdatedAt: updated,
				Name:      "This is a template document",
				Path:      createdFileUUID,
				Template:  &isTemplate,
			},
			mockBehavior: func(r *servicemocks.MockDocumentService, document dto.Document) {
				r.EXPECT().
					Get(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(dto.Document{
						ID:        123,
						UserID:    userID,
						TreeID:    1,
						CreatedAt: created,
						UpdatedAt: updated,
						Name:      "This is a template document",
						Path:      createdFileUUID,
						Template:  &isTemplate,
					}, nil)
			},
			expectedStatusCode: 200,
			expectedResponseBody: fmt.Sprintf(`{"id":123,"createdAt":"%s","updatedAt":"%s","name":"This is a template document","path":"%s","template":true}`,
				created.Format(time.RFC3339Nano),
				updated.Format(time.RFC3339Nano),
				createdFileUUID.String()),
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := servicemocks.NewMockDocumentService(c)
			tt.mockBehavior(repo, tt.inputDocument)

			services := &service.Services{DocumentService: repo}
			handler := Handler{services, nil, nil}

			// Init Endpoint
			r := gin.New()
			r.GET("/api/v1/tree/:treeID/document/:docID", handler.readDocument)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/tree/%d/document/%d", 1, 1), nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_deleteDocument(t *testing.T) {
	type mockBehavior func(r *servicemocks.MockDocumentService, document dto.Document)

	updated := time.Now()
	created := updated.Add(-time.Hour)
	isTemplate := true
	userID := uuid.New().String()
	createdFileUUID := uuid.New()

	tests := []struct {
		name                 string
		inputDocument        dto.Document
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Failed. Database. Duplicate Key",
			inputDocument: dto.Document{
				ID:     123,
				TreeID: 1,
				Name:   "This is a template document",
			},
			mockBehavior: func(r *servicemocks.MockDocumentService, document dto.Document) {
				r.EXPECT().
					Delete(gomock.Any(), gomock.Any()).
					Return(dto.Document{
						ID:     123,
						TreeID: 1,
						Name:   "This is a template document",
					}, gorm.ErrDuplicatedKey)
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"reason":"duplicated key not allowed"}`,
		},
		{
			name: "Failed. Database. Invalid Value",
			inputDocument: dto.Document{
				ID:     123,
				TreeID: 1,
				Name:   "This is a template document",
			},
			mockBehavior: func(r *servicemocks.MockDocumentService, document dto.Document) {
				r.EXPECT().
					Delete(gomock.Any(), gomock.Any()).
					Return(dto.Document{
						ID:     123,
						TreeID: 1,
						Name:   "This is a template document",
					}, gorm.ErrInvalidValue)
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"reason":"invalid value, should be pointer to struct or slice"}`,
		},
		{
			name: "Success. Update document",
			inputDocument: dto.Document{
				ID:        123,
				UserID:    userID,
				TreeID:    1,
				UpdatedAt: updated,
				Name:      "This is a template document",
				Path:      createdFileUUID,
				Template:  &isTemplate,
			},
			mockBehavior: func(r *servicemocks.MockDocumentService, document dto.Document) {
				r.EXPECT().
					Delete(gomock.Any(), gomock.Any()).
					Return(dto.Document{
						ID:        123,
						UserID:    userID,
						TreeID:    1,
						CreatedAt: created,
						UpdatedAt: updated,
						Name:      "This is a template document",
						Path:      createdFileUUID,
						Template:  &isTemplate,
					}, nil)
			},
			expectedStatusCode: 200,
			expectedResponseBody: fmt.Sprintf(`{"id":123,"createdAt":"%s","updatedAt":"%s","name":"This is a template document","path":"%s","template":true}`,
				created.Format(time.RFC3339Nano),
				updated.Format(time.RFC3339Nano),
				createdFileUUID.String()),
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := servicemocks.NewMockDocumentService(c)
			tt.mockBehavior(repo, tt.inputDocument)

			services := &service.Services{DocumentService: repo}
			handler := Handler{services, nil, nil}

			// Init Endpoint
			r := gin.New()
			r.DELETE("/api/v1/tree/:treeID/document/:docID", handler.deleteDocument)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/tree/%d/document/%d", 1, 1), nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_shareDocument(t *testing.T) {
	type mockBehavior func(r *servicemocks.MockDocumentService, document dto.Document)

	updated := time.Now()
	created := updated.Add(-time.Hour)
	isTemplate := true
	userID := uuid.New().String()
	createdFileUUID := uuid.New()

	tests := []struct {
		name                 string
		inputDocument        dto.Document
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Failed. Database. Duplicate Key",
			inputDocument: dto.Document{
				ID:     123,
				TreeID: 1,
				Name:   "This is a template document",
			},
			mockBehavior: func(r *servicemocks.MockDocumentService, document dto.Document) {
				r.EXPECT().
					Share(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(dto.Document{
						ID:     123,
						TreeID: 1,
						Name:   "This is a template document",
					}, gorm.ErrDuplicatedKey)
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"reason":"duplicated key not allowed"}`,
		},
		{
			name: "Failed. Database. Invalid Value",
			inputDocument: dto.Document{
				ID:     123,
				TreeID: 1,
				Name:   "This is a template document",
			},
			mockBehavior: func(r *servicemocks.MockDocumentService, document dto.Document) {
				r.EXPECT().
					Share(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(dto.Document{
						ID:     123,
						TreeID: 1,
						Name:   "This is a template document",
					}, gorm.ErrInvalidValue)
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"reason":"invalid value, should be pointer to struct or slice"}`,
		},
		{
			name: "Success. Update document",
			inputDocument: dto.Document{
				ID:        123,
				UserID:    userID,
				TreeID:    1,
				UpdatedAt: updated,
				Name:      "This is a template document",
				Path:      createdFileUUID,
				Template:  &isTemplate,
			},
			mockBehavior: func(r *servicemocks.MockDocumentService, document dto.Document) {
				r.EXPECT().
					Share(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(dto.Document{
						ID:        123,
						UserID:    userID,
						TreeID:    1,
						CreatedAt: created,
						UpdatedAt: updated,
						Name:      "This is a template document",
						Path:      createdFileUUID,
						Template:  &isTemplate,
					}, nil)
			},
			expectedStatusCode: 200,
			expectedResponseBody: fmt.Sprintf(`{"id":123,"createdAt":"%s","updatedAt":"%s","name":"This is a template document","path":"%s","template":true}`,
				created.Format(time.RFC3339Nano),
				updated.Format(time.RFC3339Nano),
				createdFileUUID.String()),
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := servicemocks.NewMockDocumentService(c)
			tt.mockBehavior(repo, tt.inputDocument)

			services := &service.Services{DocumentService: repo}
			handler := Handler{services, nil, nil}

			// Init Endpoint
			r := gin.New()
			r.GET("/api/v1/tree/:treeID/document/:docID/share", handler.shareDocument)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/tree/%d/document/%d/share", 1, 1), nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}
