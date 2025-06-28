package model

import (
	"time"

	"gorm.io/gorm"
)

type UserSkin struct {
	Id        UUID               `gorm:"primary_key; type: uuid; not null; default:uuid_generate_v4()"`
	UserId    UUID               `gorm:"type: uuid; not null"`
	User      User               `gorm:"foreignKey:UserId;references:Id"`
	SkinId    UUID               `gorm:"type: uuid; not null"`
	Skin      MasterHiroyukiSkin `gorm:"foreignKey:SkinId;references:Id"`
	IsUsing   bool               `gorm:"type: bool; not null; default: false"`
	IsHaving  bool               `gorm:"type: bool; not null; default: false"`
	CreatedAt time.Time          `gorm:"type: timestamp; autoCreateTime; not null; default:CURRENT_TIMESTAMP;<-:create"`
	UpdatedAt time.Time          `gorm:"type: timestamp; autoUpdateTime;<-:update"`
	DeletedAt gorm.DeletedAt     `gorm:"type: timestamp; index"`
}
