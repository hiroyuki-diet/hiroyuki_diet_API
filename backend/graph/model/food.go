package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Food struct {
	Id              uuid.UUID      `gorm:"primary_key; type: uuid; not null; default:uuid_generate_v4()"`
	Name            string         `gorm:"type: varchar(50); not null"`
	EstimateCarorie int            `gorm:"type: int; not null"`
	LastUsedDate    time.Time      `gorm:"type: date; not null"`
	CreatedAt       time.Time      `gorm:"type: timestamp; autoCreateTime; not null; default:CURRENT_TIMESTAMP;<-:create"`
	UpdatedAt       time.Time      `gorm:"type: timestamp; autoUpdateTime;<-:update"`
	DeletedAt       gorm.DeletedAt `gorm:"type: timestamp; index"`
}
