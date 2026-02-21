package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/user/daily-records-backend/models"
	"github.com/user/daily-records-backend/utils"
)

// AddRecord 添加单条记录
func AddRecord(c *gin.Context) {
	var record models.Record
	// 参数绑定与校验
	if err := c.ShouldBindJSON(&record); err != nil {
		utils.ValidationError(c, "行动描述不能为空且长度不超过50字，运动时长需为正数")
		return
	}

	userID := c.GetString("user_id")
	record.UserID = userID
	record.Tag = models.ValidateTag(record.Tag) // 标签校验与修正

	// 插入 Supabase
	var result []models.Record
	_, err := utils.Client.From("daily_records").Insert(record, false, "", "", "").ExecuteTo(&result)
	if err != nil {
		utils.Error(c, 500, "保存记录失败: "+err.Error())
		return
	}

	utils.Success(c, result[0])
}

// BatchAddRecords 批量添加记录（离线同步）
func BatchAddRecords(c *gin.Context) {
	var body struct {
		Records []models.Record `json:"records" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		utils.ValidationError(c, "请求格式不正确")
		return
	}

	userID := c.GetString("user_id")
	successCount := 0
	var failedList []models.Record

	for _, req := range body.Records {
		req.UserID = userID
		req.Tag = models.ValidateTag(req.Tag)

		var res []models.Record
		_, err := utils.Client.From("daily_records").Insert(req, false, "", "", "").ExecuteTo(&res)
		if err != nil {
			failedList = append(failedList, req)
		} else {
			successCount++
		}
	}

	utils.Success(c, gin.H{
		"success_count": successCount,
		"failed_list":   failedList,
	})
}

// GetTodayRecords 获取今天的所有记录
func GetTodayRecords(c *gin.Context) {
	userID := c.GetString("user_id")
	today := time.Now().Format("2006-01-02")

	start := today + " 00:00:00"
	end := today + " 23:59:59"

	var records []models.Record
	_, err := utils.Client.From("daily_records").
		Select("*", "exact", false).
		Eq("user_id", userID).
		Gte("created_at", start).
		Lte("created_at", end).
		Order("created_at", &utils.OrderOptions{Ascending: false}).
		ExecuteTo(&records)

	if err != nil {
		utils.Error(c, 500, "获取今天记录失败")
		return
	}

	utils.Success(c, records)
}

// GetDateRecords 获取指定日期的记录
func GetDateRecords(c *gin.Context) {
	userID := c.GetString("user_id")
	dateStr := c.Query("date") // 格式: 2026-02-21
	if dateStr == "" {
		utils.ValidationError(c, "请指定日期")
		return
	}

	start := dateStr + " 00:00:00"
	end := dateStr + " 23:59:59"

	var records []models.Record
	_, err := utils.Client.From("daily_records").
		Select("*", "exact", false).
		Eq("user_id", userID).
		Gte("created_at", start).
		Lte("created_at", end).
		Order("created_at", &utils.OrderOptions{Ascending: false}).
		ExecuteTo(&records)

	if err != nil {
		utils.Error(c, 500, "获取日期记录失败")
		return
	}

	utils.Success(c, records)
}

// DeleteRecord 删除单条记录
func DeleteRecord(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("user_id")

	// 注意：Supabase 表开启了 RLS，正常情况下无需额外校验 user_id 匹配，
	// 但为了代码严谨，可通过 Eq 明确限制。
	_, err := utils.Client.From("daily_records").
		Delete("", "").
		Eq("id", id).
		Eq("user_id", userID).
		Execute()

	if err != nil {
		utils.Error(c, 500, "删除记录失败")
		return
	}

	utils.Success(c, "删除成功")
}
