package main

import (
	"fmt"

	"github.com/moXXcha/hiroyuki_diet_API/graph/model"
	"github.com/moXXcha/hiroyuki_diet_API/utils"
)

func main() {
	var err error
	db := utils.InitDB()

	exerciseModel := model.Exercise{}
	food := model.Food{}
	meal := model.Meal{}
	profile := model.Profile{}
	signupToken := model.SignUpToken{}
	userAchienvement := model.UserAchievement{}
	userVoice := model.UserHiroyukiVoice{}
	userItem := model.UserItem{}
	userSkin := model.UserSkin{}
	user := model.User{}

	err = signupToken.Seeder(db)
	if err != nil {
		panic(err)
	}
	err = user.Seeder(db)
	if err != nil {
		panic(err)
	}
	err = profile.Seeder(db)
	if err != nil {
		panic(err)
	}
	err = food.Seeder(db)
	if err != nil {
		panic(err)
	}
	err = meal.Seeder(db)
	if err != nil {
		panic(err)
	}
	err = exerciseModel.Seeder(db)
	if err != nil {
		panic(err)
	}
	err = userAchienvement.Seeder(db)
	if err != nil {
		panic(err)
	}
	err = userVoice.Seeder(db)
	if err != nil {
		panic(err)
	}
	err = userItem.Seeder(db)
	if err != nil {
		panic(err)
	}
	err = userSkin.Seeder(db)
	if err != nil {
		panic(err)
	}
	fmt.Println("db seeded")
}
