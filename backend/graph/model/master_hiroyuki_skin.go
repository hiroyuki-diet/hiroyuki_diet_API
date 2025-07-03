package model

import (
	"fmt"
	"time"

	"github.com/moXXcha/hiroyuki_diet_API/utils"
	"gorm.io/gorm"
)

type MasterHiroyukiSkin struct {
	Id           UUID           `gorm:"primary_key; type: uuid; not null; default:uuid_generate_v4()"`
	Name         string         `gorm:"type: varchar(50); not null"`
	Part         utils.SkinPart `gorm:"type: skin_part; not null"`
	SkinImage    string         `gorm:"type: varchar(50); not null"`
	Description  string         `gorm:"type: text; not null"`
	ReleaseLevel int            `gorm:"type: int; not null"`
	CreatedAt    time.Time      `gorm:"type: timestamp; autoCreateTime; not null; default:CURRENT_TIMESTAMP;<-:create"`
	UpdatedAt    time.Time      `gorm:"type: timestamp; autoUpdateTime;<-:update"`
	DeletedAt    gorm.DeletedAt `gorm:"type: timestamp; index"`
}

func (*MasterHiroyukiSkin) GetSkins(id UUID, isUsingSkin bool, db *gorm.DB) ([]*SkinResponse, error) {
	var skins []*SkinResponse

	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}

	query := db.
		Table("master_hiroyuki_skins").
		Select(`
			master_hiroyuki_skins.id,
			master_hiroyuki_skins.name,
			master_hiroyuki_skins.description,
			master_hiroyuki_skins.part,
			master_hiroyuki_skins.release_level,
			COALESCE(user_skins.is_using, false) AS is_using,
			CASE WHEN user_skins.id IS NOT NULL THEN true ELSE false END AS is_having`).
		Joins("LEFT JOIN user_skins ON master_hiroyuki_skins.id = user_skins.skin_id AND user_skins.user_id = ?", id)

	if isUsingSkin {
		query = query.Where("user_skins.is_using = ?", true)
	}

	err := query.Scan(&skins).Error
	if err != nil {
		return nil, err
	}
	return skins, nil
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
