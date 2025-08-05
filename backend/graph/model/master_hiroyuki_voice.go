package model

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
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
		// 既にデータが存在する場合は処理をスキップ
		var count int64
		tx.Model(&MasterHiroyukiVoice{}).Count(&count)
		if count > 0 {
			return nil
		}

		// CSVファイルを開く
		file, err := os.Open("seeder/master_voice.csv")
		if err != nil {
			// docker-composeからの実行パスを考慮
			file, err = os.Open("backend/seeder/master_voice.csv")
			if err != nil {
				return fmt.Errorf("failed to open master_voice.csv: %w", err)
			}
		}
		defer file.Close()

		reader := csv.NewReader(file)
		reader.Read() // ヘッダー行をスキップ

		// MasterFieldを事前にすべて取得しておく
		var masterFields []MasterField
		if err := tx.Find(&masterFields).Error; err != nil {
			return err
		}
		// Field名をキーにしたマップを作成
		fieldMap := make(map[string]MasterField)
		for _, mf := range masterFields {
			fieldMap[string(mf.Field)] = mf
		}

		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}

			releaseLevel, err := strconv.Atoi(record[2])
			if err != nil {
				return err
			}

			fieldName := record[3]
			field, ok := fieldMap[fieldName]
			if !ok {
				return fmt.Errorf("field not found: %s", fieldName)
			}

			voice := MasterHiroyukiVoice{
				Name:         record[0],
				VoiceUrl:     record[1],
				ReleaseLevel: releaseLevel,
				VoiceFields:  []MasterField{field}, // 関連付け
			}

			// 同じ名前のデータが存在しない場合のみ作成
			if err := tx.FirstOrCreate(&voice, MasterHiroyukiVoice{Name: voice.Name}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
