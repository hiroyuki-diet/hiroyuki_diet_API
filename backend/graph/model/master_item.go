package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MasterItem struct {
	Id          uuid.UUID      `gorm:"primary_key; type: uuid; not null; default:uuid_generate_v4()"`
	Name        string         `gorm:"type: varchar(50); not null"`
	Description string         `gorm:"type: text; not null"`
	ItemImage   string         `gorm:"type: varchar(50); not null"`
	CreatedAt   time.Time      `gorm:"type: timestamp; autoCreateTime; not null; default:CURRENT_TIMESTAMP;<-:create"`
	UpdatedAt   time.Time      `gorm:"type: timestamp; autoUpdateTime;<-:update"`
	DeletedAt   gorm.DeletedAt `gorm:"type: timestamp; index"`
}

func (*MasterItem) FirstCreate(db *gorm.DB) error {
	items := []MasterItem{
		{
			Name:        "チートデイチケット",
			Description: "頑張ったごほうび！使ったら一日やすんでいいよ",
			ItemImage:   "",
		},
	}

	for i := range items {
		result := db.FirstOrCreate(&items[i], MasterItem{Name: items[i].Name})
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}
