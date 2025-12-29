package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/moXXcha/hiroyuki_diet_API/utils"
	"gorm.io/gorm"
)

type Profile struct {
	Id                      UUID           `gorm:"primary_key; type: uuid; not null; default:uuid_generate_v4()"`
	UserId                  UUID           `gorm:"type: uuid; not null"`
	User                    User           `gorm:"foreignKey:UserId;references:Id"`
	UserName                string         `gorm:"type: varchar(50); not null"`
	Age                     int            `gorm:"type: int; not null"`
	Gender                  utils.Gender   `gorm:"type: gender; not null"`
	Weight                  int            `gorm:"type: int; not null"`
	Height                  int            `gorm:"type: int; not null"`
	TargetWeight            int            `gorm:"type: int; not null"`
	TargetDailyExerciseTime int            `gorm:"type: int; not null"`
	TargetDailyCarorie      int            `gorm:"type: int; not null"`
	Favorability            int            `gorm:"type: int; not null"`
	IsCreated               bool           `gorm:"type: bool; not null; default: false"`
	CreatedAt               time.Time      `gorm:"type: timestamp; autoCreateTime; not null; default:CURRENT_TIMESTAMP;<-:create"`
	UpdatedAt               time.Time      `gorm:"type: timestamp; autoUpdateTime;<-:update"`
	DeletedAt               gorm.DeletedAt `gorm:"type: timestamp; index"`
}

func (*Profile) GetInfo(id UUID, db *gorm.DB) (*Profile, error) {
	var profile *Profile

	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}

	result := db.Where("user_id = ?", id).First(&profile)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("PROFILE_NOT_FOUND")
		}
		return nil, result.Error
	}

	return profile, nil
}

func (*Profile) Seeder(db *gorm.DB) error {
	var count int64

	// main.goが実行される度にレコードが生成されないようにする。
	db.Model(&Profile{}).Count(&count)
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

	profile := Profile{UserId: user.Id, UserName: "こなみるく", Age: 30, Gender: "woman", Weight: 30, Height: 165, TargetWeight: 20, TargetDailyExerciseTime: 1, TargetDailyCarorie: 1000, Favorability: 1, IsCreated: true}

	err = db.Create(&profile).Error

	if err != nil {
		return err
	}

	return nil
}

func (*Profile) Create(input InputProfile, db *gorm.DB) (*UUID, error) {
	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}

	profile := Profile{
		UserName:                input.UserName,
		UserId:                  input.UserID,
		Age:                     input.Age,
		Gender:                  utils.Gender(input.Gender),
		Weight:                  input.Weight,
		Height:                  input.Height,
		TargetWeight:            input.TargetWeight,
		TargetDailyCarorie:      input.TargetDailyCarorie,
		TargetDailyExerciseTime: input.TargetDailyExerciseTime,
		Favorability:            0,
		IsCreated:               true,
	}

	err := db.Create(&profile).Error

	if err != nil {
		return nil, err
	}

	return &profile.Id, nil
}

func (*Profile) Edit(input InputProfile, db *gorm.DB) (*UUID, error) {
	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}

	var profileId UUID
	err := db.Transaction(func(tx *gorm.DB) error {
		var profile Profile
		if err := tx.Where("user_id = ?", input.UserID).First(&profile).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("profile not found")
			}
			return err
		}

		profileInput := Profile{
			UserName:                input.UserName,
			UserId:                  input.UserID,
			Age:                     input.Age,
			Gender:                  utils.Gender(input.Gender),
			Weight:                  input.Weight,
			Height:                  input.Height,
			TargetWeight:            input.TargetWeight,
			TargetDailyCarorie:      input.TargetDailyCarorie,
			TargetDailyExerciseTime: input.TargetDailyExerciseTime,
			Favorability:            0,
			IsCreated:               true,
		}

		if err := tx.Model(&Profile{}).Where("user_id = ?", input.UserID).Updates(&profileInput).Error; err != nil {
			return err
		}

		profileId = profile.Id
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &profileId, nil
}
