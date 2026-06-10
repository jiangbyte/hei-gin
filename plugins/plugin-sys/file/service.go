package file

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"hei-gin/sdk/config"
	"hei-gin/sdk/db"
	"hei-gin/sdk/result"
	"hei-gin/sdk/storage"
	"hei-gin/sdk/utils"

	"github.com/gin-gonic/gin"
)

var allowedExtensions = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true, ".svg": true, ".ico": true,
	".bmp": true, ".tiff": true,
	".doc": true, ".docx": true, ".xls": true, ".xlsx": true, ".ppt": true, ".pptx": true, ".pdf": true,
	".txt": true, ".csv": true, ".md": true,
	".zip": true, ".rar": true, ".7z": true, ".tar": true, ".gz": true,
	".mp3": true, ".wav": true, ".ogg": true,
	".mp4": true, ".avi": true, ".mkv": true, ".mov": true, ".webm": true,
	".json": true, ".xml": true, ".yaml": true, ".yml": true,
}

func isAllowedExtension(ext string) bool {
	return allowedExtensions[strings.ToLower(ext)]
}

func isImageExt(ext string) bool {
	switch strings.ToLower(ext) {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp", ".svg", ".ico", ".bmp", ".tiff":
		return true
	}
	return false
}

func formatFileSize(bytes int64) (kb int64, info string) {
	if bytes < 1024 {
		return 0, fmt.Sprintf("%d B", bytes)
	}
	kb = bytes / 1024
	if kb < 1024 {
		return kb, fmt.Sprintf("%d KB", kb)
	}
	mb := float64(kb) / 1024
	return kb, fmt.Sprintf("%.1f MB", mb)
}

func maxUploadSize() int64 {
	if config.C != nil && config.C.App.UploadMaxSize > 0 {
		return config.C.App.UploadMaxSize
	}
	return 50 << 20
}

// hashReader wraps an io.Reader to compute SHA256 on-the-fly during streaming.
type hashReader struct {
	reader io.Reader
	hash   hash.Hash
}

func (r *hashReader) Read(p []byte) (int, error) {
	n, err := r.reader.Read(p)
	if n > 0 {
		r.hash.Write(p[:n])
	}
	return n, err
}

func (r *hashReader) Sum() string {
	return hex.EncodeToString(r.hash.Sum(nil))
}

func newHashReader(reader io.Reader) *hashReader {
	return &hashReader{reader: reader, hash: sha256.New()}
}

// ===== CRUD =====

func Page(c *gin.Context, param *FilePageParam) {
	ctx := c.Request.Context()
	if param.Current < 1 {
		param.Current = 1
	}
	if param.Size < 1 {
		param.Size = 10
	}
	if param.Size > 100 {
		param.Size = 100
	}

	query := db.DB.WithContext(ctx).Model(&SysFile{})
	if param.Engine != "" {
		query = query.Where("engine = ?", param.Engine)
	}
	if param.Bucket != "" {
		query = query.Where("bucket = ?", param.Bucket)
	}
	if param.Keyword != "" {
		kw := "%" + param.Keyword + "%"
		query = query.Where("original_name LIKE ? OR original_name LIKE ?", kw, kw)
	}

	var total int64
	query.Count(&total)

	var records []SysFile
	query.Order("created_at DESC").Limit(param.Size).Offset((param.Current - 1) * param.Size).Find(&records)

	vos := make([]*FileVO, len(records))
	for i := range records {
		vos[i] = records[i].ToVO()
	}
	result.PageDataResult(c, vos, total, param.Current, param.Size)
}

func Detail(c *gin.Context, id string) *FileVO {
	if id == "" {
		return nil
	}
	var entity SysFile
	if err := db.DB.WithContext(c.Request.Context()).First(&entity, "id = ?", id).Error; err != nil {
		return nil
	}
	return entity.ToVO()
}

func Remove(c *gin.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	return db.DB.WithContext(c.Request.Context()).Where("id IN ?", ids).Delete(&SysFile{}).Error
}

func RemoveAbsolute(c *gin.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	ctx := c.Request.Context()
	var files []SysFile
	db.DB.WithContext(ctx).Where("id IN ?", ids).Find(&files)
	for _, f := range files {
		if f.Engine != "" {
			if eng := storage.GetStorage(f.Engine); eng != nil {
				eng.Delete(f.Bucket, f.FileKey)
			}
		}
	}
	return db.DB.WithContext(ctx).Where("id IN ?", ids).Delete(&SysFile{}).Error
}

// ===== Single-file Upload (streaming) =====

