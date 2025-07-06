package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/moXXcha/hiroyuki_diet_API/utils"
	"gorm.io/gorm"
)

type Meal struct {
	Id           UUID           `gorm:"primary_key; type: uuid; not null; default:uuid_generate_v4()"`
	UserId       UUID           `gorm:"type: uuid; not null"`
	User         User           `gorm:"foreignKey:UserId;references:Id"`
	MealType     utils.MealType `gorm:"type: meal_type; not null"`
	TotalCalorie int            `gorm:"type: int; not null"`
	Foods        []Food         `gorm:"many2many:food_meals"`
	CreatedAt    time.Time      `gorm:"type: timestamp; autoCreateTime; not null; default:CURRENT_TIMESTAMP;<-:create"`
	UpdatedAt    time.Time      `gorm:"type: timestamp; autoUpdateTime;<-:update"`
	DeletedAt    gorm.DeletedAt `gorm:"type: timestamp; index"`
}

func (*Meal) GetAll(id UUID, db *gorm.DB) ([]*Meal, error) {
	var meals []*Meal

	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}

	err := db.Preload("Foods").Where("user_id = ?", id).Find(&meals).Error

	if err != nil {
		return nil, err
	}

	return meals, nil
}

func (*Meal) GetById(id UUID, db *gorm.DB) (*Meal, error) {
	var meal *Meal

	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}

	err := db.Preload("Foods").Where("id =?", id).First(&meal).Error

	if err != nil {
		return nil, err
	}

	return meal, nil
}

func (*Meal) Create(input InputMeal, db *gorm.DB) (*UUID, error) {
	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}

	var foods []Food

	err := db.Where("id IN ?", input.Foods).Find(&foods).Error

	if err != nil {
		return nil, err
	}

	totalCalorie := 0
	for _, food := range foods {
		totalCalorie += food.EstimateCalorie
	}

	mealInput := Meal{
		UserId:       *input.UserID,
		MealType:     input.MealType,
		TotalCalorie: totalCalorie,
		Foods:        foods,
	}

	err = db.Create(&mealInput).Error

	if err != nil {
		return nil, err
	}

	todayStr := time.Now().Format("2006-01-02")
	layout := "2006-01-02"

	t, err := time.Parse(layout, todayStr)
	if err != nil {
		return nil, err
	}

	err = db.Model(&Food{}).Where("id IN ?", input.Foods).Update("last_used_date", t).Error

	if err != nil {
		return nil, err
	}

	return &mealInput.Id, nil
}

func (*Meal) Seeder(db *gorm.DB) error {
	var count int64

	// main.goが実行される度にレコードが生成されないようにする。
	db.Model(&Meal{}).Count(&count)
	if count > 0 {
		return nil
	}

	var food []Food
	err := db.First(&food).Error

	if errors.Is(err, gorm.ErrRecordNotFound) || len(food) == 0 {
		return fmt.Errorf("food not found")
	}
	if err != nil {
		return err
	}

	var user User
	err = db.First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("user not found")
	}
	if err != nil {
		return err
	}

	meal := Meal{UserId: user.Id, MealType: "breakfast", TotalCalorie: 1000, Foods: food}
	err = db.Create(&meal).Error

	if err != err {
		return err
	}

	return nil
}
