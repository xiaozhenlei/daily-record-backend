package handlers

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/user/daily-records-backend/models"
	"github.com/user/daily-records-backend/utils"
)

// GetYearlyStats 获取年度统计
func GetYearlyStats(c *gin.Context) {
	userID := c.GetString("user_id")
	yearStr := c.Query("year")
	if yearStr == "" {
		yearStr = strconv.Itoa(time.Now().Year())
	}

	// 尝试从缓存获取
	cacheKey := utils.GenerateKey(userID, "yearly_stats", yearStr)
	if cached := utils.GlobalCache.Get(cacheKey); cached != nil {
		utils.Success(c, cached)
		return
	}

	// 计算时间范围
	start := yearStr + "-01-01 00:00:00"
	end := yearStr + "-12-31 23:59:59"

	var records []models.Record
	_, err := utils.Client.From("daily_records").
		Select("*", "exact", false).
		Eq("user_id", userID).
		Gte("created_at", start).
		Lte("created_at", end).
		ExecuteTo(&records)

	if err != nil {
		utils.Error(c, 500, "获取年度数据失败")
		return
	}

	// 计算统计数据
	stats := models.YearlyStatsResponse{
		TagStats:     make([]models.YearlyTagStat, 0),
		MonthlyTrend: make([]models.MonthlyTrend, 12),
	}

	for i := 0; i < 12; i++ {
		stats.MonthlyTrend[i].Month = i + 1
	}

	tagMap := make(map[string]*models.YearlyTagStat)
	totalDuration := 0
	totalRecords := len(records)

	for _, r := range records {
		totalDuration += r.Duration

		// 标签统计
		if _, ok := tagMap[r.Tag]; !ok {
			tagMap[r.Tag] = &models.YearlyTagStat{Tag: r.Tag}
		}
		tagMap[r.Tag].Count++
		tagMap[r.Tag].Duration += r.Duration

		// 月度趋势
		if len(r.CreatedAt) >= 7 {
			month, _ := strconv.Atoi(r.CreatedAt[5:7])
			if month >= 1 && month <= 12 {
				stats.MonthlyTrend[month-1].Count++
				stats.MonthlyTrend[month-1].Duration += r.Duration
			}
		}
	}

	stats.TotalRecords = totalRecords
	stats.TotalDuration = totalDuration

	for _, s := range tagMap {
		if totalDuration > 0 {
			s.Ratio = float64(s.Duration) / float64(totalDuration)
		}
		stats.TagStats = append(stats.TagStats, *s)
	}

	// 存入缓存
	utils.GlobalCache.Set(cacheKey, stats)
	utils.Success(c, stats)
}

// GetMonthlyStats 获取月度统计
func GetMonthlyStats(c *gin.Context) {
	userID := c.GetString("user_id")
	yearStr := c.Query("year")
	monthStr := c.Query("month")

	if yearStr == "" || monthStr == "" {
		now := time.Now()
		yearStr = strconv.Itoa(now.Year())
		monthStr = strconv.Itoa(int(now.Month()))
	}

	// 补全月份两位
	if len(monthStr) == 1 {
		monthStr = "0" + monthStr
	}

	cacheKey := utils.GenerateKey(userID, "monthly_stats", yearStr+"-"+monthStr)
	if cached := utils.GlobalCache.Get(cacheKey); cached != nil {
		utils.Success(c, cached)
		return
	}

	// 简单的月度末尾逻辑 (或者直接用下个月减一秒)
	start := yearStr + "-" + monthStr + "-01 00:00:00"

	// 计算结束时间
	year, _ := strconv.Atoi(yearStr)
	month, _ := strconv.Atoi(monthStr)
	nextMonth := time.Date(year, time.Month(month+1), 1, 0, 0, 0, 0, time.UTC)
	end := nextMonth.Add(-time.Second).Format("2006-01-02 15:04:05")

	var records []models.Record
	_, err := utils.Client.From("daily_records").
		Select("*", "exact", false).
		Eq("user_id", userID).
		Gte("created_at", start).
		Lte("created_at", end).
		ExecuteTo(&records)

	if err != nil {
		utils.Error(c, 500, "获取月度数据失败")
		return
	}

	totalRecords := len(records)
	tagMap := make(map[string]*models.MonthlyTagStat)

	for _, r := range records {
		if _, ok := tagMap[r.Tag]; !ok {
			tagMap[r.Tag] = &models.MonthlyTagStat{Tag: r.Tag}
		}
		tagMap[r.Tag].Count++
		tagMap[r.Tag].Duration += r.Duration
	}

	stats := models.MonthlyStatsResponse{
		TagStats: make([]models.MonthlyTagStat, 0),
	}

	// 计算日均记录数 (在这个月已经过去的天数中)
	daysInMonth := nextMonth.Add(-time.Second).Day()
	stats.DailyAverage = float64(totalRecords) / float64(daysInMonth)

	for _, s := range tagMap {
		stats.TagStats = append(stats.TagStats, *s)
	}

	utils.GlobalCache.Set(cacheKey, stats)
	utils.Success(c, stats)
}
