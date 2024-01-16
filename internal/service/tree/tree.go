package tree

import (
	"context"
	"fmt"
	"gitlab.com/a5805/ondeu/ondeu-back/internal/repository"
	"gitlab.com/a5805/ondeu/ondeu-back/pkg/modules"
	"gitlab.com/a5805/ondeu/ondeu-back/pkg/modules/dto"
)

type Service struct {
	repos repository.TreeRepository
}

func NewService(repos repository.TreeRepository) *Service {
	return &Service{
		repos: repos,
	}
}

func (s *Service) Create(ctx context.Context, in dto.Tree) (dto.Tree, error) {
	userId, ok := ctx.Value(modules.UserID).(string)
	if !ok {
		return in, fmt.Errorf("unauthorized action is prohibited")
	}
	in.UserID = userId

	return s.repos.Create(ctx, in)
}

func (s *Service) Get(ctx context.Context, tree dto.Tree) (dto.Tree, error) {
	//userId, ok := ctx.Value(modules.UserID).(string)
	//if !ok {
	//	return nil, fmt.Errorf("unauthorized action is prohibited")
	//}

	tree, err := s.repos.Get(ctx, tree)
	if err != nil {
		return dto.Tree{}, err
	}
	return tree, nil
}

func (s *Service) List(ctx context.Context, tree dto.Tree) ([]dto.Tree, error) {
	userId, ok := ctx.Value(modules.UserID).(string)
	if !ok {
		return nil, fmt.Errorf("unauthorized action is prohibited")
	}

	tree.UserID = userId

	trees, err := s.repos.List(ctx, tree)
	if err != nil {
		return nil, err
	}
	return trees, nil
}

func (s *Service) Delete(ctx context.Context, doc dto.Tree) (dto.Tree, error) {
	userId, ok := ctx.Value(modules.UserID).(string)
	if !ok {
		return doc, fmt.Errorf("unauthorized action is prohibited")
	}

	doc.UserID = userId

	return s.repos.Delete(ctx, doc)
}

func (s *Service) Update(ctx context.Context, doc dto.Tree) (dto.Tree, error) {
	userId, ok := ctx.Value(modules.UserID).(string)
	if !ok {
		return doc, fmt.Errorf("unauthorized action is prohibited")
	}

	doc.UserID = userId

	return s.repos.Update(ctx, doc)
}

func (s *Service) GetTreeIDs(ctx context.Context, trees []dto.Tree) []uint {
	var treeIds []uint
	for _, tree := range trees {
		treeIds = append(treeIds, tree.ID)
	}
	return treeIds
}

func (s *Service) FormTree(ctx context.Context, trees []dto.Tree, docs []dto.Document) []dto.Tree {
	for _, doc := range docs {
		for i, tree := range trees {
			if doc.TreeID == tree.ID {
				trees[i].Documents = append(trees[i].Documents, doc)
			}
		}
	}
	return trees
}
