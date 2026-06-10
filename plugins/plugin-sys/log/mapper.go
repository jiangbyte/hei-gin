package log

import "hei-gin/sdk/utils"

// SysLogToLogVO maps SysLog -> LogVO
func SysLogToLogVO(src *SysLog) *LogVO {
	if src == nil {
		return nil
	}

	dst := &LogVO{}

	dst.ID = src.ID
	dst.Category = utils.PtrStrToStr(src.Category)
	dst.Name = utils.PtrStrToStr(src.Name)
	dst.ExeStatus = utils.PtrStrToStr(src.ExeStatus)
	dst.ExeMessage = utils.PtrStrToStr(src.ExeMessage)
	dst.OpIP = utils.PtrStrToStr(src.OpIP)
	dst.OpAddress = utils.PtrStrToStr(src.OpAddress)
	dst.OpBrowser = utils.PtrStrToStr(src.OpBrowser)
	dst.OpOs = utils.PtrStrToStr(src.OpOs)
	dst.ClassName = utils.PtrStrToStr(src.ClassName)
	dst.MethodName = utils.PtrStrToStr(src.MethodName)
	dst.ReqMethod = utils.PtrStrToStr(src.ReqMethod)
	dst.ReqURL = utils.PtrStrToStr(src.ReqURL)
	dst.ParamJSON = utils.PtrStrToStr(src.ParamJSON)
	dst.ResultJSON = utils.PtrStrToStr(src.ResultJSON)
	dst.TraceID = utils.PtrStrToStr(src.TraceID)
	dst.OpUser = utils.PtrStrToStr(src.OpUser)
	dst.SignData = utils.PtrStrToStr(src.SignData)
	dst.CreatedBy = utils.PtrStrToStr(src.CreatedBy)
	dst.UpdatedBy = utils.PtrStrToStr(src.UpdatedBy)

	dst.OpTime = utils.FormatDateTimePtr(src.OpTime)
	dst.CreatedAt = utils.FormatDateTimePtr(src.CreatedAt)
	dst.UpdatedAt = utils.FormatDateTimePtr(src.UpdatedAt)

	return dst
}

// LogVOToSysLog maps LogVO -> SysLog
func LogVOToSysLog(src *LogVO) *SysLog {
	if src == nil {
		return nil
	}

	dst := &SysLog{}

	dst.ID = src.ID
	dst.Category = utils.StrToStrPtr(src.Category)
	dst.Name = utils.StrToStrPtr(src.Name)
	dst.ExeStatus = utils.StrToStrPtr(src.ExeStatus)
	dst.ExeMessage = utils.StrToStrPtr(src.ExeMessage)
	dst.OpIP = utils.StrToStrPtr(src.OpIP)
	dst.OpAddress = utils.StrToStrPtr(src.OpAddress)
	dst.OpBrowser = utils.StrToStrPtr(src.OpBrowser)
	dst.OpOs = utils.StrToStrPtr(src.OpOs)
	dst.ClassName = utils.StrToStrPtr(src.ClassName)
	dst.MethodName = utils.StrToStrPtr(src.MethodName)
	dst.ReqMethod = utils.StrToStrPtr(src.ReqMethod)
	dst.ReqURL = utils.StrToStrPtr(src.ReqURL)
	dst.ParamJSON = utils.StrToStrPtr(src.ParamJSON)
	dst.ResultJSON = utils.StrToStrPtr(src.ResultJSON)
	dst.TraceID = utils.StrToStrPtr(src.TraceID)
	dst.OpUser = utils.StrToStrPtr(src.OpUser)
	dst.SignData = utils.StrToStrPtr(src.SignData)
	dst.CreatedBy = utils.StrToStrPtr(src.CreatedBy)
	dst.UpdatedBy = utils.StrToStrPtr(src.UpdatedBy)

	dst.OpTime = utils.ParseDateTimePtr(&src.OpTime)
	dst.CreatedAt = utils.ParseDateTimePtr(&src.CreatedAt)
	dst.UpdatedAt = utils.ParseDateTimePtr(&src.UpdatedAt)

	return dst
}
