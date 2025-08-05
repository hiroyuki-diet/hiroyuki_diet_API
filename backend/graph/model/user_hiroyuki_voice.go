package model

import (
	"errors"
	"fmt"
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

func (*UserHiroyukiVoice) Seeder(db *gorm.DB) error {
	var count int64

	// main.goが実行される度にレコードが生成されないようにする。
	db.Model(&UserHiroyukiVoice{}).Count(&count)
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

	var voice MasterHiroyukiVoice
	err = db.First(&voice).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("voice not found")
	}

	if err != nil {
		return err
	}

	userHiroyukiVoice := UserHiroyukiVoice{UserId: user.Id, VoiceId: voice.Id, IsHaving: true}

	err = db.Create(&userHiroyukiVoice).Error

	if err != nil {
		return err
	}

	return nil
}
