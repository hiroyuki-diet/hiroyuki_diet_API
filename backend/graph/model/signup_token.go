package model

import (
	"time"

	"gorm.io/gorm"
)

type SignUpToken struct {
	Id          UUID           `gorm:"primary_key; type: uuid; not null; default:uuid_generate_v4()"`
	Token       int            `gorm:"type: int; not null"`
	SurviveTime int            `gorm:"type: int; not null"`
	CreatedAt   time.Time      `gorm:"type: timestamp; autoCreateTime; not null; default:CURRENT_TIMESTAMP;<-:create"`
	UpdatedAt   time.Time      `gorm:"type: timestamp; autoUpdateTime;<-:update"`
	DeletedAt   gorm.DeletedAt `gorm:"type: timestamp; index"`
}

func (*SignUpToken) Seeder(db *gorm.DB) error {
	var count int64

	// main.goが実行される度にレコードが生成されないようにする。
	db.Model(&SignUpToken{}).Count(&count)
	if count > 0 {
		return nil
	}

	signUpToken := SignUpToken{Token: 123456, SurviveTime: 1}

	err := db.Create(&signUpToken).Error

	if err != nil {
		return err
	}

	return nil
}
