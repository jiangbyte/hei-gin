package log

import (
	"time"

	"gorm.io/gorm"

	"hei-gin/sdk/db"
	"hei-gin/sdk/exception"
	"hei-gin/sdk/result"
	"hei-gin/sdk/utils"

	"github.com/gin-gonic/gin"
)

func LogPage(c *gin.Context, p *LogPageParam) {
	ctx := c.Request.Context()
	if p.Current < 1 {
		p.Current = 1
	}
	if p.Size < 1 {
		p.Size = 10
	}
	if p.Size > 100 {
		p.Size = 100
	}

	q := db.DB.WithContext(ctx).Model(&SysLog{})
	if p.Category != "" {
		q = q.Where("category = ?", p.Category)
	}
	if p.ExeStatus != "" {
		q = q.Where("exe_status = ?", p.ExeStatus)
	}
	if p.Keyword != "" {
		like := "%" + p.Keyword + "%"
		q = q.Where("name LIKE ? OR op_user LIKE ? OR op_ip LIKE ?", like, like, like)
	}

	var total int64
	q.Count(&total)

	var rows []SysLog
	q.Order("created_at DESC").Limit(p.Size).Offset((p.Current - 1) * p.Size).Find(&rows)

	vos := make([]*LogVO, len(rows))
	for i, r := range rows {
		vos[i] = SysLogToLogVO(&r)
	}
	result.PageDataResult(c, vos, total, p.Current, p.Size)
}

func LogCreate(c *gin.Context, vo *LogVO) {
	ctx := c.Request.Context()

	e := LogVOToSysLog(vo)
	if err := db.DB.WithContext(ctx).Create(&e).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("添加日志失败: "+err.Error(), 500))
		return
	}
}

func LogModify(c *gin.Context, vo *LogVO) {
	ctx := c.Request.Context()
	if vo.ID == "" {
		result.WriteError(c, exception.NewBusinessError("ID不能为空", 400))
		return
	}

	var e SysLog
	if err := db.DB.WithContext(ctx).First(&e, "id = ?", vo.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			result.WriteError(c, exception.NewBusinessError("数据不存在", 400))
			return
		}
		result.WriteError(c, exception.NewBusinessError("查询日志失败: "+err.Error(), 500))
		return
	}

	up := map[string]interface{}{}
	if vo.Category != "" {
		up["category"] = vo.Category
	}
	if vo.Name != "" {
		up["name"] = vo.Name
	}
	if vo.ExeStatus != "" {
		up["exe_status"] = vo.ExeStatus
	}
	if vo.ExeMessage != "" {
		up["exe_message"] = vo.ExeMessage
	}
	if err := db.DB.WithContext(ctx).Model(&SysLog{}).Where("id = ?", vo.ID).Updates(up).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("编辑日志失败: "+err.Error(), 500))
		return
	}
}

func LogRemove(c *gin.Context, param *utils.IdsParam) {
	ids := param.IDs
	if len(ids) == 0 {
		return
	}
	ctx := c.Request.Context()
	if err := db.DB.WithContext(ctx).Where("id IN ?", ids).Delete(&SysLog{}).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("删除日志失败: "+err.Error(), 500))
		return
	}
}

func LogDetail(c *gin.Context, id string) *LogVO {
	if id == "" {
		result.WriteError(c, exception.NewBusinessError("ID不能为空", 400))
		return nil
	}
	ctx := c.Request.Context()
	var e SysLog
	if err := db.DB.WithContext(ctx).First(&e, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			result.WriteError(c, exception.NewBusinessError("数据不存在", 400))
			return nil
		}
		result.WriteError(c, exception.NewBusinessError("查询日志详情失败: "+err.Error(), 500))
		return nil
	}
	return SysLogToLogVO(&e)
}

func LogDeleteByCategory(c *gin.Context, param *LogDeleteByCategoryParam) {
	ctx := c.Request.Context()
	if err := db.DB.WithContext(ctx).Where("category = ?", param.Category).Delete(&SysLog{}).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("按分类删除日志失败: "+err.Error(), 500))
		return
	}
}

