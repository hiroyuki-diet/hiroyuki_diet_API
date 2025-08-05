package utils

type Gender string
type MealType string
type SkinPart string
type Field string

const (
	Man   Gender = "man"
	Woman Gender = "woman"

	Breakfast MealType = "breakfast"
	Lunch     MealType = "lunch"
	Dinner    MealType = "dinner"
	Snacking  MealType = "snacking"

	Head SkinPart = "head"
	Face SkinPart = "face"
	Body SkinPart = "body"

	Login               Field = "login"
	Signin              Field = "signin"
	Home                Field = "home"
	Meal                Field = "meal"
	MealForm            Field = "meal_form"
	MealEdit            Field = "meal_edit"
	Data                Field = "data"
	Profile             Field = "profile"
	ProfileEdit         Field = "profile_edit"
	Exercise            Field = "exercise"
	ExerciseComplete    Field = "exercise_complete"
	Achievement         Field = "achievement"
	AchievementComplete Field = "achievement_complete"
	ChibiHiroyuki       Field = "chibi_hiroyuki"
)
