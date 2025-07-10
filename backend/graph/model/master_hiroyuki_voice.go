package model

import (
	"fmt"
	"time"

	"github.com/moXXcha/hiroyuki_diet_API/utils"
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

func (*MasterHiroyukiVoice) GetVoices(id UUID, fields []utils.Field, db *gorm.DB) ([]*HiroyukiVoiceResponse, error) {

	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}

	strFields := make([]string, len(fields))
	for i, f := range fields {
		strFields[i] = string(f)
	}

	var results []*HiroyukiVoiceResponse

	err := db.
		Table("master_hiroyuki_voices").
		Select(`
		master_hiroyuki_voices.id,
		master_hiroyuki_voices.name,
		master_hiroyuki_voices.voice_url,
		master_hiroyuki_voices.release_level,
		COALESCE(user_hiroyuki_voices.user_id IS NOT NULL, false) AS is_having
	`).
		Joins("JOIN voice_fields ON voice_fields.master_hiroyuki_voice_id = master_hiroyuki_voices.id").
		Joins("JOIN master_fields ON master_fields.id = voice_fields.master_field_id").
		Joins("LEFT JOIN user_hiroyuki_voices ON user_hiroyuki_voices.voice_id = master_hiroyuki_voices.id AND user_hiroyuki_voices.user_id = ?", id).
		Where("master_fields.field IN ?", strFields).
		Group("master_hiroyuki_voices.id, user_hiroyuki_voices.user_id").
		Preload("VoiceFields").
		Find(&results).Error

	if err != nil {
		return nil, err
	}

	for i := range results {
		var voice MasterHiroyukiVoice
		if err := db.Preload("VoiceFields").First(&voice, "id = ?", results[i].Id).Error; err == nil {
			fields := make([]utils.Field, 0, len(voice.VoiceFields))
			for _, f := range voice.VoiceFields {
				fields = append(fields, f.Field)
			}
			results[i].Fields = fields
		}
	}

	return results, err

}

func (*MasterHiroyukiVoice) FirstCreate(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var fields []MasterField
		fieldsStr := []string{"home", "chibi_hiroyuki"}

		if err := tx.Where("field IN ?", fieldsStr).Find(&fields).Error; err != nil {
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
			if err := tx.FirstOrCreate(&voicies[i], MasterHiroyukiVoice{Name: voicies[i].Name}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
