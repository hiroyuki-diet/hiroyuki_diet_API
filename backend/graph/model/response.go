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
	Id      UUID
	Name    string
	IsClear bool
}

type HiroyukiVoiceResponse struct {
	Id           UUID
	VoiceUrl     string
	ReleaseLevel int
	Fields       []MasterField
	IsHaving     bool
}
