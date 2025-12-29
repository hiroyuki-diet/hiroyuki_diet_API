package model

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Exercise struct {
	Id          UUID           `gorm:"primary_key; type: uuid; not null; default:uuid_generate_v4()"`
	UserId      UUID           `gorm:"type: uuid; uniqueIndex:idx_user_date; not null"`
	User        User           `gorm:"foreignKey:UserId;references:Id"`
	Time        int            `gorm:"type: int; not null"`
	IsCompleted bool           `gorm:"type: boolean; not null; default:false"`
	Date        time.Time      `gorm:"type: date; uniqueIndex:idx_user_date; not null; default:CURRENT_TIMESTAMP;<-:create"`
	CreatedAt   time.Time      `gorm:"type: timestamp; autoCreateTime; not null; default:CURRENT_TIMESTAMP;<-:create"`
	UpdatedAt   time.Time      `gorm:"type: timestamp; autoUpdateTime;<-:update"`
	DeletedAt   gorm.DeletedAt `gorm:"type: timestamp; index"`
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

func (*Exercise) Create(input InputExercise, db *gorm.DB) (*UUID, error) {
	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}
	todayStr := time.Now().Format("2006-01-02")

	layout := "2006-01-02"

	t, err := time.Parse(layout, todayStr)
	if err != nil {
		return nil, err
	}

	isCompleted := false
	if input.IsCompleted != nil {
		isCompleted = *input.IsCompleted
	}

	exercise := Exercise{
		UserId:      *input.UserID,
		Time:        input.Time,
		IsCompleted: isCompleted,
		Date:        t,
	}

	err = db.Create(&exercise).Error

	if err != nil {
		return nil, err
	}

	return &exercise.Id, nil
}

func (*Exercise) Edit(input InputExercise, db *gorm.DB) (*UUID, error) {
	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}

	var exerciseId UUID
	err := db.Transaction(func(tx *gorm.DB) error {
		todayStr := time.Now().Format("2006-01-02")
		layout := "2006-01-02"

		t, err := time.Parse(layout, todayStr)
		if err != nil {
			return err
		}

		var exercise Exercise

		if err := tx.Model(&Exercise{}).Where("user_id = ?", input.UserID).Where("date = ?", t).First(&exercise).Error; err != nil {
			return err
		}

		updates := map[string]interface{}{
			"time": input.Time,
		}
		if input.IsCompleted != nil {
			updates["is_completed"] = *input.IsCompleted
		}

		if err := tx.Model(&Exercise{}).Where("user_id = ?", input.UserID).Where("date = ?", t).Updates(updates).Error; err != nil {
			return err
		}

		if err := tx.Where("user_id = ?", input.UserID).Where("date = ?", t).First(&exercise).Error; err != nil {
			return err
		}

		exerciseId = exercise.Id
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &exerciseId, nil
}
