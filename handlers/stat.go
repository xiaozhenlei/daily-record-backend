package handlers

import (
	"fmt"
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/user/daily-records-backend/models"
	"github.com/user/daily-records-backend/utils"
)

// GetWeekStat 获取周统计
func GetWeekStat(c *gin.Context) {
	userID := c.GetString("user_id")
	weekStart := c.Query("week_start") // 2026-02-16
	weekEnd := c.Query("week_end")     // 2026-02-22

	if weekStart == "" || weekEnd == "" {
		utils.ValidationError(c, "需提供 week_start 和 week_end")
		return
	}

	// 尝试从缓存获取
	cacheKey := utils.GenerateKey(userID, "week", weekStart+"_"+weekEnd)
	if cached := utils.GlobalCache.Get(cacheKey); cached != nil {
		utils.Success(c, cached)
		return
	}

	// 查询数据
	var records []models.Record
	_, err := utils.Client.From("daily_records").
		Select("*", "exact", false).
		Eq("user_id", userID).
		Gte("created_at", weekStart+" 00:00:00").
		Lte("created_at", weekEnd+" 23:59:59").
		ExecuteTo(&records)

	if err != nil {
		utils.Error(c, 500, "查询数据失败")
		return
	}

	// 聚合统计
	tagMap := make(map[string]*models.WeekStat)
	for _, r := range records {
		if _, ok := tagMap[r.Tag]; !ok {
			tagMap[r.Tag] = &models.WeekStat{Tag: r.Tag}
		}
		tagMap[r.Tag].Count++
		tagMap[r.Tag].TotalHours += float64(r.Duration) / 60.0
	}

	var stats []models.WeekStat
	for _, s := range tagMap {
		s.TotalHours = math.Round(s.TotalHours*100) / 100
		stats = append(stats, *s)
	}

	// 存入缓存
	utils.GlobalCache.Set(cacheKey, stats)
	utils.Success(c, stats)
}

// GetYearStat 获取年统计
func GetYearStat(c *gin.Context) {
	userID := c.GetString("user_id")
	yearStr := c.Query("year") // 2026
	if yearStr == "" {
		utils.ValidationError(c, "需提供 year")
		return
	}

	// 尝试从缓存获取
	cacheKey := utils.GenerateKey(userID, "year", yearStr)
	if cached := utils.GlobalCache.Get(cacheKey); cached != nil {
		utils.Success(c, cached)
		return
	}

	// 查询全年数据
	var records []models.Record
	_, err := utils.Client.From("daily_records").
		Select("*", "exact", false).
		Eq("user_id", userID).
		Like("created_at", yearStr+"-%").
		ExecuteTo(&records)

	if err != nil {
		utils.Error(c, 500, "查询全年数据失败")
		return
	}

	// 聚合逻辑 (含标签比例和每月分布)
	yearStat := models.YearStat{
		MonthHours: make([]models.MonthHour, 12),
	}
	for i := 0; i < 12; i++ {
		yearStat.MonthHours[i].Month = i + 1
	}

	tagMap := make(map[string]*models.YearTagStat)
	var totalHours float64 = 0

	for _, r := range records {
		// 标签聚合
		if _, ok := tagMap[r.Tag]; !ok {
			tagMap[r.Tag] = &models.YearTagStat{Tag: r.Tag}
		}
		hours := float64(r.Duration) / 60.0
		tagMap[r.Tag].Count++
		tagMap[r.Tag].TotalHours += hours
		totalHours += hours

		// 月份聚合 (解析 created_at 格式 "2026-02-21...")
		if len(r.CreatedAt) >= 7 {
			month, _ := strconv.Atoi(r.CreatedAt[5:7])
			if month >= 1 && month <= 12 {
				yearStat.MonthHours[month-1].TotalHours += hours
			}
		}
	}

	// 计算比例和极值
	maxH, minH := -1.0, 10000000.0
	for i, mh := range yearStat.MonthHours {
		mh.TotalHours = math.Round(mh.TotalHours*100) / 100
		yearStat.MonthHours[i] = mh
		if mh.TotalHours > maxH {
			maxH = mh.TotalHours
			yearStat.MaxMonth = mh.Month
			yearStat.MaxHours = mh.TotalHours
		}
		if mh.TotalHours < minH {
			minH = mh.TotalHours
			yearStat.MinMonth = mh.Month
			yearStat.MinHours = mh.TotalHours
		}
	}

	for _, ts := range tagMap {
		if totalHours > 0 {
			ts.Ratio = math.Round((ts.TotalHours/totalHours)*10000) / 100
		}
		ts.TotalHours = math.Round(ts.TotalHours*100) / 100
		yearStat.TagStats = append(yearStat.TagStats, *ts)
	}

	// 存入缓存
	utils.GlobalCache.Set(cacheKey, yearStat)
	utils.Success(c, yearStat)
}

// ExportWeek 导出周文本总结
func ExportWeek(c *gin.Context) {
	userID := c.GetString("user_id")
	weekStart := c.Query("week_start")
	weekEnd := c.Query("week_end")

	var records []models.Record
	utils.Client.From("daily_records").Select("*", "exact", false).Eq("user_id", userID).Gte("created_at", weekStart).Lte("created_at", weekEnd).ExecuteTo(&records)

	summary := fmt.Sprintf("📅 周总结 (%s ~ %s)\n\n", weekStart, weekEnd)
	total := 0
	for _, r := range records {
		summary += fmt.Sprintf("- [%s] %s (%d min)\n", r.Tag, r.Content, r.Duration)
		total += r.Duration
	}
	summary += fmt.Sprintf("\n总计用时: %.1f 小时", float64(total)/60.0)

	c.String(200, summary)
}

// ExportYear 导出年文本总结
func ExportYear(c *gin.Context) {
	userID := c.GetString("user_id")
	year := c.Query("year")

	var records []models.Record
	utils.Client.From("daily_records").Select("*", "exact", false).Eq("user_id", userID).Like("created_at", year+"-%").ExecuteTo(&records)

	summary := fmt.Sprintf("🏆 %s年度精进报告\n\n", year)
	tagTotal := make(map[string]int)
	for _, r := range records {
		tagTotal[r.Tag] += r.Duration
	}

	summary += "核心产出统计:\n"
	for tag, dur := range tagTotal {
		summary += fmt.Sprintf("- %s: %.1f 小时\n", tag, float64(dur)/60.0)
	}

	c.String(200, summary)
}
