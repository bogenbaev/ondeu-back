package service

import (
	"context"
	"github.com/Nerzal/gocloak/v8"
	"gitlab.com/a5805/ondeu/ondeu-back/internal/remote"
	"gitlab.com/a5805/ondeu/ondeu-back/internal/repository"
	"gitlab.com/a5805/ondeu/ondeu-back/internal/service/documents"
	"gitlab.com/a5805/ondeu/ondeu-back/internal/service/information"
	"gitlab.com/a5805/ondeu/ondeu-back/internal/service/tree"
	keycloak2 "gitlab.com/a5805/ondeu/ondeu-back/pkg/gocloak"
	"gitlab.com/a5805/ondeu/ondeu-back/pkg/modules"
	"gitlab.com/a5805/ondeu/ondeu-back/pkg/modules/dto"
	"mime/multipart"
	"time"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type DocumentService interface {
	// Create creates a new document
	Create(ctx context.Context, in dto.Document, file *multipart.FileHeader) (dto.Document, error)
	// Get returns a document
	Get(ctx context.Context, doc dto.Document, download bool) (dto.Document, error)
	// Update updates a document
	Update(ctx context.Context, doc dto.Document) (dto.Document, error)
	// Delete deletes a document
	Delete(ctx context.Context, doc dto.Document) (dto.Document, error)

	// Share create share link for a document
	Share(ctx context.Context, doc dto.Document, duration time.Duration) (dto.Document, error)
	// ListByTree returns a slice of documents by tree id
	ListByTree(ctx context.Context, ids []uint) ([]dto.Document, error)
	// ListByGroups returns a slice of documents by group id
	ListByGroups(ctx context.Context, ids []uint) ([]dto.Document, error)
}

type TreeService interface {
	// Create creates a new tree
	Create(ctx context.Context, tree dto.Tree) (dto.Tree, error)
	// Get returns a specific tree
	Get(ctx context.Context, tree dto.Tree) (dto.Tree, error)
	// List returns all tree
	List(ctx context.Context, tree dto.Tree) ([]dto.Tree, error)
	// Update deletes a tree
	Update(ctx context.Context, tree dto.Tree) (dto.Tree, error)
	// Delete deletes a tree
	Delete(ctx context.Context, tree dto.Tree) (dto.Tree, error)

	// GetTreeIDs returns a slice of tree ids
	GetTreeIDs(ctx context.Context, trees []dto.Tree) []uint

	// FormTree returns a slice of trees with documents
	FormTree(ctx context.Context, trees []dto.Tree, docs []dto.Document) []dto.Tree
}

type InformationService interface {
	// GetRoles returns a slice of users
	GetRoles(ctx context.Context) ([]*gocloak.Role, error)
}

type Services struct {
	TreeService
	DocumentService
	InformationService
}

func NewServices(cfg *modules.AppConfigs, keycloak keycloak2.IKeycloak, repos *repository.Repository, remotes *remote.Remote) *Services {
	return &Services{
		TreeService:        tree.NewService(repos.TreeRepository),
		DocumentService:    documents.NewService(repos.DocumentRepository, remotes),
		InformationService: information.NewService(cfg.Keycloak, keycloak),
	}
}
