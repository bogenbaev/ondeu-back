package tree

import (
	"context"
	"github.com/sirupsen/logrus"
	"gitlab.com/a5805/ondeu/ondeu-back/pkg/modules/dto"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (fm *Repository) Create(ctx context.Context, tree dto.Tree) (dto.Tree, error) {
	return tree, fm.db.WithContext(ctx).Create(&tree).Error
}

func (fm *Repository) Get(ctx context.Context, tree dto.Tree) (dto.Tree, error) {
	return tree, fm.db.WithContext(ctx).Find(&tree).Error
}

func (fm *Repository) List(ctx context.Context, tree dto.Tree) ([]dto.Tree, error) {
	sql := `WITH RECURSIVE cte AS (
		SELECT t1.id, t1.parent_id, t1.name,
			   t1.created_at, t1.updated_at, t1.role, t1.template, t1.group
		FROM   trees t1
		WHERE  t1.parent_id = ? and t1.user_id = ?
	
		UNION  ALL
		SELECT t2.id, t2.parent_id, t2.name,
			   t2.created_at, t2.updated_at, t2.role, t2.template, t2.group
		FROM trees t2 JOIN cte c ON t2.parent_id = c.id and t2.user_id = ?
	) SELECT * from cte;`

	var trees []dto.Tree
	if err := fm.db.WithContext(ctx).
		Model(dto.Tree{}).
		Preload("Documents").
		Raw(sql, tree.ID, tree.UserID, tree.UserID).
		Scan(&trees).
		Error; err != nil {
		return nil, err
	}
	return trees, nil
}

func (fm *Repository) Update(ctx context.Context, tree dto.Tree) (dto.Tree, error) {
	logrus.Debugf("[input]: %+v", tree)

	tx := fm.db.WithContext(ctx).Model(dto.Tree{}).
		Where("id = ?", tree.ID).
		Where("user_id = ?", tree.UserID).Save(&tree)
	if tx.Error != nil {
		logrus.Errorf("[error]: %+v", tx.Error)
	}

	if tx.RowsAffected == 0 {
		return tree, gorm.ErrRecordNotFound
	}

	return tree, nil
}

func (fm *Repository) Delete(ctx context.Context, tree dto.Tree) (dto.Tree, error) {
	logrus.Debugf("[input]: %+v", tree)

	tx := fm.db.WithContext(ctx).Model(dto.Tree{}).
		Where("user_id = ?", tree.UserID).Delete(&tree)
	if tx.Error != nil {
		logrus.Errorf("[error]: %+v", tx.Error)
	}

	if tx.RowsAffected == 0 {
		return tree, gorm.ErrRecordNotFound
	}

	if err := fm.db.WithContext(ctx).
		Table("tree_documents").
		Where("tree_id = ?", tree.ID).
		Delete(&dto.TreeDocuments{}).Error; err != nil {
		return tree, err
	}

	// TODO: add logic of delete files after delete tree

	return tree, nil
}
