package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/moXXcha/hiroyuki_diet_API/utils"
	"gorm.io/gorm"
)

type Profile struct {
	Id                      uuid.UUID      `gorm:"primary_key; type: uuid; not null; default:uuid_generate_v4()"`
	UserId                  uuid.UUID      `gorm:"type: uuid; not null"`
	User                    User           `gorm:"foreignKey:UserId;references:Id"`
	UserName                string         `gorm:"type: varchar(50); not null"`
	Age                     int            `gorm:"type: int; not null"`
	Gender                  utils.Gender   `gorm:"type: gender; not null"`
	Weight                  int            `gorm:"type: int; not null"`
	Height                  int            `gorm:"type: int; not null"`
	TargetWeight            int            `gorm:"type: int; not null"`
	TargetDailyExerciseTime int            `gorm:"type: int; not null"`
	TargetDailyCarorie      int            `gorm:"type: int; not null"`
	Favorability            int            `gorm:"type: int; not null"`
	IsCreated               bool           `gorm:"type: bool; not null; default: false"`
	CreatedAt               time.Time      `gorm:"type: timestamp; autoCreateTime; not null; default:CURRENT_TIMESTAMP;<-:create"`
	UpdatedAt               time.Time      `gorm:"type: timestamp; autoUpdateTime;<-:update"`
	DeletedAt               gorm.DeletedAt `gorm:"type: timestamp; index"`
}
