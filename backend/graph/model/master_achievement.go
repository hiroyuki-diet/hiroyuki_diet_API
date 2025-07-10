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
	return db.Transaction(func(tx *gorm.DB) error {
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
			if err := tx.FirstOrCreate(&achievements[i], MasterAchievement{Name: achievements[i].Name}).Error; err != nil {
				return err
			}
		}
		return nil
	})
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

	var achievementId UUID
	err := db.Transaction(func(tx *gorm.DB) error {
		var userAchievement UserAchievement
		if err := tx.Where("user_id = ?", input.UserID).Where("achievement_id = ?", input.AchievementID).First(&userAchievement).Error; err == nil {
			return fmt.Errorf("already receipt")
		}

		achievement := UserAchievement{
			UserId:        input.UserID,
			AchievementId: input.AchievementID,
		}
		if err := tx.Model(&UserAchievement{}).Create(&achievement).Error; err != nil {
			return err
		}

		achievementId = achievement.Id
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &achievementId, nil
}
