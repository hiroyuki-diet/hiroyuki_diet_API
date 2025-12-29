package model

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

var jstWeight = time.FixedZone("JST", 9*60*60)

type WeightHistory struct {
	Id        UUID           `gorm:"primary_key; type: uuid; not null; default:uuid_generate_v4()"`
	UserId    UUID           `gorm:"type: uuid; not null"`
	User      User           `gorm:"foreignKey:UserId;references:Id"`
	Weight    int            `gorm:"type: int; not null"`
	Date      time.Time      `gorm:"type: date; not null; default:CURRENT_TIMESTAMP;<-:create"`
	CreatedAt time.Time      `gorm:"type: timestamp; autoCreateTime; not null; default:CURRENT_TIMESTAMP;<-:create"`
	UpdatedAt time.Time      `gorm:"type: timestamp; autoUpdateTime;<-:update"`
	DeletedAt gorm.DeletedAt `gorm:"type: timestamp; index"`
}

func (*WeightHistory) GetHistories(offset string, limit string, userId UUID, db *gorm.DB) ([]*WeightHistory, error) {
	var histories []*WeightHistory

	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}

	err := db.Where("user_id = ?", userId).Where("date BETWEEN ? AND ?", offset, limit).Order("date asc").Find(&histories).Error

	if err != nil {
		return nil, err
	}

	return histories, nil
}

func (*WeightHistory) Create(userId UUID, weight int, db *gorm.DB) (*UUID, error) {
	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}

	// Check if there's already an entry for today (use JST for consistency with client)
	now := time.Now().In(jstWeight)
	today := now.Format("2006-01-02")
	var existing WeightHistory
	err := db.Where("user_id = ? AND date = ?", userId, today).First(&existing).Error

	if err == nil {
		// Update existing entry
		if err := db.Model(&WeightHistory{}).Where("id = ?", existing.Id).Update("weight", weight).Error; err != nil {
			return nil, err
		}
		return &existing.Id, nil
	}

	// Create new entry with explicit JST date
	history := WeightHistory{
		UserId: userId,
		Weight: weight,
		Date:   now,
	}

	if err := db.Create(&history).Error; err != nil {
		return nil, err
	}

	return &history.Id, nil
}
