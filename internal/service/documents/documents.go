package documents

import (
	"context"
	"fmt"
	"gitlab.com/a5805/ondeu/ondeu-back/internal/remote"
	"gitlab.com/a5805/ondeu/ondeu-back/internal/repository"
	"gitlab.com/a5805/ondeu/ondeu-back/pkg/modules"
	"gitlab.com/a5805/ondeu/ondeu-back/pkg/modules/dto"
	"mime"
	"mime/multipart"
	"path/filepath"
	"time"
)

type Service struct {
	repos   repository.DocumentRepository
	remotes remote.DocumentsRemote
}

func NewService(repos repository.DocumentRepository, remotes remote.DocumentsRemote) *Service {
	return &Service{
		repos:   repos,
		remotes: remotes,
	}
}

func (s *Service) Create(ctx context.Context, document dto.Document, file *multipart.FileHeader) (dto.Document, error) {
	userID, ok := ctx.Value(modules.UserID).(string)
	if !ok {
		return document, fmt.Errorf("unauthorized action is prohibited")
	}

	content, err := file.Open()
	if err != nil {
		return document, err
	}
	defer content.Close()

	if document.Name == "" {
		document.Name = file.Filename
	}

	document.UserID = userID
	document.RequestContent = content
	document.Size = file.Size
	document.Type, _, _ = mime.ParseMediaType(file.Header.Get("Content-Type"))
	document.Extension = filepath.Ext(file.Filename)

	stored, err := s.repos.Create(ctx, document)
	if err != nil {
		return document, err
	}

	uploaded, err := s.remotes.Upload(ctx, stored)
	if err != nil {
		return uploaded, err
	}

	return uploaded, nil
}

func (s *Service) Get(ctx context.Context, doc dto.Document, download bool) (dto.Document, error) {
	userId, ok := ctx.Value(modules.UserID).(string)
	if !ok {
		return doc, fmt.Errorf("unauthorized action is prohibited")
	}

	doc.UserID = userId

	stored, err := s.repos.Get(ctx, doc)
	if err != nil {
		return stored, err
	}

	if !download {
		return stored, nil
	}

	return s.remotes.Get(ctx, stored)
}

func (s *Service) Delete(ctx context.Context, doc dto.Document) (dto.Document, error) {
	userId, ok := ctx.Value(modules.UserID).(string)
	if !ok {
		return doc, fmt.Errorf("unauthorized action is prohibited")
	}

	doc.UserID = userId

	document, err := s.repos.Get(ctx, doc)
	if err != nil {
		return document, err
	}

	document, err = s.remotes.Delete(ctx, document)
	if err != nil {
		return document, err
	}

	return s.repos.Delete(ctx, doc)
}

func (s *Service) Update(ctx context.Context, doc dto.Document) (dto.Document, error) {
	userId, ok := ctx.Value(modules.UserID).(string)
	if !ok {
		return doc, fmt.Errorf("unauthorized action is prohibited")
	}

	doc.UserID = userId

	return s.repos.Update(ctx, doc)
}

func (s *Service) ListByTree(ctx context.Context, ids []uint) ([]dto.Document, error) {
	return s.repos.ListByTree(ctx, ids)
}

func (s *Service) ListByGroups(ctx context.Context, groupIds []uint) ([]dto.Document, error) {
	return s.repos.ListByGroups(ctx, groupIds)
}

func (s *Service) Share(ctx context.Context, doc dto.Document, duration time.Duration) (dto.Document, error) {
	userId, ok := ctx.Value(modules.UserID).(string)
	if !ok {
		return doc, fmt.Errorf("unauthorized action is prohibited")
	}

	doc.UserID = userId

	document, err := s.repos.Get(ctx, doc)
	if err != nil {
		return document, err
	}

	return s.remotes.Share(ctx, document, duration)
}
