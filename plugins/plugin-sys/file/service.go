package file

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"hei-gin/sdk/config"
	"hei-gin/sdk/db"
	"hei-gin/sdk/exception"
	"hei-gin/sdk/result"
	"hei-gin/sdk/storage"
	"hei-gin/sdk/utils"

	"github.com/gin-gonic/gin"
)

var allowedExtensions = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true, ".ico": true,
	".bmp": true, ".tiff": true,
	".doc": true, ".docx": true, ".xls": true, ".xlsx": true, ".ppt": true, ".pptx": true, ".pdf": true,
	".txt": true, ".csv": true, ".md": true,
	".zip": true, ".rar": true, ".7z": true, ".tar": true, ".gz": true,
	".mp3": true, ".wav": true, ".ogg": true,
	".mp4": true, ".avi": true, ".mkv": true, ".mov": true, ".webm": true,
	".json": true, ".xml": true, ".yaml": true, ".yml": true,
}

const chunkSize int64 = 5 << 20

type chunkUploadState struct {
	Engine      string `json:"engine"`
	Bucket      string `json:"bucket"`
	FileKey     string `json:"file_key"`
	Name        string `json:"name"`
	FileSize    int64  `json:"file_size"`
	TotalChunks int    `json:"total_chunks"`
	OwnerID     string `json:"owner_id"`
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

func validateUploadMeta(fileName string, fileSize int64) (string, error) {
	if strings.TrimSpace(fileName) == "" {
		return "", fmt.Errorf("文件名不能为空")
	}
	if fileSize <= 0 {
		return "", fmt.Errorf("文件大小必须大于0")
	}
	if fileSize > maxUploadSize() {
		return "", fmt.Errorf("文件大小超过限制 (%d MB)", maxUploadSize()/(1<<20))
	}
	ext := strings.ToLower(filepath.Ext(fileName))
	if ext == "" || !isAllowedExtension(ext) {
		return "", fmt.Errorf("不支持的文件类型: %s", ext)
	}
	return ext, nil
}

func chunkStateKey(uploadID string) string {
	return "hei:file:chunk:" + uploadID
}

func currentOwnerID(c *gin.Context) string {
	if v, exists := c.Get("login_id"); exists {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func saveChunkState(c *gin.Context, uploadID string, state chunkUploadState) error {
	if db.Redis == nil {
		return fmt.Errorf("Redis 不可用，无法初始化分片上传")
	}
	data, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("保存分片上传状态失败: %w", err)
	}
	if err := db.Redis.SetEx(c.Request.Context(), chunkStateKey(uploadID), data, 24*time.Hour).Err(); err != nil {
		return fmt.Errorf("保存分片上传状态失败: %w", err)
	}
	return nil
}

func loadChunkState(c *gin.Context, uploadID string) (*chunkUploadState, error) {
	if strings.TrimSpace(uploadID) == "" {
		return nil, fmt.Errorf("upload_id 不能为空")
	}
	if db.Redis == nil {
		return nil, fmt.Errorf("Redis 不可用，无法读取分片上传状态")
	}
	data, err := db.Redis.Get(c.Request.Context(), chunkStateKey(uploadID)).Bytes()
	if err != nil {
		return nil, fmt.Errorf("分片上传会话不存在或已过期")
	}
	var state chunkUploadState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("分片上传状态无效")
	}
	if state.FileKey == "" || state.TotalChunks <= 0 {
		return nil, fmt.Errorf("分片上传状态无效")
	}
	ownerID := currentOwnerID(c)
	if state.OwnerID != "" && ownerID != "" && state.OwnerID != ownerID {
		return nil, fmt.Errorf("无权访问该分片上传会话")
	}
	return &state, nil
}

func deleteChunkState(c *gin.Context, uploadID string) {
	if db.Redis != nil && uploadID != "" {
		_ = db.Redis.Del(c.Request.Context(), chunkStateKey(uploadID)).Err()
	}
}

func storeStream(ctx context.Context, eng storage.Engine, bucket, fileKey string, reader io.Reader, size int64) (string, error) {
	if ss, ok := eng.(storage.ContextSizedStreamer); ok {
		return ss.StoreStreamWithContext(ctx, bucket, fileKey, reader, size)
	}
	if ss, ok := eng.(storage.SizedStreamer); ok {
		return ss.StoreStreamWithSize(bucket, fileKey, reader, size)
	}
	return eng.StoreStream(bucket, fileKey, reader)
}

func supportsChunkUpload(eng storage.Engine) bool {
	if _, ok := eng.(storage.ContextChunkedUploader); ok {
		return true
	}
	_, ok := eng.(storage.ChunkedUploader)
	return ok
}

func initChunkUpload(ctx context.Context, eng storage.Engine, bucket, fileKey string, totalChunks int) (string, error) {
	if cu, ok := eng.(storage.ContextChunkedUploader); ok {
		return cu.InitChunkUploadWithContext(ctx, bucket, fileKey, totalChunks)
	}
	return eng.(storage.ChunkedUploader).InitChunkUpload(bucket, fileKey, totalChunks)
}

func uploadChunk(ctx context.Context, eng storage.Engine, bucket, fileKey, uploadID string, chunk storage.ChunkInfo) error {
	if cu, ok := eng.(storage.ContextChunkedUploader); ok {
		return cu.UploadChunkWithContext(ctx, bucket, fileKey, uploadID, chunk)
	}
	return eng.(storage.ChunkedUploader).UploadChunk(bucket, fileKey, uploadID, chunk)
}

func completeChunkUpload(ctx context.Context, eng storage.Engine, bucket, fileKey, uploadID string) (string, error) {
	if cu, ok := eng.(storage.ContextChunkedUploader); ok {
		return cu.CompleteChunkUploadWithContext(ctx, bucket, fileKey, uploadID)
	}
	return eng.(storage.ChunkedUploader).CompleteChunkUpload(bucket, fileKey, uploadID)
}

func abortChunkUpload(ctx context.Context, eng storage.Engine, bucket, fileKey, uploadID string) error {
	if cu, ok := eng.(storage.ContextChunkedUploader); ok {
		return cu.AbortChunkUploadWithContext(ctx, bucket, fileKey, uploadID)
	}
	return eng.(storage.ChunkedUploader).AbortChunkUpload(bucket, fileKey, uploadID)
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

type contextReader struct {
	ctx    context.Context
	reader io.Reader
}

func (r *contextReader) Read(p []byte) (int, error) {
	select {
	case <-r.ctx.Done():
		return 0, r.ctx.Err()
	default:
		return r.reader.Read(p)
	}
}

// ===== Page =====

func FilePage(c *gin.Context, p *FilePageParam) {
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

	q := db.DB.WithContext(ctx).Model(&SysFile{})
	if p.Engine != "" {
		q = q.Where("engine = ?", p.Engine)
	}
	if p.Bucket != "" {
		q = q.Where("bucket = ?", p.Bucket)
	}
	if p.Keyword != "" {
		kw := "%" + p.Keyword + "%"
		q = q.Where("name LIKE ? OR name LIKE ?", kw, kw)
	}

	var total int64
	q.Count(&total)

	var rows []SysFile
	q.Order("created_at DESC").Limit(p.Size).Offset((p.Current - 1) * p.Size).Find(&rows)

	vos := make([]*FileVO, len(rows))
	for i, r := range rows {
		vos[i] = SysFileToFileVO(&r)
	}
	result.PageDataResult(c, vos, total, p.Current, p.Size)
}

// ===== Detail =====

func FileDetail(c *gin.Context, id string) *FileVO {
	if id == "" {
		result.WriteError(c, exception.NewBusinessError("ID不能为空", 400))
		return nil
	}
	ctx := c.Request.Context()
	var e SysFile
	if err := db.DB.WithContext(ctx).First(&e, "id = ?", id).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("文件不存在", 400))
		return nil
	}
	return SysFileToFileVO(&e)
}

