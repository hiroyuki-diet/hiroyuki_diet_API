package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/moXXcha/hiroyuki_diet_API/utils"
	"gorm.io/gorm"
)

type MasterFields struct {
	Id        uuid.UUID      `gorm:"primary_key; type: uuid; not null; default:uuid_generate_v4()"`
	Field     utils.Field    `gorm:"type: field; not null"`
	CreatedAt time.Time      `gorm:"type: timestamp; autoCreateTime; not null; default:CURRENT_TIMESTAMP;<-:create"`
	UpdatedAt time.Time      `gorm:"type: timestamp; autoUpdateTime;<-:update"`
	DeletedAt gorm.DeletedAt `gorm:"type: timestamp; index"`
}

func (*MasterFields) FirstCreate(db *gorm.DB) error {
	fields := []MasterFields{
		{
			Field: "login",
		},
		{
			Field: "signin",
		},
		{
			Field: "home",
		},
		{
			Field: "meal",
		},
		{
			Field: "meal_form",
		},
		{
			Field: "meal_edit",
		},
		{
			Field: "data",
		},
		{
			Field: "profile",
		},
		{
			Field: "profile_edit",
		},
		{
			Field: "exercise",
		},
		{
			Field: "exercise_complete",
		},
		{
			Field: "achievement",
		},
		{
			Field: "achievement_complete",
		},
		{
			Field: "chibi_hiroyuki",
		},
	}

	for i := range fields {
		result := db.FirstOrCreate(&fields[i], MasterFields{Field: fields[i].Field})
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}
