package model

import (
	"time"

	"gorm.io/gorm"
)

type MasterExercise struct {
	Id        UUID           `gorm:"primary_key; type: uuid; not null; default:uuid_generate_v4()"`
	Name      string         `gorm:"type: varchar(50); not null"`
	Mets      int            `gorm:"type: int; not null"`
	CreatedAt time.Time      `gorm:"type: timestamp; autoCreateTime; not null; default:CURRENT_TIMESTAMP;<-:create"`
	UpdatedAt time.Time      `gorm:"type: timestamp; autoUpdateTime;<-:update"`
	DeletedAt gorm.DeletedAt `gorm:"type: timestamp; index"`
}

func (*MasterExercise) FirstCreate(db *gorm.DB) error {
	exercisies := []MasterExercise{
		{
			Name: "腹筋10回",
			Mets: 5,
		},
		{
			Name: "腕立て10回",
			Mets: 7,
		},
	}

	for i := range exercisies {
		result := db.FirstOrCreate(&exercisies[i], MasterExercise{Name: exercisies[i].Name})
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}
