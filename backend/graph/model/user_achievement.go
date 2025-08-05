package model

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type UserAchievement struct {
	Id            UUID              `gorm:"primary_key; type: uuid; not null; default:uuid_generate_v4()"`
	UserId        UUID              `gorm:"type: uuid; not null"`
	User          User              `gorm:"foreignKey:UserId;references:Id"`
	AchievementId UUID              `gorm:"type: uuid; not null"`
	Achievement   MasterAchievement `gorm:"foreignKey:AchievementId;references:Id"`
	CreatedAt     time.Time         `gorm:"type: timestamp; autoCreateTime; not null; default:CURRENT_TIMESTAMP;<-:create"`
	UpdatedAt     time.Time         `gorm:"type: timestamp; autoUpdateTime;<-:update"`
	DeletedAt     gorm.DeletedAt    `gorm:"type: timestamp; index"`
}

func (*UserAchievement) Seeder(db *gorm.DB) error {
	var count int64

	// main.goが実行される度にレコードが生成されないようにする。
	db.Model(&UserAchievement{}).Count(&count)
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

	var achievement MasterAchievement
	err = db.First(&achievement).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("achievement not found")
	}

	if err != nil {
		return nil
	}

	userAchievement := UserAchievement{UserId: user.Id, AchievementId: achievement.Id}

	err = db.Create(&userAchievement).Error

	if err != nil {
		return err
	}

	return nil
}