func LogLoginBarChart(c *gin.Context) *BarChartData {
	ctx := c.Request.Context()
	now := time.Now()
	since := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, -6)

	var records []SysLog
	db.DB.WithContext(ctx).Where("category IN ?", []string{"LOGIN", "LOGOUT"}).Where("op_time >= ?", since).Find(&records)

	days := make([]string, 7)
	for i := 0; i < 7; i++ {
		days[i] = since.AddDate(0, 0, i).Format("2006-01-02")
	}

	loginMap := make(map[string]int)
	logoutMap := make(map[string]int)
	for _, r := range records {
		if r.OpTime != nil && r.Category != "" {
			dayStr := r.OpTime.Format("2006-01-02")
			switch r.Category {
			case "LOGIN":
				loginMap[dayStr]++
			case "LOGOUT":
				logoutMap[dayStr]++
			}
		}
	}

	loginData := make([]int, 7)
	logoutData := make([]int, 7)
	for i, d := range days {
		loginData[i] = loginMap[d]
		logoutData[i] = logoutMap[d]
	}

	return &BarChartData{
		Days: days,
		Series: []CategorySeries{
			{Name: "登录", Data: loginData},
			{Name: "登出", Data: logoutData},
		},
	}
}

func LogLoginPieChart(c *gin.Context) *PieChartData {
	ctx := c.Request.Context()

	var loginTotal int64
	db.DB.WithContext(ctx).Model(&SysLog{}).Where("category = ?", "LOGIN").Count(&loginTotal)

	var logoutTotal int64
	db.DB.WithContext(ctx).Model(&SysLog{}).Where("category = ?", "LOGOUT").Count(&logoutTotal)

	return &PieChartData{
		Data: []CategoryTotal{
			{Category: "登录", Total: int(loginTotal)},
			{Category: "登出", Total: int(logoutTotal)},
		},
	}
}

func LogOpBarChart(c *gin.Context) *BarChartData {
	ctx := c.Request.Context()
	now := time.Now()
	since := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, -6)

	var records []SysLog
	db.DB.WithContext(ctx).Where("category IN ?", []string{"OPERATE", "EXCEPTION"}).Where("op_time >= ?", since).Find(&records)

	days := make([]string, 7)
	for i := 0; i < 7; i++ {
		days[i] = since.AddDate(0, 0, i).Format("2006-01-02")
	}

	operateMap := make(map[string]int)
	exceptionMap := make(map[string]int)
	for _, r := range records {
		if r.OpTime != nil && r.Category != "" {
			dayStr := r.OpTime.Format("2006-01-02")
			switch r.Category {
			case "OPERATE":
				operateMap[dayStr]++
			case "EXCEPTION":
				exceptionMap[dayStr]++
			}
		}
	}

	operateData := make([]int, 7)
	exceptionData := make([]int, 7)
	for i, d := range days {
		operateData[i] = operateMap[d]
		exceptionData[i] = exceptionMap[d]
	}

	return &BarChartData{
		Days: days,
		Series: []CategorySeries{
			{Name: "操作", Data: operateData},
			{Name: "异常", Data: exceptionData},
		},
	}
}

func LogOpPieChart(c *gin.Context) *PieChartData {
	ctx := c.Request.Context()

	var operateTotal int64
	db.DB.WithContext(ctx).Model(&SysLog{}).Where("category = ?", "OPERATE").Count(&operateTotal)

	var exceptionTotal int64
	db.DB.WithContext(ctx).Model(&SysLog{}).Where("category = ?", "EXCEPTION").Count(&exceptionTotal)

	return &PieChartData{
		Data: []CategoryTotal{
			{Category: "操作", Total: int(operateTotal)},
			{Category: "异常", Total: int(exceptionTotal)},
		},
	}
}
