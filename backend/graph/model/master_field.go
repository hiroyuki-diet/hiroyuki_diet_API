package model

import (
	"time"

	"github.com/moXXcha/hiroyuki_diet_API/utils"
	"gorm.io/gorm"
)

type MasterField struct {
	Id        UUID           `gorm:"primary_key; type: uuid; not null; default:uuid_generate_v4()"`
	Field     utils.Field    `gorm:"type: field; unique; not null"`
	CreatedAt time.Time      `gorm:"type: timestamp; autoCreateTime; not null; default:CURRENT_TIMESTAMP;<-:create"`
	UpdatedAt time.Time      `gorm:"type: timestamp; autoUpdateTime;<-:update"`
	DeletedAt gorm.DeletedAt `gorm:"type: timestamp; index"`
}

func (*MasterField) FirstCreate(db *gorm.DB) error {
	fields := []MasterField{
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
		err := db.FirstOrCreate(&fields[i], "field = ?", fields[i].Field).Error
		if err != nil {
			return err
		}
	}

	return nil
}
