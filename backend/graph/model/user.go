package model

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/moXXcha/hiroyuki_diet_API/utils"
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

func (*User) Seeder(db *gorm.DB) error {
	var count int64

	// main.goが実行される度にレコードが生成されないようにする。
	db.Model(&User{}).Count(&count)
	if count > 0 {
		return nil
	}

	var signUpToken SignUpToken
	err := db.First(&signUpToken).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("signupToken not found")
	}

	if err != nil {
		return err
	}

	user := User{Email: "konami@example.com", Password: "test", Level: 1, SignUpTokenId: signUpToken.Id, IsTokenAuthenticated: true, ExperiencePoint: 50}

	err = db.Create(&user).Error

	if err != nil {
		return err
	}

	return nil
}

func (*User) SignUp(input Auth, db *gorm.DB) (*UUID, error) {
	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}

	tx := db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	var existingUser User
	if err := tx.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		tx.Rollback()
		return nil, fmt.Errorf("このメールアドレスは既に使用されています")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return nil, err
	}

	// トークンの生成
	token := rand.Intn(900000) + 100000
	signUpToken := SignUpToken{
		Token:       token,
		SurviveTime: 1,
	}

	if err := tx.Create(&signUpToken).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	user := User{
		Email:                input.Email,
		Password:             hashedPassword,
		Level:                0,
		ExperiencePoint:      0,
		SignUpTokenId:        signUpToken.Id,
		IsTokenAuthenticated: false,
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &user.Id, nil
}

func (*User) TokenAuth(input InputTokenAuth, db *gorm.DB) (*JWTTokenResponse, error) {
	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}

	tx := db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	// エラーが発生した場合にロールバックを確実に行う
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var user User
	if err := tx.Where("id = ?", input.UserID).First(&user).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("ユーザーが見つかりません")
		}
		return nil, err
	}

	if user.IsTokenAuthenticated {
		tx.Rollback()
		return nil, fmt.Errorf("すでに認証済みです")
	}

	var signUpToken SignUpToken
	if err := tx.Where("id = ?", user.SignUpTokenId).First(&signUpToken).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("トークンが見つかりません")
		}
		return nil, err
	}

	if signUpToken.Token != input.Token {
		tx.Rollback()
		return nil, fmt.Errorf("トークンが一致しません")
	}

	user.IsTokenAuthenticated = true
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// JWTトークンを生成
	token, err := utils.GenerateJWT(user.Id.String(), user.IsTokenAuthenticated, time.Hour*24)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to generate JWT: %w", err)
	}

	// すべての処理が成功したらコミット
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	tokenResponse := JWTTokenResponse{
		UserId: user.Id,
		Token:  token,
	}

	return &tokenResponse, nil
}
