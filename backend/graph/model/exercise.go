package model

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Exercise struct {
	Id        UUID           `gorm:"primary_key; type: uuid; not null; default:uuid_generate_v4()"`
	UserId    UUID           `gorm:"type: uuid; not null"`
	User      User           `gorm:"foreignKey:UserId;references:Id"`
	Time      int            `gorm:"type: int; not null"`
	Date      time.Time      `gorm:"type: date; not null; default:CURRENT_TIMESTAMP;<-:create"`
	CreatedAt time.Time      `gorm:"type: timestamp; autoCreateTime; not null; default:CURRENT_TIMESTAMP;<-:create"`
	UpdatedAt time.Time      `gorm:"type: timestamp; autoUpdateTime;<-:update"`
	DeletedAt gorm.DeletedAt `gorm:"type: timestamp; index"`
}

func (*Exercise) Seeder(db *gorm.DB) error {
	var count int64

	// main.goが実行される度にレコードが生成されないようにする。
	db.Model(&Exercise{}).Count(&count)
	if count > 0 {
		return nil
	}

	var user User
	err := db.First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("user not found")
	}

	if err != nil {
		return nil
	}

	exercise := Exercise{UserId: user.Id, Time: 1}

	err = db.Create(&exercise).Error

	if err != nil {
		return err
	}

	return nil
}

func (*Exercise) GetInfo(id UUID, offset string, limit string, db *gorm.DB) ([]*Exercise, error) {
	var exercises []*Exercise

	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}

	err := db.
		Where("user_id = ?", id).
		Where("date BETWEEN ? AND ?", offset, limit).
		Order("date asc").
		Find(&exercises).Error

	if err != nil {
		return nil, err
	}

	return exercises, nil
}