// ===== Remove =====

func FileRemove(c *gin.Context, param *utils.IdsParam) {
	ids := param.IDs
	if len(ids) == 0 {
		return
	}
	if err := db.DB.WithContext(c.Request.Context()).Where("id IN ?", ids).Delete(&SysFile{}).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("删除文件失败: "+err.Error(), 500))
		return
	}
}

func FileRemoveAbsolute(c *gin.Context, param *utils.IdsParam) {
	ids := param.IDs
	if len(ids) == 0 {
		return
	}
	ctx := c.Request.Context()
	var files []SysFile
	db.DB.WithContext(ctx).Where("id IN ?", ids).Find(&files)
	for _, f := range files {
		if f.Engine != "" {
			if eng := storage.GetStorage(f.Engine); eng != nil {
				if deleter, ok := eng.(storage.ContextDeleter); ok {
					_ = deleter.DeleteWithContext(ctx, f.Bucket, f.FileKey)
				} else {
					_ = eng.Delete(f.Bucket, f.FileKey)
				}
			}
		}
	}
	if err := db.DB.WithContext(ctx).Where("id IN ?", ids).Delete(&SysFile{}).Error; err != nil {
		result.WriteError(c, exception.NewBusinessError("删除文件失败: "+err.Error(), 500))
		return
	}
}

