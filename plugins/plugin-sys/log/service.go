package log

import (
	"time"

	"gorm.io/gorm"

	"hei-gin/sdk/exception"
	"hei-gin/sdk/result"
	"hei-gin/sdk/utils"

	"github.com/gin-gonic/gin"
)

type service struct {
	repo *repository
}

func (s *service) LogPage(c *gin.Context, p *LogPageParam) {
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

	rows, total := s.repo.Page(ctx, p)

	vos := make([]*LogVO, len(rows))
	for i, r := range rows {
		vos[i] = SysLogToLogVO(&r)
	}
	result.PageDataResult(c, vos, total, p.Current, p.Size)
}

func (s *service) LogCreate(c *gin.Context, vo *LogVO) {
	ctx := c.Request.Context()

	e := LogVOToSysLog(vo)
	if err := s.repo.Create(ctx, e); err != nil {
		result.WriteError(c, exception.NewBusinessError("添加日志失败: "+err.Error(), 500))
		return
	}
}

func (s *service) LogModify(c *gin.Context, vo *LogVO) {
	ctx := c.Request.Context()
	if vo.ID == "" {
		result.WriteError(c, exception.NewBusinessError("ID不能为空", 400))
		return
	}

	_, err := s.repo.FindByID(ctx, vo.ID)
	if err != nil {
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
	if err := s.repo.UpdateByID(ctx, vo.ID, up); err != nil {
		result.WriteError(c, exception.NewBusinessError("编辑日志失败: "+err.Error(), 500))
		return
	}
}

func (s *service) LogRemove(c *gin.Context, param *utils.IdsParam) {
	ids := param.IDs
	if len(ids) == 0 {
		return
	}
	ctx := c.Request.Context()
	if err := s.repo.DeleteByIDs(ctx, ids); err != nil {
		result.WriteError(c, exception.NewBusinessError("删除日志失败: "+err.Error(), 500))
		return
	}
}

func (s *service) LogDetail(c *gin.Context, id string) *LogVO {
	if id == "" {
		result.WriteError(c, exception.NewBusinessError("ID不能为空", 400))
		return nil
	}
	ctx := c.Request.Context()
	e, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.WriteError(c, exception.NewBusinessError("数据不存在", 400))
			return nil
		}
		result.WriteError(c, exception.NewBusinessError("查询日志详情失败: "+err.Error(), 500))
		return nil
	}
	return SysLogToLogVO(e)
}

func (s *service) LogDeleteByCategory(c *gin.Context, param *LogDeleteByCategoryParam) {
	ctx := c.Request.Context()
	if err := s.repo.DeleteByCategory(ctx, param.Category); err != nil {
		result.WriteError(c, exception.NewBusinessError("按分类删除日志失败: "+err.Error(), 500))
		return
	}
}

func (s *service) LogLoginBarChart(c *gin.Context) *BarChartData {
	ctx := c.Request.Context()
	now := time.Now()
	since := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, -6)

	records := s.repo.ListByCategoriesSince(ctx, []string{"LOGIN", "LOGOUT"}, since)

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

func (s *service) LogLoginPieChart(c *gin.Context) *PieChartData {
	ctx := c.Request.Context()

	loginTotal := s.repo.CountByCategory(ctx, "LOGIN")
	logoutTotal := s.repo.CountByCategory(ctx, "LOGOUT")

	return &PieChartData{
		Data: []CategoryTotal{
			{Category: "登录", Total: int(loginTotal)},
			{Category: "登出", Total: int(logoutTotal)},
		},
	}
}

func (s *service) LogOpBarChart(c *gin.Context) *BarChartData {
	ctx := c.Request.Context()
	now := time.Now()
	since := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, -6)

	records := s.repo.ListByCategoriesSince(ctx, []string{"OPERATE", "EXCEPTION"}, since)

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

func (s *service) LogOpPieChart(c *gin.Context) *PieChartData {
	ctx := c.Request.Context()

	operateTotal := s.repo.CountByCategory(ctx, "OPERATE")
	exceptionTotal := s.repo.CountByCategory(ctx, "EXCEPTION")

	return &PieChartData{
		Data: []CategoryTotal{
			{Category: "操作", Total: int(operateTotal)},
			{Category: "异常", Total: int(exceptionTotal)},
		},
	}
}

func LogPage(c *gin.Context, p *LogPageParam) {
	defaultModule.service.LogPage(c, p)
}

func LogCreate(c *gin.Context, vo *LogVO) {
	defaultModule.service.LogCreate(c, vo)
}

func LogModify(c *gin.Context, vo *LogVO) {
	defaultModule.service.LogModify(c, vo)
}

func LogRemove(c *gin.Context, param *utils.IdsParam) {
	defaultModule.service.LogRemove(c, param)
}

func LogDetail(c *gin.Context, id string) *LogVO {
	return defaultModule.service.LogDetail(c, id)
}

func LogDeleteByCategory(c *gin.Context, param *LogDeleteByCategoryParam) {
	defaultModule.service.LogDeleteByCategory(c, param)
}

func LogLoginBarChart(c *gin.Context) *BarChartData {
	return defaultModule.service.LogLoginBarChart(c)
}

func LogLoginPieChart(c *gin.Context) *PieChartData {
	return defaultModule.service.LogLoginPieChart(c)
}

func LogOpBarChart(c *gin.Context) *BarChartData {
	return defaultModule.service.LogOpBarChart(c)
}

func LogOpPieChart(c *gin.Context) *PieChartData {
	return defaultModule.service.LogOpPieChart(c)
}
