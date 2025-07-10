package model

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type MasterItem struct {
	Id          UUID           `gorm:"primary_key; type: uuid; not null; default:uuid_generate_v4()"`
	Name        string         `gorm:"type: varchar(50); not null"`
	Description string         `gorm:"type: text; not null"`
	ItemImage   string         `gorm:"type: varchar(50); not null"`
	CreatedAt   time.Time      `gorm:"type: timestamp; autoCreateTime; not null; default:CURRENT_TIMESTAMP;<-:create"`
	UpdatedAt   time.Time      `gorm:"type: timestamp; autoUpdateTime;<-:update"`
	DeletedAt   gorm.DeletedAt `gorm:"type: timestamp; index"`
}

func (*MasterItem) FirstCreate(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		items := []MasterItem{
			{
				Name:        "チートデイチケット",
				Description: "頑張ったごほうび！使ったら一日やすんでいいよ",
				ItemImage:   "",
			},
		}

		for i := range items {
			if err := tx.FirstOrCreate(&items[i], MasterItem{Name: items[i].Name}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (*MasterItem) Use(input InputUseItem, db *gorm.DB) (*UUID, error) {
	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}

	var userItemId UUID
	err := db.Transaction(func(tx *gorm.DB) error {
		var userItem UserItem
		if err := tx.Where("user_id = ? AND item_id = ?", input.UserID, input.ItemID).First(&userItem).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("user item not found")
			}
			return err
		}

		if userItem.Count < input.Count {
			return fmt.Errorf("not enough items")
		}

		userItem.Count -= input.Count
		if err := tx.Save(&userItem).Error; err != nil {
			return err
		}

		userItemId = userItem.Id
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &userItemId, nil
}

func (*MasterItem) GetAllByUserId(id UUID, db *gorm.DB) ([]*ItemResponse, error) {

	var responses []*ItemResponse

	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}

	err := db.Table("master_items").
		Select(`master_items.id as id, master_items.name, master_items.description, 
	        COALESCE(user_items.count, 0) as count`).
		Joins("LEFT JOIN user_items ON user_items.item_id = master_items.id AND user_items.user_id = ?", id).
		Scan(&responses).Error

	if err != nil {
		return nil, err
	}
	return responses, nil
}
