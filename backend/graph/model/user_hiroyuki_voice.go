package model

import (
	"time"

	"gorm.io/gorm"
)

type UserHiroyukiVoice struct {
	Id        UUID                `gorm:"primary_key; type: uuid; not null; default:uuid_generate_v4()"`
	UserId    UUID                `gorm:"type: uuid; not null"`
	User      User                `gorm:"foreignKey:UserId;references:Id"`
	VoiceId   UUID                `gorm:"type: uuid; not null"`
	Voice     MasterHiroyukiVoice `gorm:"foreignKey:VoiceId;references:Id"`
	IsHaving  bool                `gorm:"type: bool; not null; default: false"`
	CreatedAt time.Time           `gorm:"type: timestamp; autoCreateTime; not null; default:CURRENT_TIMESTAMP;<-:create"`
	UpdatedAt time.Time           `gorm:"type: timestamp; autoUpdateTime;<-:update"`
	DeletedAt gorm.DeletedAt      `gorm:"type: timestamp; index"`
}
