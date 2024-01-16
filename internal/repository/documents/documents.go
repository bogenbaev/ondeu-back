package documents

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

func (fm *Repository) Get(ctx context.Context, doc dto.Document) (dto.Document, error) {
	logrus.Debugf("[input]: %+v", doc)

	var count int64
	if err := fm.db.WithContext(ctx).
		Table("tree_documents").
		Where("tree_id = ?", doc.TreeID).
		Where("document_id = ?", doc.ID).
		Count(&count).Error; err != nil {
		return doc, err
	}

	logrus.Debugf("[count of documents]: %+v", count)
	if count == 0 {
		return doc, gorm.ErrRecordNotFound
	}

	var documentCount int64
	if err := fm.db.WithContext(ctx).
		Model(dto.Document{}).
		Where("user_id = ?", doc.UserID).
		Find(&doc).Count(&documentCount).Error; err != nil {
		logrus.Errorf("[error]: %+v", err)
	}

	if documentCount == 0 {
		return doc, gorm.ErrRecordNotFound
	}

	return doc, nil
}

func (fm *Repository) Create(ctx context.Context, doc dto.Document) (dto.Document, error) {
	logrus.Debugf("[input]: %+v", doc)
	if err := fm.db.WithContext(ctx).
		Model(dto.Document{}).
		Create(&doc).
		Error; err != nil {
		return doc, err
	}

	return doc, fm.db.WithContext(ctx).
		Table("tree_documents").
		Create(&dto.TreeDocuments{TreeID: doc.TreeID, DocumentID: doc.ID}).
		Error
}

func (fm *Repository) ListByTree(ctx context.Context, ids []uint) ([]dto.Document, error) {
	sql := `select * from documents d 
    		join tree_documents td 
    		on d.id = td.document_id 
         	where td.tree_id in ?;`

	var document []dto.Document
	if err := fm.db.WithContext(ctx).
		Model(dto.Document{}).
		Raw(sql, ids).
		Scan(&document).
		Error; err != nil {
		return nil, err
	}
	return document, nil
}

func (fm *Repository) ListByGroups(ctx context.Context, ids []uint) ([]dto.Document, error) {
	sql := `select * from documents d 
    		join group_documents gd 
    		on d.id = gd.document_id 
         	where gd.group_id in ?;`

	var document []dto.Document
	if err := fm.db.WithContext(ctx).
		Model(dto.Document{}).
		Raw(sql, ids).
		Scan(&document).
		Error; err != nil {
		return nil, err
	}
	return document, nil
}

func (fm *Repository) Delete(ctx context.Context, doc dto.Document) (dto.Document, error) {
	logrus.Debugf("[input]: %+v", doc)

	res := fm.db.WithContext(ctx).
		Table("tree_documents").
		Where("tree_id = ?", doc.TreeID).
		Where("document_id = ?", doc.ID).
		Delete(&dto.TreeDocuments{})
	if res.Error != nil {
		return doc, res.Error
	}

	if res.RowsAffected == 0 {
		return doc, gorm.ErrRecordNotFound
	}

	tx := fm.db.WithContext(ctx).Model(dto.Document{}).
		Where("user_id = ?", doc.UserID).Delete(&doc)
	if tx.Error != nil {
		logrus.Errorf("[error]: %+v", tx.Error)
	}

	if tx.RowsAffected == 0 {
		return doc, gorm.ErrRecordNotFound
	}

	return doc, nil
}

func (fm *Repository) Update(ctx context.Context, doc dto.Document) (dto.Document, error) {
	logrus.Debugf("[input]: %+v", doc)

	tx := fm.db.WithContext(ctx).Model(dto.Document{}).
		Where("id = ?", doc.ID).
		Where("user_id = ?", doc.UserID).Save(&doc)
	if tx.Error != nil {
		logrus.Errorf("[error]: %+v", tx.Error)
	}

	if tx.RowsAffected == 0 {
		return doc, gorm.ErrRecordNotFound
	}

	return doc, nil
}

func (fm *Repository) FindByCondition(ctx context.Context, field, param string) ([]dto.Document, error) {
	logrus.Debugf("[input]: %+v, %+v", field, param)

	docs := make([]dto.Document, 0)
	tx := fm.db.WithContext(ctx).Model(dto.Document{}).Find(&docs).Where("?=%?%", field, param).Save(&docs)
	if tx != nil {
		logrus.Errorf("[error]: %+v", tx.Error)
	}

	if tx.RowsAffected == 0 {
		return docs, gorm.ErrRecordNotFound
	}

	return docs, nil
}
