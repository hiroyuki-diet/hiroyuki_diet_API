package model

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type UserSkin struct {
	Id        UUID               `gorm:"primary_key; type: uuid; not null; default:uuid_generate_v4()"`
	UserId    UUID               `gorm:"type: uuid; not null"`
	User      User               `gorm:"foreignKey:UserId;references:Id"`
	SkinId    UUID               `gorm:"type: uuid; not null"`
	Skin      MasterHiroyukiSkin `gorm:"foreignKey:SkinId;references:Id"`
	IsUsing   bool               `gorm:"type: bool; not null; default: false"`
	CreatedAt time.Time          `gorm:"type: timestamp; autoCreateTime; not null; default:CURRENT_TIMESTAMP;<-:create"`
	UpdatedAt time.Time          `gorm:"type: timestamp; autoUpdateTime;<-:update"`
	DeletedAt gorm.DeletedAt     `gorm:"type: timestamp; index"`
}

func (*UserSkin) Seeder(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var count int64

		// main.goが実行される度にレコードが生成されないようにする。
		tx.Model(&UserSkin{}).Count(&count)
		if count > 0 {
			return nil
		}

		var user *User
		if err := tx.First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("user not found")
			}
			return err
		}

		var skin MasterHiroyukiSkin
		if err := tx.First(&skin).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("skin not found")
			}
			return err
		}

		userSkin := UserSkin{UserId: user.Id, SkinId: skin.Id, IsUsing: true}

		if err := tx.Create(&userSkin).Error; err != nil {
			return err
		}

		return nil
	})
}
