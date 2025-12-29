package achievement

import (
	"github.com/moXXcha/hiroyuki_diet_API/graph/model"
	"gorm.io/gorm"
)

type AchievementChecker struct {
	DB *gorm.DB
}

func NewAchievementChecker(db *gorm.DB) *AchievementChecker {
	return &AchievementChecker{DB: db}
}

// CheckExerciseAchievements checks and grants exercise-related achievements
// Returns list of newly achieved achievement names
func (c *AchievementChecker) CheckExerciseAchievements(userId model.UUID) ([]string, error) {
	var newlyAchieved []string

	// Get total exercise time for user
	var totalSeconds int64
	err := c.DB.Model(&model.Exercise{}).
		Where("user_id = ? AND deleted_at IS NULL", userId).
		Select("COALESCE(SUM(time), 0)").
		Scan(&totalSeconds).Error
	if err != nil {
		return nil, err
	}

	totalHours := float64(totalSeconds) / 3600.0

	// Get today's exercise time
	var todaySeconds int64
	err = c.DB.Model(&model.Exercise{}).
		Where("user_id = ? AND date = CURRENT_DATE AND deleted_at IS NULL", userId).
		Select("COALESCE(SUM(time), 0)").
		Scan(&todaySeconds).Error
	if err != nil {
		return nil, err
	}

	todayHours := float64(todaySeconds) / 3600.0

	// Define exercise achievements
	exerciseAchievements := []struct {
		Name      string
		Condition func() bool
	}{
		{"とりあえず動いてる", func() bool { return totalHours >= 1 }},
		{"ちゃんと続けてる", func() bool { return totalHours >= 10 }},
		{"なんか頑張ってる", func() bool { return totalHours >= 50 }},
		{"もうプロトレーナー？", func() bool { return totalHours >= 100 }},
		{"急にやる気出した？", func() bool { return todayHours >= 1 }},
	}

	for _, achievement := range exerciseAchievements {
		if achievement.Condition() {
			granted, err := c.grantAchievementIfNotExists(userId, achievement.Name)
			if err != nil {
				return nil, err
			}
			if granted {
				newlyAchieved = append(newlyAchieved, achievement.Name)
			}
		}
	}

	return newlyAchieved, nil
}

// grantAchievementIfNotExists grants an achievement if user doesn't have it yet
// Returns true if newly granted, false if already exists
func (c *AchievementChecker) grantAchievementIfNotExists(userId model.UUID, achievementName string) (bool, error) {
	// Find the achievement by name
	var achievement model.MasterAchievement
	err := c.DB.Where("name = ? AND deleted_at IS NULL", achievementName).First(&achievement).Error
	if err != nil {
		return false, err
	}

	// Check if user already has this achievement
	var count int64
	err = c.DB.Model(&model.UserAchievement{}).
		Where("user_id = ? AND achievement_id = ? AND deleted_at IS NULL", userId, achievement.Id).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	if count > 0 {
		return false, nil // Already has achievement
	}

	// Grant the achievement
	userAchievement := model.UserAchievement{
		UserId:        userId,
		AchievementId: achievement.Id,
	}
	err = c.DB.Create(&userAchievement).Error
	if err != nil {
		return false, err
	}

	return true, nil
}

// CheckLoginAchievements checks and grants login streak achievements
// Returns list of newly achieved achievement names
func (c *AchievementChecker) CheckLoginAchievements(userId model.UUID) ([]string, error) {
	var newlyAchieved []string

	// Get current login streak
	streak, err := model.GetLoginStreak(userId, c.DB)
	if err != nil {
		return nil, err
	}

	// Define login achievements
	loginAchievements := []struct {
		Name      string
		Condition func() bool
	}{
		{"はじめまして", func() bool { return streak >= 1 }},
		{"三日坊主じゃない人", func() bool { return streak >= 3 }},
		{"いい感じに続けてる人", func() bool { return streak >= 7 }},
		{"意外と根性あるかも", func() bool { return streak >= 14 }},
		{"ここまで続くとは…", func() bool { return streak >= 30 }},
		{"本気でやる気になった？", func() bool { return streak >= 90 }},
	}

	for _, achievement := range loginAchievements {
		if achievement.Condition() {
			granted, err := c.grantAchievementIfNotExists(userId, achievement.Name)
			if err != nil {
				// Skip if achievement not found in master data
				continue
			}
			if granted {
				newlyAchieved = append(newlyAchieved, achievement.Name)
			}
		}
	}

	return newlyAchieved, nil
}
