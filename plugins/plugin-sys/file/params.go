package file

// FileVO 文件视图对象
type FileVO struct {
	ID            string `json:"id"`
	Engine        string `json:"engine"`
	Bucket        string `json:"bucket"`
	FileKey       string `json:"file_key"`
	Name          string `json:"name"`
	Suffix        string `json:"suffix"`
	SizeKb        int64  `json:"size_kb"`
	SizeInfo      string `json:"size_info"`
	ObjName       string `json:"obj_name"`
	StoragePath   string `json:"storage_path"`
	DownloadPath  string `json:"download_path"`
	IsDownloadAuth bool `json:"is_download_auth"`
	Thumbnail     string `json:"thumbnail"`
	Checksum      string `json:"checksum"`
	ChecksumAlgo  string `json:"checksum_algo"`
	ExtJson       string `json:"ext_json"`
	CreatedAt     string `json:"created_at"`
	CreatedBy     string `json:"created_by"`
	UpdatedAt     string `json:"updated_at"`
	UpdatedBy     string `json:"updated_by"`
}

// FilePageParam 文件分页参数
type FilePageParam struct {
	Current int    `json:"current" form:"current"`
	Size    int    `json:"size" form:"size"`
	Keyword string `json:"keyword" form:"keyword"`
	Engine  string `json:"engine" form:"engine"`
	Bucket  string `json:"bucket" form:"bucket"`
}

// FileUploadResult 文件上传结果
type FileUploadResult struct {
	ID           string `json:"id"`
	Engine       string `json:"engine"`
	Bucket       string `json:"bucket"`
	FileKey      string `json:"file_key"`
	Name         string `json:"original_name"`
	Suffix       string `json:"file_suffix"`
	SizeKb       int64  `json:"file_size_kb"`
	SizeInfo     string `json:"size_info"`
	DownloadPath string `json:"download_path"`
	Thumbnail    string `json:"thumbnail"`
}

// ChunkUploadInitParam 分片上传初始化参数
type ChunkUploadInitParam struct {
	FileName    string `json:"file_name" form:"file_name" binding:"required"`
	FileSize    int64  `json:"file_size" form:"file_size" binding:"required"`
	TotalChunks int    `json:"total_chunks" form:"total_chunks" binding:"required"`
	Engine      string `json:"engine" form:"engine"`
	Bucket      string `json:"bucket" form:"bucket"`
}

// ChunkUploadResult 分片上传初始化结果
type ChunkUploadResult struct {
	UploadID    string `json:"upload_id"`
	FileKey     string `json:"file_key"`
	ChunkSize   int64  `json:"chunk_size"`
	TotalChunks int    `json:"total_chunks"`
}

// ChunkUploadParam 上传分片参数
type ChunkUploadParam struct {
	UploadID    string `json:"upload_id" form:"upload_id" binding:"required"`
	ChunkIndex  int    `json:"chunk_index" form:"chunk_index" binding:"required"`
	TotalChunks int    `json:"total_chunks" form:"total_chunks"`
	Checksum    string `json:"checksum" form:"checksum"`
}

// ChunkCompleteParam 完成分片上传参数
type ChunkCompleteParam struct {
	UploadID     string `json:"upload_id" form:"upload_id" binding:"required"`
	FileKey      string `json:"file_key" form:"file_key" binding:"required"`
	Name         string `json:"original_name" form:"original_name" binding:"required"`
	FileSize     int64  `json:"file_size" form:"file_size" binding:"required"`
	Engine       string `json:"engine" form:"engine"`
	Bucket       string `json:"bucket" form:"bucket"`
}

// ChunkAbortParam 取消分片上传参数
type ChunkAbortParam struct {
	UploadID string `json:"upload_id" form:"upload_id" binding:"required"`
	FileKey  string `json:"file_key" form:"file_key"`
	Engine   string `json:"engine" form:"engine"`
	Bucket   string `json:"bucket" form:"bucket"`
}
