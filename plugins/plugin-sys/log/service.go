package log

import (
	"time"

	"hei-gin/sdk/db"
	"hei-gin/sdk/exception"
	"hei-gin/sdk/result"
	"gorm.io/gorm"
	"hei-gin/sdk/utils"
	"github.com/gin-gonic/gin"
)

func Page(c *gin.Context, param *LogPageParam) {
	ctx := c.Request.Context()
	if param.Current < 1 { param.Current = 1 }
	if param.Size < 1 { param.Size = 10 }
	if param.Size > 100 { param.Size = 100 }

	query := db.DB.WithContext(ctx).Model(&SysLog{})
	if param.Category != "" { query = query.Where("category = ?", param.Category) }
	if param.Keyword != "" {
		kw := "%" + param.Keyword + "%"
		query = query.Where("name LIKE ? OR op_user LIKE ? OR op_ip LIKE ?", kw, kw, kw)
	}

	var total int64
	query.Count(&total)

	var records []SysLog
	query.Order("created_at DESC").Limit(param.Size).Offset((param.Current - 1) * param.Size).Find(&records)

	vos := make([]*LogVO, len(records))
	for i, r := range records { vos[i] = entToVO(&r) }
	result.PageDataResult(c, vos, total, param.Current, param.Size)
}

func LoginBarChart(c *gin.Context) *BarChartData {
	ctx := c.Request.Context()
	now := time.Now()
	since := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, -6)

	var records []SysLog
	db.DB.WithContext(ctx).Where("category IN ?", []string{"LOGIN", "LOGOUT"}).Where("op_time >= ?", since).Find(&records)

	days := make([]string, 7)
	for i := 0; i < 7; i++ { days[i] = since.AddDate(0, 0, i).Format("2006-01-02") }

	loginMap := make(map[string]int)
	logoutMap := make(map[string]int)
	for _, r := range records {
		if r.OpTime != nil && r.Category != nil {
			dayStr := r.OpTime.Format("2006-01-02")
			switch *r.Category {
			case "LOGIN": loginMap[dayStr]++
			case "LOGOUT": logoutMap[dayStr]++
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

func LoginPieChart(c *gin.Context) *PieChartData {
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

func OpBarChart(c *gin.Context) *BarChartData {
	ctx := c.Request.Context()
	now := time.Now()
	since := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, -6)

	var records []SysLog
	db.DB.WithContext(ctx).Where("category IN ?", []string{"OPERATE", "EXCEPTION"}).Where("op_time >= ?", since).Find(&records)

	days := make([]string, 7)
	for i := 0; i < 7; i++ { days[i] = since.AddDate(0, 0, i).Format("2006-01-02") }

	operateMap := make(map[string]int)
	exceptionMap := make(map[string]int)
	for _, r := range records {
		if r.OpTime != nil && r.Category != nil {
			dayStr := r.OpTime.Format("2006-01-02")
			switch *r.Category {
			case "OPERATE": operateMap[dayStr]++
			case "EXCEPTION": exceptionMap[dayStr]++
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

func OpPieChart(c *gin.Context) *PieChartData {
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

func entToVO(entity *SysLog) *LogVO {
	if entity == nil { return nil }
	vo := &LogVO{ID: entity.ID}
	if entity.Category != nil { vo.Category = *entity.Category }
	if entity.Name != nil { vo.Name = *entity.Name }
	if entity.ExeStatus != nil { vo.ExeStatus = *entity.ExeStatus }
	if entity.ExeMessage != nil { vo.ExeMessage = *entity.ExeMessage }
	if entity.OpIP != nil { vo.OpIP = *entity.OpIP }
	if entity.OpAddress != nil { vo.OpAddress = *entity.OpAddress }
	if entity.OpBrowser != nil { vo.OpBrowser = *entity.OpBrowser }
	if entity.OpOs != nil { vo.OpOs = *entity.OpOs }
	if entity.ClassName != nil { vo.ClassName = *entity.ClassName }
	if entity.MethodName != nil { vo.MethodName = *entity.MethodName }
	if entity.ReqMethod != nil { vo.ReqMethod = *entity.ReqMethod }
	if entity.ReqURL != nil { vo.ReqURL = *entity.ReqURL }
	if entity.ParamJSON != nil { vo.ParamJSON = *entity.ParamJSON }
	if entity.ResultJSON != nil { vo.ResultJSON = *entity.ResultJSON }
	if entity.TraceID != nil { vo.TraceID = *entity.TraceID }
	if entity.OpUser != nil { vo.OpUser = *entity.OpUser }
	if entity.SignData != nil { vo.SignData = *entity.SignData }
	if entity.CreatedBy != nil { vo.CreatedBy = *entity.CreatedBy }
	if entity.UpdatedBy != nil { vo.UpdatedBy = *entity.UpdatedBy }
	if entity.OpTime != nil { vo.OpTime = entity.OpTime.Format("2006-01-02 15:04:05") }
	if entity.CreatedAt != nil { vo.CreatedAt = entity.CreatedAt.Format("2006-01-02 15:04:05") }
	if entity.UpdatedAt != nil { vo.UpdatedAt = entity.UpdatedAt.Format("2006-01-02 15:04:05") }
	return vo
}


func Create(c *gin.Context, vo *LogVO, userID string) {
	ctx := c.Request.Context()
	now := time.Now()
	entity := SysLog{
		ID: utils.GenerateID(), CreatedAt: &now, UpdatedAt: &now,
	}
	if vo.Category != "" { entity.Category = &vo.Category }
	if vo.Name != "" { entity.Name = &vo.Name }
	if vo.ExeStatus != "" { entity.ExeStatus = &vo.ExeStatus }
	if vo.ExeMessage != "" { entity.ExeMessage = &vo.ExeMessage }
	if vo.OpIP != "" { entity.OpIP = &vo.OpIP }
	if vo.OpAddress != "" { entity.OpAddress = &vo.OpAddress }
	if userID != "" { entity.CreatedBy = &userID; entity.UpdatedBy = &userID }
	if err := db.DB.WithContext(ctx).Create(&entity).Error; err != nil {
		panic(exception.NewBusinessError("添加日志失败: "+err.Error(), 500))
	}
}

func Modify(c *gin.Context, vo *LogVO, userID string) {
	ctx := c.Request.Context()
	var entity SysLog
	if err := db.DB.WithContext(ctx).First(&entity, "id = ?", vo.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound { panic(exception.NewBusinessError("数据不存在", 400)) }
		panic(exception.NewBusinessError("查询日志失败: "+err.Error(), 500))
	}
	up := map[string]interface{}{"updated_at": time.Now()}
	if vo.Category != "" { up["category"] = vo.Category }
	if vo.Name != "" { up["name"] = vo.Name }
	if vo.ExeStatus != "" { up["exe_status"] = vo.ExeStatus }
	if userID != "" { up["updated_by"] = userID }
	if err := db.DB.WithContext(ctx).Model(&SysLog{}).Where("id = ?", vo.ID).Updates(up).Error; err != nil {
		panic(exception.NewBusinessError("编辑日志失败: "+err.Error(), 500))
	}
}

func Remove(c *gin.Context, ids []string) {
	if len(ids) == 0 { return }
	ctx := c.Request.Context()
	if err := db.DB.WithContext(ctx).Where("id IN ?", ids).Delete(&SysLog{}).Error; err != nil {
		panic(exception.NewBusinessError("删除日志失败: "+err.Error(), 500))
	}
}

func Detail(c *gin.Context, id string) *LogVO {
	if id == "" { return nil }
	ctx := c.Request.Context()
	var entity SysLog
	if err := db.DB.WithContext(ctx).First(&entity, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound { return nil }
		panic(exception.NewBusinessError("查询日志详情失败: "+err.Error(), 500))
	}
	return entToVO(&entity)
}

func DeleteByCategory(c *gin.Context, param *LogDeleteByCategoryParam) {
	ctx := c.Request.Context()
	if err := db.DB.WithContext(ctx).Where("category = ?", param.Category).Delete(&SysLog{}).Error; err != nil {
		panic(exception.NewBusinessError("按分类删除日志失败: "+err.Error(), 500))
	}
}

func VisLineChart(c *gin.Context) *BarChartData {
	return LoginBarChart(c)
}

func VisPieChart(c *gin.Context) *PieChartData {
	return LoginPieChart(c)
}

