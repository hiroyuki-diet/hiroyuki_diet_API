package model

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Food struct {
	Id              UUID           `gorm:"primary_key; type: uuid; not null; default:uuid_generate_v4()"`
	Name            string         `gorm:"type: varchar(50); not null"`
	EstimateCalorie int            `gorm:"type: int; not null"`
	LastUsedDate    time.Time      `gorm:"type: date; not null"`
	CreatedAt       time.Time      `gorm:"type: timestamp; autoCreateTime; not null; default:CURRENT_TIMESTAMP;<-:create"`
	UpdatedAt       time.Time      `gorm:"type: timestamp; autoUpdateTime;<-:update"`
	DeletedAt       gorm.DeletedAt `gorm:"type: timestamp; index"`
}

func (*Food) GetAll(db *gorm.DB) ([]*Food, error) {
	var foods []*Food

	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}

	result := db.Find(&foods)

	if result.Error != nil {
		return nil, result.Error
	}

	return foods, nil
}
