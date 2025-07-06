package model

import (
	"fmt"
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

func (*MasterAchievement) GetAchievement(id UUID, db *gorm.DB) ([]*AchievementResponse, error) {
	var achievements []*AchievementResponse

	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}

	err := db.
		Table("master_achievements").
		Select(`
			master_achievements.id,
			master_achievements.name,
			COALESCE(user_achievements.is_clear, false) AS is_clear
		`).
		Joins(`
			LEFT JOIN user_achievements 
			ON master_achievements.id = user_achievements.achievement_id 
			AND user_achievements.user_id = ?
		`, id).
		Scan(&achievements).Error

	return achievements, err
}

func (*MasterAchievement) Receipt(input InputAchievement, db *gorm.DB) (*UUID, error) {
	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}

	var userAchievement UserAchievement
	err := db.Where("user_id = ?", input.UserID).Where("achievement_id = ?", input.AchievementID).First(&userAchievement).Error

	if err == nil {
		return nil, fmt.Errorf("allready receipt")
	}

	achievement := UserAchievement{
		UserId:        input.UserID,
		AchievementId: input.AchievementID,
	}
	err = db.Model(&UserAchievement{}).Create(&achievement).Error

	if err != nil {
		return nil, err
	}

	return &achievement.Id, nil
}
