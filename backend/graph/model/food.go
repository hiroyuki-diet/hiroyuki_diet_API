package model

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type Food struct {
	Id              UUID           `gorm:"primary_key; type: uuid; not null; default:uuid_generate_v4()"`
	Name            string         `gorm:"type: varchar(50); not null"`
	EstimateCalorie int            `gorm:"type: int; not null"`
	LastUsedDate    time.Time      `gorm:"type: date; not null"`
	CreatedAt       time.Time      `gorm:"type: timestamp; autoCreateTime; not null; default:CURRENT_TIMESTAMP;<-:create"`
	UpdatedAt       time.Time      `gorm:"type: timestamp; autoUpdateTime;<-:update"`
	DeletedAt       gorm.DeletedAt `gorm:"type: timestamp; index"`
}

func (*Food) GetAll(db *gorm.DB) ([]*Food, error) {
	var foods []*Food

	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}

	result := db.Find(&foods)

	if result.Error != nil {
		return nil, result.Error
	}

	return foods, nil
}

func (*Food) Seeder(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var count int64

		// main.goが実行される度にレコードが生成されないようにする。
		tx.Model(&Food{}).Count(&count)
		if count > 0 {
			return nil
		}

		// CSVファイルを開く
		file, err := os.Open("seeder/foods.csv")
		if err != nil {
			return fmt.Errorf("failed to open foods.csv: %w", err)
		}
		defer file.Close()

		reader := csv.NewReader(file)
		records, err := reader.ReadAll()
		if err != nil {
			return fmt.Errorf("failed to read foods.csv: %w", err)
		}

		defaultDate, err := time.Parse("2006-01-02", "2025-01-01")
		if err != nil {
			return err
		}

		// ヘッダー行をスキップしてデータを処理
		for i, record := range records {
			if i == 0 {
				continue // ヘッダー行をスキップ
			}

			if len(record) < 2 {
				continue
			}

			name := record[0]
			calorie, err := strconv.Atoi(record[1])
			if err != nil {
				return fmt.Errorf("failed to parse calorie for %s: %w", name, err)
			}

			food := Food{
				Name:            name,
				EstimateCalorie: calorie,
				LastUsedDate:    defaultDate,
			}

			if err := tx.Create(&food).Error; err != nil {
				return fmt.Errorf("failed to create food %s: %w", name, err)
			}
		}

		return nil
	})
}
