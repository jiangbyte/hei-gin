package log

import "hei-gin/sdk/utils"

// SysLogToLogVO 将 log.SysLog 映射到 log.LogVO
func SysLogToLogVO(src *SysLog) *LogVO {
	if src == nil {
		return nil
	}

	dst := &LogVO{}

	dst.ID = src.ID
	dst.Category = src.Category
	dst.Name = src.Name
	dst.ExeStatus = src.ExeStatus
	dst.ExeMessage = src.ExeMessage
	dst.OpIP = src.OpIP
	dst.OpAddress = src.OpAddress
	dst.OpBrowser = src.OpBrowser
	dst.OpOs = src.OpOs
	dst.ClassName = src.ClassName
	dst.MethodName = src.MethodName
	dst.ReqMethod = src.ReqMethod
	dst.ReqURL = src.ReqURL
	dst.ParamJSON = src.ParamJSON
	dst.ResultJSON = src.ResultJSON
	dst.TraceID = src.TraceID
	dst.OpUser = src.OpUser
	dst.SignData = src.SignData
	dst.CreatedBy = src.CreatedBy
	dst.UpdatedBy = src.UpdatedBy

	// *time.Time → string manual conversion
	dst.OpTime = utils.FormatDateTimePtr(src.OpTime)
	dst.CreatedAt = utils.FormatDateTimePtr(src.CreatedAt)
	dst.UpdatedAt = utils.FormatDateTimePtr(src.UpdatedAt)

	return dst
}

// LogVOToSysLog 将 log.LogVO 映射到 log.SysLog
func LogVOToSysLog(src *LogVO) *SysLog {
	if src == nil {
		return nil
	}

	dst := &SysLog{}

	dst.ID = src.ID
	dst.Category = src.Category
	dst.Name = src.Name
	dst.ExeStatus = src.ExeStatus
	dst.ExeMessage = src.ExeMessage
	dst.OpIP = src.OpIP
	dst.OpAddress = src.OpAddress
	dst.OpBrowser = src.OpBrowser
	dst.OpOs = src.OpOs
	dst.ClassName = src.ClassName
	dst.MethodName = src.MethodName
	dst.ReqMethod = src.ReqMethod
	dst.ReqURL = src.ReqURL
	dst.ParamJSON = src.ParamJSON
	dst.ResultJSON = src.ResultJSON
	dst.TraceID = src.TraceID
	dst.OpUser = src.OpUser
	dst.SignData = src.SignData
	dst.CreatedBy = src.CreatedBy
	dst.UpdatedBy = src.UpdatedBy

	// string → *time.Time manual conversion
	dst.OpTime = utils.ParseDateTimePtr(&src.OpTime)
	dst.CreatedAt = utils.ParseDateTimePtr(&src.CreatedAt)
	dst.UpdatedAt = utils.ParseDateTimePtr(&src.UpdatedAt)

	return dst
}
