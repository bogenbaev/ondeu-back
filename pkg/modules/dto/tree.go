package dto

import (
	"time"
)

type Tree struct {
	ID        uint       `gorm:"<-:create;primarykey" json:"id,omitempty"`
	UserID    string     `json:"-" gorm:"<-:create;varchar(255)"`
	DocID     uint       `gorm:"<-:create;foreignkey" json:"-"`
	ParentID  uint       `json:"parentID" gorm:"<-:create;"`
	CreatedAt time.Time  `json:"createdAt,omitempty" gorm:"<-:create"`
	UpdatedAt time.Time  `json:"updatedAt,omitempty"`
	Name      string     `gorm:"varchar(2000)" json:"name" binding:"required"`
	Role      string     `json:"role" form:"role" binding:"required" gorm:"varchar(255)"`
	Template  *bool      `json:"template" form:"template,omitempty"  gorm:"default:false"`
	Group     *bool      `json:"group" form:"group,omitempty" gorm:"default:false"`
	Documents []Document `json:"documents" gorm:"many2many:tree_documents;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type TreeDocuments struct {
	TreeID     uint `gorm:"primarykey"`
	DocumentID uint `gorm:"primarykey"`
}
