package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/moXXcha/hiroyuki_diet_API/utils"
	"gorm.io/gorm"
)

type MasterHiroyukiSkin struct {
	Id           uuid.UUID      `gorm:"primary_key; type: uuid; not null; default:uuid_generate_v4()"`
	Name         string         `gorm:"type: varchar(50); not null"`
	Part         utils.SkinPart `gorm:"type: skin_part; not null"`
	SkinImage    string         `gorm:"type: varchar(50); not null"`
	Description  string         `gorm:"type: text; not null"`
	ReleaseLevel int            `gorm:"type: int; not null"`
	CreatedAt    time.Time      `gorm:"type: timestamp; autoCreateTime; not null; default:CURRENT_TIMESTAMP;<-:create"`
	UpdatedAt    time.Time      `gorm:"type: timestamp; autoUpdateTime;<-:update"`
	DeletedAt    gorm.DeletedAt `gorm:"type: timestamp; index"`
}

func (*MasterHiroyukiSkin) FirstCreate(db *gorm.DB) error {
	skins := []MasterHiroyukiSkin{
		{
			Name:         "鬼のツノ",
			Part:         "head",
			SkinImage:    "",
			Description:  "鬼のツノ。げきおこ",
			ReleaseLevel: 5,
		},
		{
			Name:         "まるメガネ",
			Part:         "face",
			SkinImage:    "",
			Description:  "かわいいまるメガネ。知的に見えるかも？",
			ReleaseLevel: 5,
		},
		{
			Name:         "論破Tシャツ",
			Part:         "body",
			SkinImage:    "",
			Description:  "なにやってんですか。運動してください",
			ReleaseLevel: 5,
		},
	}

	for i := range skins {
		result := db.FirstOrCreate(&skins[i], MasterHiroyukiSkin{Name: skins[i].Name})
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}
