package remote

import (
	"context"
	"github.com/aws/aws-sdk-go/service/s3"
	"gitlab.com/a5805/ondeu/ondeu-back/internal/remote/documents"
	"gitlab.com/a5805/ondeu/ondeu-back/pkg/modules"
	"gitlab.com/a5805/ondeu/ondeu-back/pkg/modules/dto"
	"time"
)

type DocumentsRemote interface {
	// Upload uploads a document to spaces
	Upload(ctx context.Context, doc dto.Document) (dto.Document, error)
	// Get returns a document from spaces
	Get(ctx context.Context, doc dto.Document) (dto.Document, error)
	// Delete deletes a document from spaces
	Delete(ctx context.Context, doc dto.Document) (dto.Document, error)
	// Share create share link for a document
	Share(ctx context.Context, doc dto.Document, duration time.Duration) (dto.Document, error)
}

type Remote struct {
	DocumentsRemote
}

func NewRemote(s3 *s3.S3, cfg *modules.ObjectStorage) *Remote {
	return &Remote{
		DocumentsRemote: documents.NewRemote(s3, cfg),
	}
}
