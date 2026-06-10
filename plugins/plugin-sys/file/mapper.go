package file

import "hei-gin/sdk/utils"

// SysFileToFileVO 将 file.SysFile 映射到 file.FileVO
func SysFileToFileVO(src *SysFile) *FileVO {
	if src == nil {
		return nil
	}

	dst := &FileVO{}

	dst.ID = src.ID
	dst.Engine = src.Engine
	dst.Bucket = src.Bucket
	dst.FileKey = src.FileKey
	dst.Name = src.Name
	dst.Suffix = src.Suffix
	dst.SizeKb = src.SizeKb
	dst.SizeInfo = src.SizeInfo
	dst.ObjName = src.ObjName
	dst.StoragePath = src.StoragePath
	dst.DownloadPath = src.DownloadPath
	dst.IsDownloadAuth = src.IsDownloadAuth
	dst.Thumbnail = src.Thumbnail
	dst.Checksum = src.Checksum
	dst.ChecksumAlgo = src.ChecksumAlgo
	dst.ExtJson = src.ExtJson
	dst.CreatedBy = src.CreatedBy
	dst.UpdatedBy = src.UpdatedBy

	// time.Time → string manual conversion
	dst.CreatedAt = utils.FormatDateTime(src.CreatedAt)
	dst.UpdatedAt = utils.FormatDateTime(src.UpdatedAt)

	return dst
}

// FileVOToSysFile 将 file.FileVO 映射到 file.SysFile
func FileVOToSysFile(src *FileVO) *SysFile {
	if src == nil {
		return nil
	}

	dst := &SysFile{}

	dst.ID = src.ID
	dst.Engine = src.Engine
	dst.Bucket = src.Bucket
	dst.FileKey = src.FileKey
	dst.Name = src.Name
	dst.Suffix = src.Suffix
	dst.SizeKb = src.SizeKb
	dst.SizeInfo = src.SizeInfo
	dst.ObjName = src.ObjName
	dst.StoragePath = src.StoragePath
	dst.DownloadPath = src.DownloadPath
	dst.IsDownloadAuth = src.IsDownloadAuth
	dst.Thumbnail = src.Thumbnail
	dst.Checksum = src.Checksum
	dst.ChecksumAlgo = src.ChecksumAlgo
	dst.ExtJson = src.ExtJson
	dst.CreatedBy = src.CreatedBy
	dst.UpdatedBy = src.UpdatedBy

	return dst
}
