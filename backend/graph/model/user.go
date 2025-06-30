package model

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type User struct {
	Id                   UUID           `gorm:"primary_key; type: uuid; not null; default:uuid_generate_v4()"`
	Email                string         `gorm:"type: varchar(50); not null"`
	Password             string         `gorm:"type: text; not null"`
	Level                int            `gorm:"type: int; not null"`
	SignUpTokenId        UUID           `gorm:"type: uuid; not null"`
	SignUpToken          SignUpToken    `gorm:"foreignKey:SignUpTokenId;references:Id"`
	IsTokenAuthenticated bool           `gorm:"type: bool; not null; default: false"`
	ExperiencePoint      int            `gorm:"type: int; not null; default: 0"`
	CreatedAt            time.Time      `gorm:"type: timestamp; autoCreateTime; not null; default:CURRENT_TIMESTAMP;<-:create"`
	UpdatedAt            time.Time      `gorm:"type: timestamp; autoUpdateTime;<-:update"`
	DeletedAt            gorm.DeletedAt `gorm:"type: timestamp; index"`
}

func (*User) GetInfo(id UUID, db *gorm.DB) (*User, error) {
	var user *User

	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}

	result := db.Preload("SignUpToken").First(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}