func Upload(c *gin.Context) (*FileUploadResult, error) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		return nil, fmt.Errorf("上传文件失败: %w", err)
	}
	defer file.Close()

	if header.Size > maxUploadSize() {
		return nil, fmt.Errorf("文件大小超过限制 (%d MB)", maxUploadSize()/(1<<20))
	}

	engineType := c.PostForm("engine")
	if engineType == "" {
		engineType = "LOCAL"
	}
	bucket := c.PostForm("bucket")
	if bucket == "" {
		bucket = "DEFAULT"
	}

	now := time.Now()
	ext := filepath.Ext(header.Filename)
	if !isAllowedExtension(ext) {
		return nil, fmt.Errorf("不支持的文件类型: %s", ext)
	}
	fileKey := utils.GenerateID() + ext

	eng := storage.GetStorage(engineType)
	if eng == nil {
		return nil, fmt.Errorf("不支持的存储类型: %s", engineType)
	}

	// Stream file to storage while computing SHA256 on-the-fly
	hr := newHashReader(file)
	storagePath, err := eng.StoreStream(bucket, fileKey, hr)
	if err != nil {
		return nil, fmt.Errorf("保存文件失败: %w", err)
	}

	checksum := hr.Sum()
	fileSizeKb, sizeInfo := formatFileSize(header.Size)
	downloadPath := storage.GetURL(engineType, bucket, fileKey)

	thumbnail := ""
	if isImageExt(ext) {
		thumbnail = downloadPath
	}

	entity := SysFile{
		ID:           utils.GenerateID(),
		Engine:       engineType,
		Bucket:       bucket,
		FileKey:      fileKey,
		ObjName:      fileKey,
		Name: header.Filename,
		Suffix:   ext,
		SizeKb:   fileSizeKb,
		SizeInfo:     sizeInfo,
		StoragePath:  storagePath,
		DownloadPath: downloadPath,
		Thumbnail:    thumbnail,
		Checksum:     checksum,
		ChecksumAlgo: "sha256",
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := db.DB.Create(&entity).Error; err != nil {
		return nil, fmt.Errorf("保存文件记录失败: %w", err)
	}

	return &FileUploadResult{
		ID:           entity.ID,
		Engine:       entity.Engine,
		Bucket:       entity.Bucket,
		FileKey:      entity.FileKey,
		Name: entity.Name,
		Suffix:   entity.Suffix,
		SizeKb:   entity.SizeKb,
		SizeInfo:     entity.SizeInfo,
		DownloadPath: entity.DownloadPath,
		Thumbnail:    entity.Thumbnail,
	}, nil
}

// ===== Chunked Upload =====

func InitChunkUpload(c *gin.Context, param *ChunkUploadInitParam) (*ChunkUploadResult, error) {
	engineType := param.Engine
	if engineType == "" {
		engineType = "LOCAL"
	}
	bucket := param.Bucket
	if bucket == "" {
		bucket = "DEFAULT"
	}

	eng := storage.GetStorage(engineType)
	if eng == nil {
		return nil, fmt.Errorf("不支持的存储类型: %s", engineType)
	}

	ext := filepath.Ext(param.FileName)
	fileKey := utils.GenerateID() + ext

	var uploadID string
	if cu, ok := eng.(storage.ChunkedUploader); ok {
		id, err := cu.InitChunkUpload(bucket, fileKey, param.TotalChunks)
		if err != nil {
			return nil, fmt.Errorf("初始化分片上传失败: %w", err)
		}
		uploadID = id
	} else {
		uploadID = utils.GenerateID()
		tmpDir := filepath.Join(os.TempDir(), "chunk_"+uploadID)
		if err := os.MkdirAll(tmpDir, 0755); err != nil {
			return nil, fmt.Errorf("创建临时目录失败: %w", err)
		}
	}

	chunkSize := param.FileSize / int64(param.TotalChunks)
	if param.FileSize%int64(param.TotalChunks) != 0 {
		chunkSize++
	}

	return &ChunkUploadResult{
		UploadID:    uploadID,
		FileKey:     fileKey,
		ChunkSize:   chunkSize,
		TotalChunks: param.TotalChunks,
	}, nil
}

func UploadChunk(c *gin.Context, param *ChunkUploadParam) error {
	engineType, _ := c.Get("_chunk_engine")
	engineStr, _ := engineType.(string)
	if engineStr == "" {
		engineStr = c.PostForm("engine")
		if engineStr == "" {
			engineStr = "LOCAL"
		}
	}
	bucket, _ := c.Get("_chunk_bucket")
	bucketStr, _ := bucket.(string)
	fileKey, _ := c.Get("_chunk_fileKey")
	fileKeyStr, _ := fileKey.(string)

	eng := storage.GetStorage(engineStr)
	if eng == nil {
		return fmt.Errorf("不支持的存储类型: %s", engineStr)
	}

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		return fmt.Errorf("读取分片文件失败: %w", err)
	}
	defer file.Close()

	if cu, ok := eng.(storage.ChunkedUploader); ok {
		chunk := storage.ChunkInfo{
			UploadID:    param.UploadID,
			ChunkIndex:  param.ChunkIndex,
			TotalChunks: param.TotalChunks,
			Checksum:    param.Checksum,
			Data:        file,
		}
		if err := cu.UploadChunk(bucketStr, fileKeyStr, param.UploadID, chunk); err != nil {
			return fmt.Errorf("上传分片失败: %w", err)
		}
	} else {
		tmpDir := filepath.Join(os.TempDir(), "chunk_"+param.UploadID)
		chunkFile := filepath.Join(tmpDir, fmt.Sprintf("chunk_%06d", param.ChunkIndex))
		data, err := io.ReadAll(file)
		if err != nil {
			return fmt.Errorf("读取分片数据失败: %w", err)
		}
		if err := os.WriteFile(chunkFile, data, 0644); err != nil {
			return fmt.Errorf("保存分片文件失败: %w", err)
		}
	}
	return nil
}