// ===== Single-file Upload (streaming) =====

func FileUpload(c *gin.Context) (*FileUploadResult, error) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		return nil, fmt.Errorf("上传文件失败: %w", err)
	}
	defer file.Close()

	ext, err := validateUploadMeta(header.Filename, header.Size)
	if err != nil {
		return nil, err
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
	fileKey := utils.GenerateID() + ext

	eng := storage.GetStorage(engineType)
	if eng == nil {
		return nil, fmt.Errorf("不支持的存储类型: %s", engineType)
	}

	// Stream file to storage while computing SHA256 on-the-fly
	hr := newHashReader(file)
	storagePath, err := storeStream(c.Request.Context(), eng, bucket, fileKey, hr, header.Size)
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

		Engine:       engineType,
		Bucket:       bucket,
		FileKey:      fileKey,
		ObjName:      fileKey,
		Name:         header.Filename,
		Suffix:       ext,
		SizeKb:       fileSizeKb,
		SizeInfo:     sizeInfo,
		StoragePath:  storagePath,
		DownloadPath: downloadPath,
		Thumbnail:    thumbnail,
		Checksum:     checksum,
		ChecksumAlgo: "sha256",
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := db.DB.WithContext(c.Request.Context()).Create(&entity).Error; err != nil {
		return nil, fmt.Errorf("保存文件记录失败: %w", err)
	}

	return &FileUploadResult{
		ID:           entity.ID,
		Engine:       entity.Engine,
		Bucket:       entity.Bucket,
		FileKey:      entity.FileKey,
		Name:         entity.Name,
		Suffix:       entity.Suffix,
		SizeKb:       entity.SizeKb,
		SizeInfo:     entity.SizeInfo,
		DownloadPath: entity.DownloadPath,
		Thumbnail:    entity.Thumbnail,
	}, nil
}

// ===== Chunk Upload =====

func FileInitChunkUpload(c *gin.Context, param *ChunkUploadInitParam) (*ChunkUploadResult, error) {
	ext, err := validateUploadMeta(param.FileName, param.FileSize)
	if err != nil {
		return nil, err
	}
	if param.TotalChunks <= 0 {
		return nil, fmt.Errorf("total_chunks 必须大于0")
	}
	expectedChunks := int((param.FileSize + chunkSize - 1) / chunkSize)
	if param.TotalChunks != expectedChunks {
		return nil, fmt.Errorf("total_chunks 与文件大小不匹配")
	}

	engineType := param.Engine
	if engineType == "" {
		engineType = "LOCAL"
	}
	bucket := param.Bucket
	if bucket == "" {
		bucket = "DEFAULT"
	}

	fileKey := utils.GenerateID() + ext

	eng := storage.GetStorage(engineType)
	if eng == nil {
		return nil, fmt.Errorf("不支持的存储类型: %s", engineType)
	}

	if supportsChunkUpload(eng) {
		ctx := c.Request.Context()
		uploadID, err := initChunkUpload(ctx, eng, bucket, fileKey, param.TotalChunks)
		if err != nil {
			return nil, fmt.Errorf("分片上传初始化失败: %w", err)
		}
		if err := saveChunkState(c, uploadID, chunkUploadState{
			Engine:      engineType,
			Bucket:      bucket,
			FileKey:     fileKey,
			Name:        param.FileName,
			FileSize:    param.FileSize,
			TotalChunks: param.TotalChunks,
			OwnerID:     currentOwnerID(c),
		}); err != nil {
			_ = abortChunkUpload(ctx, eng, bucket, fileKey, uploadID)
			return nil, err
		}
		return &ChunkUploadResult{
			UploadID:    uploadID,
			FileKey:     fileKey,
			ChunkSize:   chunkSize,
			TotalChunks: param.TotalChunks,
		}, nil
	}

	// Fallback: create temp directory for local chunk storage
	tmpDir := filepath.Join(os.TempDir(), "chunk_"+fileKey)
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return nil, fmt.Errorf("创建临时目录失败: %w", err)
	}
	if err := saveChunkState(c, fileKey, chunkUploadState{
		Engine:      engineType,
		Bucket:      bucket,
		FileKey:     fileKey,
		Name:        param.FileName,
		FileSize:    param.FileSize,
		TotalChunks: param.TotalChunks,
		OwnerID:     currentOwnerID(c),
	}); err != nil {
		_ = os.RemoveAll(tmpDir)
		return nil, err
	}

	return &ChunkUploadResult{
		UploadID:    fileKey,
		FileKey:     fileKey,
		ChunkSize:   chunkSize,
		TotalChunks: param.TotalChunks,
	}, nil
}

