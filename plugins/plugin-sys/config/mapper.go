package config

import "hei-gin/sdk/utils"

// SysConfigToConfigVO 将 config.SysConfig 映射到 config.ConfigVO
func SysConfigToConfigVO(src *SysConfig) *ConfigVO {
	if src == nil {
		return nil
	}

	dst := &ConfigVO{}

	dst.ID = src.ID
	dst.ConfigKey = src.ConfigKey
	dst.ConfigValue = src.ConfigValue
	dst.Category = src.Category
	dst.Remark = src.Remark
	dst.SortCode = src.SortCode
	dst.Extra = src.Extra
	dst.CreatedBy = src.CreatedBy
	dst.UpdatedBy = src.UpdatedBy

	// *time.Time → string manual conversion
	dst.CreatedAt = utils.FormatDateTimePtr(src.CreatedAt)
	dst.UpdatedAt = utils.FormatDateTimePtr(src.UpdatedAt)

	return dst
}

// ConfigVOToSysConfig 将 config.ConfigVO 映射到 config.SysConfig
func ConfigVOToSysConfig(src *ConfigVO) *SysConfig {
	if src == nil {
		return nil
	}

	dst := &SysConfig{}

	dst.ID = src.ID
	dst.ConfigKey = src.ConfigKey
	dst.ConfigValue = src.ConfigValue
	dst.Category = src.Category
	dst.Remark = src.Remark
	dst.SortCode = src.SortCode
	dst.Extra = src.Extra
	dst.CreatedBy = src.CreatedBy
	dst.UpdatedBy = src.UpdatedBy

	// string → *time.Time manual conversion
	dst.CreatedAt = utils.ParseDateTimePtr(&src.CreatedAt)
	dst.UpdatedAt = utils.ParseDateTimePtr(&src.UpdatedAt)

	return dst
}
