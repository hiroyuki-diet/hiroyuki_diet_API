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

func (*Food) Seeder(db *gorm.DB) error {
	var count int64

	// main.goが実行される度にレコードが生成されないようにする。
	db.Model(&Food{}).Count(&count)
	if count > 0 {
		return nil
	}

	t, err := time.Parse("2006-01-02", "2025-07-01")

	if err != nil {
		return err
	}

	food := Food{Name: "たこやき", EstimateCalorie: 1000, LastUsedDate: t}
	err = db.Create(&food).Error

	if err != nil {
		return err
	}

	return nil
}
