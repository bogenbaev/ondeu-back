package dto

import (
	"github.com/google/uuid"
	"io"
	"time"
)

type Document struct {
	ID              uint          `json:"id" gorm:"<-:create;primarykey;"`
	UserID          string        `json:"-"  gorm:"<-:create;varchar(50)"`
	TreeID          uint          `json:"-"  gorm:"->;-:migration;column:tree_id"`
	CreatedAt       time.Time     `json:"createdAt" gorm:"<-:create;"`
	UpdatedAt       time.Time     `json:"updatedAt"`
	Name            string        `form:"name,omitempty" json:"name,omitempty" gorm:"varchar(2000)"`
	Extension       string        `json:"extension,omitempty" gorm:"varchar(10);<-:create"`
	Size            int64         `json:"size,omitempty" gorm:"number;<-:create;"`
	Type            string        `json:"type,omitempty" gorm:"varchar(255);<-:create"`
	Path            uuid.UUID     `json:"path,omitempty" gorm:"<-:create;type:uuid;default:gen_random_uuid()"`
	Template        *bool         `json:"template" form:"template,omitempty" gorm:"default:false"`
	ShareLink       string        `json:"shareLink,omitempty" gorm:"-:all"`
	RequestContent  io.ReadSeeker `gorm:"-:all" json:"-"`
	ResponseContent []byte        `gorm:"-:all" json:"-"`
}
