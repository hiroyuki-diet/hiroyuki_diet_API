package model

import "github.com/moXXcha/hiroyuki_diet_API/utils"

type ItemResponse struct {
	Id          UUID
	Name        string
	Description string
	ItemImage   string
	Count       int
}

type SkinResponse struct {
	Id           UUID
	Name         string
	Description  string
	Part         utils.SkinPart
	SkinImage    string
	ReleaseLevel int
	IsUsing      bool
	IsHaving     bool
}

type AchievementResponse struct {
	Id          UUID
	Name        string
	Description string
	IsClear     bool
}

type HiroyukiVoiceResponse struct {
	Id           UUID
	Name         string
	VoiceUrl     string
	ReleaseLevel int
	Fields       []utils.Field
	IsHaving     bool
}

type JWTTokenResponse struct {
	UserId UUID
	Token  string
}

type ExerciseMutationResponse struct {
	ID              *UUID
	NewAchievements []string
}

type LoginResponse struct {
	UserId          UUID
	Token           string
	NewAchievements []string
}
