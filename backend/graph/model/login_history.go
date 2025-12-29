package model

import (
	"time"

	"gorm.io/gorm"
)

type LoginHistory struct {
	Id        UUID           `gorm:"primary_key; type: uuid; not null; default:uuid_generate_v4()"`
	UserId    UUID           `gorm:"type: uuid; not null"`
	User      User           `gorm:"foreignKey:UserId;references:Id"`
	LoginDate time.Time      `gorm:"type: date; not null"`
	CreatedAt time.Time      `gorm:"type: timestamp; autoCreateTime; not null; default:CURRENT_TIMESTAMP;<-:create"`
	DeletedAt gorm.DeletedAt `gorm:"type: timestamp; index"`
}

// RecordLogin records a login for the user on the current date
// Returns true if this is a new login for today, false if already logged in today
func RecordLogin(userId UUID, db *gorm.DB) (bool, error) {
	today := time.Now().Truncate(24 * time.Hour)

	// Check if already logged in today
	var count int64
	err := db.Model(&LoginHistory{}).
		Where("user_id = ? AND login_date = ? AND deleted_at IS NULL", userId, today).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	if count > 0 {
		return false, nil // Already logged in today
	}

	// Record new login
	loginHistory := LoginHistory{
		UserId:    userId,
		LoginDate: today,
	}
	err = db.Create(&loginHistory).Error
	if err != nil {
		return false, err
	}

	return true, nil
}

// GetLoginStreak returns the current consecutive login days for a user
func GetLoginStreak(userId UUID, db *gorm.DB) (int, error) {
	var loginDates []time.Time
	err := db.Model(&LoginHistory{}).
		Where("user_id = ? AND deleted_at IS NULL", userId).
		Order("login_date DESC").
		Pluck("login_date", &loginDates).Error
	if err != nil {
		return 0, err
	}

	if len(loginDates) == 0 {
		return 0, nil
	}

	streak := 1
	today := time.Now().Truncate(24 * time.Hour)

	// Check if the most recent login is today or yesterday
	lastLogin := loginDates[0].Truncate(24 * time.Hour)
	if lastLogin.Before(today.AddDate(0, 0, -1)) {
		// Last login was before yesterday, streak is broken
		return 0, nil
	}

	// Count consecutive days
	for i := 1; i < len(loginDates); i++ {
		currentDate := loginDates[i-1].Truncate(24 * time.Hour)
		prevDate := loginDates[i].Truncate(24 * time.Hour)

		expectedPrev := currentDate.AddDate(0, 0, -1)
		if prevDate.Equal(expectedPrev) {
			streak++
		} else {
			break
		}
	}

	return streak, nil
}
