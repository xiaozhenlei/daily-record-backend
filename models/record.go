package models

import (
	"strings"
)

// Record 每日行动记录结构体
type Record struct {
	ID         string `json:"id"`
	UserID     string `json:"user_id,omitempty"`
	Content    string `json:"content" binding:"required,max=50"`
	Tag        string `json:"tag" binding:"required"`
	Duration   int    `json:"duration" binding:"min=0"`
	CreatedAt  string `json:"created_at"`
}

// ValidateTag 验证标签并返回合法的标签
func ValidateTag(tag string) string {
	allowedTags := []string{"工作", "学习", "休闲", "家务", "其他"}
	for _, t := range allowedTags {
		if tag == t {
			return tag
		}
	}
	return "其他"
}

// WeekStat 周统计结构体
type WeekStat struct {
	Tag        string  `json:"tag"`
	Count      int     `json:"count"`
	TotalHours float64 `json:"total_hours"`
}

// YearTagStat 年标签统计
type YearTagStat struct {
	Tag        string  `json:"tag"`
	Count      int     `json:"count"`
	TotalHours float64 `json:"total_hours"`
	Ratio      float64 `json:"ratio"`
}

// MonthHour 每月耗时
type MonthHour struct {
	Month      int     `json:"month"`
	TotalHours float64 `json:"total_hours"`
}

// YearStat 年度统计结构体
type YearStat struct {
	TagStats   []YearTagStat `json:"tag_stats"`
	MonthHours []MonthHour   `json:"month_hours"`
	MaxMonth   int           `json:"max_month"`
	MaxHours   float64       `json:"max_hours"`
	MinMonth   int           `json:"min_month"`
	MinHours   float64       `json:"min_hours"`
}
