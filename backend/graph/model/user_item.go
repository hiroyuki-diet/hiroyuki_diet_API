package model

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type UserItem struct {
	Id        UUID           `gorm:"primary_key; type: uuid; not null; default:uuid_generate_v4()"`
	UserId    UUID           `gorm:"type: uuid; not null"`
	User      User           `gorm:"foreignKey:UserId;references:Id"`
	ItemId    UUID           `gorm:"type: uuid; not null"`
	Item      MasterItem     `gorm:"foreignKey:ItemId;references:Id"`
	Count     int            `gorm:"type: int; not null; default:0"`
	CreatedAt time.Time      `gorm:"type: timestamp; autoCreateTime; not null; default:CURRENT_TIMESTAMP;<-:create"`
	UpdatedAt time.Time      `gorm:"type: timestamp; autoUpdateTime;<-:update"`
	DeletedAt gorm.DeletedAt `gorm:"type: timestamp; index"`
}

func (*UserItem) Seeder(db *gorm.DB) error {
	var count int64

	// main.goが実行される度にレコードが生成されないようにする。
	db.Model(&UserItem{}).Count(&count)
	if count > 0 {
		return nil
	}

	var user User
	err := db.First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("user not found")
	}

	if err != nil {
		return err
	}

	var item MasterItem
	err = db.First(&item).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("item not found")
	}

	if err != nil {
		return err
	}

	userItem := UserItem{UserId: user.Id, ItemId: item.Id, Count: 1}

	err = db.Create(&userItem).Error

	if err != nil {
		return err
	}

	return nil
}
