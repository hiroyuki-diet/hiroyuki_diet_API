package model

import (
	"time"

	"gorm.io/gorm"
)

type Exercise struct {
	Id        UUID           `gorm:"primary_key; type: uuid; not null; default:uuid_generate_v4()"`
	UserId    UUID           `gorm:"type: uuid; not null"`
	User      User           `gorm:"foreignKey:UserId;references:Id"`
	Time      int            `gorm:"type: int; not null"`
	CreatedAt time.Time      `gorm:"type: timestamp; autoCreateTime; not null; default:CURRENT_TIMESTAMP;<-:create"`
	UpdatedAt time.Time      `gorm:"type: timestamp; autoUpdateTime;<-:update"`
	DeletedAt gorm.DeletedAt `gorm:"type: timestamp; index"`
}
