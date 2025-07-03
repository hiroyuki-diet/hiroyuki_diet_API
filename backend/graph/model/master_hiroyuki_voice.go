package model

import (
	"time"

	"gorm.io/gorm"
)

type MasterHiroyukiVoice struct {
	Id           UUID           `gorm:"primary_key; type: uuid; not null; default:uuid_generate_v4()"`
	Name         string         `gorm:"type: varchar(50); not null"`
	VoiceUrl     string         `gorm:"type: varchar(50); not null"`
	ReleaseLevel int            `gorm:"type: int; not null"`
	VoiceFields  []MasterField  `gorm:"many2many:voice_fields"`
	CreatedAt    time.Time      `gorm:"type: timestamp; autoCreateTime; not null; default:CURRENT_TIMESTAMP;<-:create"`
	UpdatedAt    time.Time      `gorm:"type: timestamp; autoUpdateTime;<-:update"`
	DeletedAt    gorm.DeletedAt `gorm:"type: timestamp; index"`
}

func (*MasterHiroyukiVoice) FirstCreate(db *gorm.DB) error {
	var fields []MasterField
	fieldsStr := []string{"home", "chibi_hiroyuki"}

	err := db.Where("field IN ?", fieldsStr).Find(&fields).Error

	if err != nil {
		return err
	}

	voicies := []MasterHiroyukiVoice{
		{
			Name:         "よろしく",
			VoiceUrl:     "",
			ReleaseLevel: 0,
			VoiceFields:  fields,
		},
	}

	for i := range voicies {
		result := db.FirstOrCreate(&voicies[i], MasterHiroyukiVoice{Name: voicies[i].Name})
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}
