package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserHiroyukiVoice struct {
	Id        uuid.UUID           `gorm:"primary_key; type: uuid; not null; default:uuid_generate_v4()"`
	UserId    uuid.UUID           `gorm:"type: uuid; not null"`
	User      User                `gorm:"foreignKey:UserId;references:Id"`
	VoiceId   uuid.UUID           `gorm:"type: uuid; not null"`
	Voice     MasterHiroyukiVoice `gorm:"foreignKey:VoiceId;references:Id"`
	IsHaving  bool                `gorm:"type: bool; not null; default: false"`
	CreatedAt time.Time           `gorm:"type: timestamp; autoCreateTime; not null; default:CURRENT_TIMESTAMP;<-:create"`
	UpdatedAt time.Time           `gorm:"type: timestamp; autoUpdateTime;<-:update"`
	DeletedAt gorm.DeletedAt      `gorm:"type: timestamp; index"`
}
