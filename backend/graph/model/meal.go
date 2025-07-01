package model

import (
	"fmt"
	"time"

	"github.com/moXXcha/hiroyuki_diet_API/utils"
	"gorm.io/gorm"
)

type Meal struct {
	Id           UUID           `gorm:"primary_key; type: uuid; not null; default:uuid_generate_v4()"`
	UserId       UUID           `gorm:"type: uuid; not null"`
	User         User           `gorm:"foreignKey:UserId;references:Id"`
	MealType     utils.MealType `gorm:"type: meal_type; not null"`
	TotalCalorie int            `gorm:"type: int; not null"`
	Foods        []Food         `gorm:"many2many:food_meals"`
	CreatedAt    time.Time      `gorm:"type: timestamp; autoCreateTime; not null; default:CURRENT_TIMESTAMP;<-:create"`
	UpdatedAt    time.Time      `gorm:"type: timestamp; autoUpdateTime;<-:update"`
	DeletedAt    gorm.DeletedAt `gorm:"type: timestamp; index"`
}

func (*Meal) GetAll(id UUID, db *gorm.DB) ([]*Meal, error) {
	var meals []*Meal

	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}

	err := db.Preload("Foods").Where("user_id = ?", id).Find(&meals).Error

	if err != nil {
		return nil, err
	}

	return meals, nil
}
