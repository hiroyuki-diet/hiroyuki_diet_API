package model

import (
	"time"

	"gorm.io/gorm"
)

type MasterAchievement struct {
	Id          UUID           `gorm:"primary_key; type: uuid; not null; default:uuid_generate_v4()"`
	Name        string         `gorm:"type: varchar(50); not null"`
	Description string         `gorm:"type: text; not null"`
	CreatedAt   time.Time      `gorm:"type: timestamp; autoCreateTime; not null; default:CURRENT_TIMESTAMP;<-:create"`
	UpdatedAt   time.Time      `gorm:"type: timestamp; autoUpdateTime;<-:update"`
	DeletedAt   gorm.DeletedAt `gorm:"type: timestamp; index"`
}

func (*MasterAchievement) FirstCreate(db *gorm.DB) error {
	achievements := []MasterAchievement{
		{
			Name:        "初ログイン",
			Description: "初回ログイン達成実績",
		},
		{
			Name:        "レベル5達成",
			Description: "レベル5達成実績",
		},
	}

	for i := range achievements {
		result := db.FirstOrCreate(&achievements[i], MasterAchievement{Name: achievements[i].Name})
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}