func FileUploadChunk(c *gin.Context, param *ChunkUploadParam) error {
	state, err := loadChunkState(c, param.UploadID)
	if err != nil {
		return err
	}
	if param.TotalChunks != 0 && param.TotalChunks != state.TotalChunks {
		return fmt.Errorf("total_chunks 与初始化信息不一致")
	}
	if param.ChunkIndex < 0 || param.ChunkIndex >= state.TotalChunks {
		return fmt.Errorf("chunk_index 超出范围")
	}

	eng := storage.GetStorage(state.Engine)
	if eng == nil {
		return fmt.Errorf("不支持的存储类型: %s", state.Engine)
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		return fmt.Errorf("读取分片文件失败: %w", err)
	}
	defer file.Close()
	if header.Size <= 0 || header.Size > chunkSize {
		return fmt.Errorf("分片大小超过限制")
	}
	if param.ChunkIndex == state.TotalChunks-1 {
		expectedLast := state.FileSize - int64(state.TotalChunks-1)*chunkSize
		if header.Size != expectedLast {
			return fmt.Errorf("最后一个分片大小不匹配")
		}
	} else if header.Size != chunkSize {
		return fmt.Errorf("分片大小不匹配")
	}

	if supportsChunkUpload(eng) {
		chunk := storage.ChunkInfo{
			UploadID:    param.UploadID,
			ChunkIndex:  param.ChunkIndex,
			TotalChunks: state.TotalChunks,
			Checksum:    param.Checksum,
			Size:        header.Size,
			Data:        &contextReader{ctx: c.Request.Context(), reader: file},
		}
		if err := uploadChunk(c.Request.Context(), eng, state.Bucket, state.FileKey, param.UploadID, chunk); err != nil {
			return fmt.Errorf("上传分片失败: %w", err)
		}
	} else {
		tmpDir := filepath.Join(os.TempDir(), "chunk_"+param.UploadID)
		if err := os.MkdirAll(tmpDir, 0755); err != nil {
			return fmt.Errorf("创建临时目录失败: %w", err)
		}
		chunkFile := filepath.Join(tmpDir, fmt.Sprintf("chunk_%06d", param.ChunkIndex))
		f, err := os.Create(chunkFile)
		if err != nil {
			return fmt.Errorf("保存分片文件失败: %w", err)
		}
		if _, err := io.Copy(f, &contextReader{ctx: c.Request.Context(), reader: file}); err != nil {
			_ = f.Close()
			return fmt.Errorf("保存分片文件失败: %w", err)
		}
		if err := f.Close(); err != nil {
			return fmt.Errorf("保存分片文件失败: %w", err)
		}
	}
	return nil
}

