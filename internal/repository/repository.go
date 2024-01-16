package repository

import (
	"context"
	"gitlab.com/a5805/ondeu/ondeu-back/internal/repository/documents"
	"gitlab.com/a5805/ondeu/ondeu-back/internal/repository/tree"
	"gitlab.com/a5805/ondeu/ondeu-back/pkg/modules/dto"
	"gorm.io/gorm"
)

type DocumentRepository interface {
	// Create creates a new document
	Create(ctx context.Context, doc dto.Document) (dto.Document, error)
	// Get returns a document
	Get(ctx context.Context, doc dto.Document) (dto.Document, error)
	// Delete deletes a document
	Delete(ctx context.Context, doc dto.Document) (dto.Document, error)
	// Update updates a document
	Update(ctx context.Context, doc dto.Document) (dto.Document, error)
	// FindByCondition Filter documents by condition
	FindByCondition(ctx context.Context, field, param string) ([]dto.Document, error)
	// ListByTree returns a slice of documents by tree id
	ListByTree(ctx context.Context, ids []uint) ([]dto.Document, error)
	// ListByGroups returns a slice of documents by group id
	ListByGroups(ctx context.Context, ids []uint) ([]dto.Document, error)
}

type TreeRepository interface {
	// Create creates a new tree
	Create(ctx context.Context, tree dto.Tree) (dto.Tree, error)
	// Get returns a tree
	Get(ctx context.Context, tree dto.Tree) (dto.Tree, error)
	// List returns a tree
	List(ctx context.Context, tree dto.Tree) ([]dto.Tree, error)
	// Update deletes a tree
	Update(ctx context.Context, tree dto.Tree) (dto.Tree, error)
	// Delete deletes a tree
	Delete(ctx context.Context, tree dto.Tree) (dto.Tree, error)
}

type Repository struct {
	DocumentRepository
	TreeRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		DocumentRepository: documents.NewRepository(db),
		TreeRepository:     tree.NewRepository(db),
	}
}
