package dto

import "time"

type Group struct {
	ID        uint       `gorm:"primarykey" json:"id,omitempty"`
	CreatedAt time.Time  `json:"createdAt,omitempty"`
	UpdatedAt time.Time  `json:"updatedAt,omitempty"`
	DeletedAt time.Time  `json:"deletedAt,omitempty"`
	Name      string     `gorm:"varchar(2000)" json:"name" binding:"required"`
	Desc      string     `gorm:"text" json:"desc,omitempty"`
	Role      string     `json:"role" gorm:"varchar(255)"`
	Documents []Document `json:"documents,omitempty" gorm:"many2many:group_documents;"`
}
