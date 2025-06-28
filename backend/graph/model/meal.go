package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/moXXcha/hiroyuki_diet_API/utils"
	"gorm.io/gorm"
)

type Meal struct {
	Id           uuid.UUID      `gorm:"primary_key; type: uuid; not null; default:uuid_generate_v4()"`
	UserId       uuid.UUID      `gorm:"type: uuid; not null"`
	User         User           `gorm:"foreignKey:UserId;references:Id"`
	MealType     utils.MealType `gorm:"type: meal_type; not null"`
	TotalCarorie int            `gorm:"type: int; not null"`
	Food         []Food         `gorm:"many2many:food_meals"`
	CreatedAt    time.Time      `gorm:"type: timestamp; autoCreateTime; not null; default:CURRENT_TIMESTAMP;<-:create"`
	UpdatedAt    time.Time      `gorm:"type: timestamp; autoUpdateTime;<-:update"`
	DeletedAt    gorm.DeletedAt `gorm:"type: timestamp; index"`
}