func FileCompleteChunkUpload(c *gin.Context, param *ChunkCompleteParam) (*FileUploadResult, error) {
	state, err := loadChunkState(c, param.UploadID)
	if err != nil {
		return nil, err
	}
	param.Engine = state.Engine
	param.Bucket = state.Bucket
	param.FileKey = state.FileKey
	param.Name = state.Name
	param.FileSize = state.FileSize

	eng := storage.GetStorage(state.Engine)
	if eng == nil {
		return nil, fmt.Errorf("不支持的存储类型: %s", state.Engine)
	}

	now := time.Now()

	var storagePath string
	if supportsChunkUpload(eng) {
		path, err := completeChunkUpload(c.Request.Context(), eng, state.Bucket, state.FileKey, param.UploadID)
		if err != nil {
			return nil, fmt.Errorf("合并分片失败: %w", err)
		}
		storagePath = path
	} else {
		tmpDir := filepath.Join(os.TempDir(), "chunk_"+param.UploadID)
		defer os.RemoveAll(tmpDir)

		path, err := mergeAndStore(c.Request.Context(), eng, state.Bucket, state.FileKey, tmpDir, state.TotalChunks)
		if err != nil {
			return nil, err
		}
		storagePath = path
	}
	deleteChunkState(c, param.UploadID)

	ext := filepath.Ext(param.Name)
	downloadPath := storage.GetURL(state.Engine, state.Bucket, state.FileKey)

	thumbnail := ""
	if isImageExt(ext) {
		thumbnail = downloadPath
	}

	fileSizeKb, sizeInfo := formatFileSize(param.FileSize)

	entity := SysFile{
		Engine:       state.Engine,
		Bucket:       state.Bucket,
		FileKey:      param.FileKey,
		ObjName:      param.FileKey,
		Name:         param.Name,
		Suffix:       ext,
		SizeKb:       fileSizeKb,
		SizeInfo:     sizeInfo,
		StoragePath:  storagePath,
		DownloadPath: downloadPath,
		Thumbnail:    thumbnail,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := db.DB.WithContext(c.Request.Context()).Create(&entity).Error; err != nil {
		return nil, fmt.Errorf("保存文件记录失败: %w", err)
	}

	return &FileUploadResult{
		ID:           entity.ID,
		Engine:       entity.Engine,
		Bucket:       entity.Bucket,
		FileKey:      entity.FileKey,
		Name:         entity.Name,
		Suffix:       entity.Suffix,
		SizeKb:       entity.SizeKb,
		SizeInfo:     entity.SizeInfo,
		DownloadPath: entity.DownloadPath,
		Thumbnail:    entity.Thumbnail,
	}, nil
}

func FileAbortChunkUpload(c *gin.Context, param *ChunkAbortParam) error {
	state, err := loadChunkState(c, param.UploadID)
	if err != nil {
		return err
	}

	eng := storage.GetStorage(state.Engine)
	if eng == nil {
		return fmt.Errorf("不支持的存储类型: %s", state.Engine)
	}

	if supportsChunkUpload(eng) {
		err := abortChunkUpload(c.Request.Context(), eng, state.Bucket, state.FileKey, param.UploadID)
		deleteChunkState(c, param.UploadID)
		return err
	}

	tmpDir := filepath.Join(os.TempDir(), "chunk_"+param.UploadID)
	deleteChunkState(c, param.UploadID)
	return os.RemoveAll(tmpDir)
}

func mergeAndStore(ctx context.Context, eng storage.Engine, bucket, fileKey, tmpDir string, totalChunks int) (string, error) {
	for i := 0; i < totalChunks; i++ {
		chunkPath := filepath.Join(tmpDir, fmt.Sprintf("chunk_%06d", i))
		if _, err := os.Stat(chunkPath); err != nil {
			return "", fmt.Errorf("分片不完整: %d", i)
		}
	}

	pr, pw := io.Pipe()
	go func() {
		defer pw.Close()
		for i := 0; i < totalChunks; i++ {
			select {
			case <-ctx.Done():
				pw.CloseWithError(ctx.Err())
				return
			default:
			}
			chunkPath := filepath.Join(tmpDir, fmt.Sprintf("chunk_%06d", i))
			f, err := os.Open(chunkPath)
			if err != nil {
				pw.CloseWithError(fmt.Errorf("打开分片文件失败: %w", err))
				return
			}
			_, err = io.Copy(pw, &contextReader{ctx: ctx, reader: f})
			f.Close()
			if err != nil {
				pw.CloseWithError(fmt.Errorf("读取分片数据失败: %w", err))
				return
			}
		}
	}()

	return storeStream(ctx, eng, bucket, fileKey, pr, -1)
}

// ===== Download =====

func FileDownload(c *gin.Context, id string) error {
	var entity SysFile
	if err := db.DB.WithContext(c.Request.Context()).First(&entity, "id = ?", id).Error; err != nil {
		return fmt.Errorf("文件不存在")
	}

	return serveFile(c, &entity)
}

func FileDownloadByKey(c *gin.Context, bucket, fileKey string) error {
	if bucket == "" || fileKey == "" || strings.Contains(bucket, "..") || strings.Contains(fileKey, "..") ||
		strings.Contains(bucket, "/") || strings.Contains(bucket, "\\") ||
		strings.Contains(fileKey, "/") || strings.Contains(fileKey, "\\") {
		return fmt.Errorf("文件不存在")
	}

	var entity SysFile
	if err := db.DB.WithContext(c.Request.Context()).
		First(&entity, "bucket = ? AND file_key = ?", bucket, fileKey).Error; err != nil {
		return fmt.Errorf("文件不存在")
	}
	if entity.IsDownloadAuth {
		if _, exists := c.Get("login_id"); !exists {
			return fmt.Errorf("未授权/未登录")
		}
	}
	return serveFile(c, &entity)
}

func serveFile(c *gin.Context, entity *SysFile) error {
	if strings.EqualFold(entity.Engine, "LOCAL") && entity.StoragePath != "" {
		c.File(entity.StoragePath)
		return nil
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