func CompleteChunkUpload(c *gin.Context, param *ChunkCompleteParam) (*FileUploadResult, error) {
	engineType := param.Engine
	if engineType == "" {
		engineType = "LOCAL"
	}
	bucket := param.Bucket
	if bucket == "" {
		bucket = "DEFAULT"
	}
	if param.FileKey == "" {
		return nil, fmt.Errorf("file_key 不能为空")
	}

	eng := storage.GetStorage(engineType)
	if eng == nil {
		return nil, fmt.Errorf("不支持的存储类型: %s", engineType)
	}

	now := time.Now()

	var storagePath string
	if cu, ok := eng.(storage.ChunkedUploader); ok {
		path, err := cu.CompleteChunkUpload(bucket, param.FileKey, param.UploadID)
		if err != nil {
			return nil, fmt.Errorf("合并分片失败: %w", err)
		}
		storagePath = path
	} else {
		tmpDir := filepath.Join(os.TempDir(), "chunk_"+param.UploadID)
		defer os.RemoveAll(tmpDir)

		path, err := mergeAndStore(eng, bucket, param.FileKey, tmpDir)
		if err != nil {
			return nil, err
		}
		storagePath = path
	}

	ext := filepath.Ext(param.Name)
	downloadPath := storage.GetURL(engineType, bucket, param.FileKey)

	thumbnail := ""
	if isImageExt(ext) {
		thumbnail = downloadPath
	}

	fileSizeKb, sizeInfo := formatFileSize(param.FileSize)

	entity := SysFile{
		ID:           utils.GenerateID(),
		Engine:       engineType,
		Bucket:       bucket,
		FileKey:      param.FileKey,
		ObjName:      param.FileKey,
		Name: param.Name,
		Suffix:   ext,
		SizeKb:   fileSizeKb,
		SizeInfo:     sizeInfo,
		StoragePath:  storagePath,
		DownloadPath: downloadPath,
		Thumbnail:    thumbnail,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := db.DB.Create(&entity).Error; err != nil {
		return nil, fmt.Errorf("保存文件记录失败: %w", err)
	}

	return &FileUploadResult{
		ID:           entity.ID,
		Engine:       entity.Engine,
		Bucket:       entity.Bucket,
		FileKey:      entity.FileKey,
		Name: entity.Name,
		Suffix:   entity.Suffix,
		SizeKb:   entity.SizeKb,
		SizeInfo:     entity.SizeInfo,
		DownloadPath: entity.DownloadPath,
		Thumbnail:    entity.Thumbnail,
	}, nil
}

func AbortChunkUpload(c *gin.Context, param *ChunkAbortParam) error {
	engineType := param.Engine
	if engineType == "" {
		engineType = "LOCAL"
	}

	eng := storage.GetStorage(engineType)
	if eng == nil {
		return fmt.Errorf("不支持的存储类型: %s", engineType)
	}

	if cu, ok := eng.(storage.ChunkedUploader); ok {
		return cu.AbortChunkUpload(param.Bucket, param.FileKey, param.UploadID)
	}

	tmpDir := filepath.Join(os.TempDir(), "chunk_"+param.UploadID)
	return os.RemoveAll(tmpDir)
}

func mergeAndStore(eng storage.Engine, bucket, fileKey, tmpDir string) (string, error) {
	entries, err := os.ReadDir(tmpDir)
	if err != nil {
		return "", fmt.Errorf("读取临时目录失败: %w", err)
	}

	pr, pw := io.Pipe()
	go func() {
		defer pw.Close()
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			chunkPath := filepath.Join(tmpDir, entry.Name())
			f, err := os.Open(chunkPath)
			if err != nil {
				pw.CloseWithError(fmt.Errorf("打开分片文件失败: %w", err))
				return
			}
			_, err = io.Copy(pw, f)
			f.Close()
			if err != nil {
				pw.CloseWithError(fmt.Errorf("读取分片数据失败: %w", err))
				return
			}
		}
	}()

	return eng.StoreStream(bucket, fileKey, pr)
}

// If storage provides a URL, use it; otherwise auto-construct from the request.

// ===== Download =====

func Download(c *gin.Context, id string) error {
	var entity SysFile
	if err := db.DB.WithContext(c.Request.Context()).First(&entity, "id = ?", id).Error; err != nil {
		return fmt.Errorf("文件不存在")
	}

	if entity.DownloadPath != "" {
		c.Redirect(302, entity.DownloadPath)
		return nil
	}

	if entity.StoragePath != "" {
		c.File(entity.StoragePath)
		return nil
	}

	return fmt.Errorf("文件路径为空")
}
