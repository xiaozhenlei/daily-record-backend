package models

// YearlyStatsResponse 年度统计返回
type YearlyStatsResponse struct {
	TotalRecords  int             `json:"total_records"`
	TotalDuration int             `json:"total_duration"`
	TagStats      []YearlyTagStat `json:"tag_stats"`
	MonthlyTrend  []MonthlyTrend  `json:"monthly_trend"`
}

// YearlyTagStat 年度标签统计
type YearlyTagStat struct {
	Tag      string  `json:"tag"`
	Count    int     `json:"count"`
	Duration int     `json:"duration"`
	Ratio    float64 `json:"ratio"`
}

// MonthlyTrend 月度趋势
type MonthlyTrend struct {
	Month    int `json:"month"`
	Count    int `json:"count"`
	Duration int `json:"duration"`
}

// MonthlyStatsResponse 月度统计返回
type MonthlyStatsResponse struct {
	DailyAverage float64          `json:"daily_average"`
	TagStats     []MonthlyTagStat `json:"tag_stats"`
}

// MonthlyTagStat 月度标签统计
type MonthlyTagStat struct {
	Tag      string `json:"tag"`
	Count    int    `json:"count"`
	Duration int    `json:"duration"`
}
