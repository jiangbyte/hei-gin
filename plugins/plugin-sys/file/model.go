package file

import "time"

// SysFile follows the design from Snowy's DevFile entity.
type SysFile struct {
	ID            string    `gorm:"primaryKey;size:32" json:"id"`
	Engine        string    `gorm:"size:32;not null" json:"engine"`                  // LOCAL, MINIO, S3
	Bucket        string    `gorm:"size:128;not null" json:"bucket"`                 // storage bucket
	FileKey       string    `gorm:"size:500;not null;uniqueIndex" json:"file_key"`   // unique object key
	Name          string    `gorm:"size:255;not null" json:"name"`                   // original filename
	Suffix        string    `gorm:"size:32" json:"suffix"`                           // file extension
	SizeKb        int64     `gorm:"default:0" json:"size_kb"`                        // file size in KB
	SizeInfo      string    `gorm:"size:32" json:"size_info"`                        // formatted size
	ObjName       string    `gorm:"size:500" json:"obj_name"`                        // object name (same as FileKey)
	StoragePath   string    `gorm:"size:500" json:"storage_path"`                    // path in storage backend
	DownloadPath  string    `gorm:"size:500" json:"download_path"`                   // HTTP download URL
	IsDownloadAuth bool     `gorm:"default:false" json:"is_download_auth"`           // requires auth to download
	Thumbnail     string    `gorm:"size:500" json:"thumbnail"`                       // thumbnail URL for images
	Checksum      string    `gorm:"size:128" json:"checksum"`                        // SHA256 hex
	ChecksumAlgo  string    `gorm:"size:16" json:"checksum_algo"`                    // "sha256"
	ExtJson       string    `gorm:"type:text" json:"ext_json"`                       // extra metadata JSON
	CreatedAt     time.Time `json:"created_at"`
	CreatedBy     string    `gorm:"size:32" json:"created_by"`
	UpdatedAt     time.Time `json:"updated_at"`
	UpdatedBy     string    `gorm:"size:32" json:"updated_by"`
}

func (SysFile) TableName() string { return "sys_file" }
