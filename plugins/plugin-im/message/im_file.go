package message

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"path/filepath"
	"strings"

	"encoding/json"
	imModel "hei-gin/plugins/plugin-im/model"
	"hei-gin/sdk/config"
	"hei-gin/sdk/exception"
	"hei-gin/sdk/result"
	"hei-gin/sdk/storage"
	"hei-gin/sdk/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FileUploadResult struct {
	URL          string `json:"url"`
	FileKey      string `json:"file_key"`
	Bucket       string `json:"bucket"`
	Engine       string `json:"engine"`
	OriginalName string `json:"original_name"`
	FileSize     int64  `json:"file_size"`
	FileType     string `json:"file_type"`
}

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

func isImageExt(ext string) bool {
	switch strings.ToLower(ext) {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp", ".svg", ".ico", ".bmp", ".tiff":
		return true
	}
	return false
}

func maxUploadSize() int64 {
	if config.C != nil && config.C.App.UploadMaxSize > 0 {
		return config.C.App.UploadMaxSize
	}
	return 50 << 20
}

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

func storeStream(eng storage.Engine, bucket, fileKey string, reader io.Reader, size int64) (string, error) {
	return storeStreamWithContext(context.Background(), eng, bucket, fileKey, reader, size)
}

func storeStreamWithContext(ctx context.Context, eng storage.Engine, bucket, fileKey string, reader io.Reader, size int64) (string, error) {
	if ss, ok := eng.(storage.ContextSizedStreamer); ok {
		return ss.StoreStreamWithContext(ctx, bucket, fileKey, reader, size)
	}
	if ss, ok := eng.(storage.SizedStreamer); ok {
		return ss.StoreStreamWithSize(bucket, fileKey, reader, size)
	}
	return eng.StoreStream(bucket, fileKey, reader)
}

func UploadFile(c *gin.Context, senderID, senderType string) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		result.WriteError(c, exception.NewBusinessError("上传文件失败: "+err.Error(), 400))
		return
	}
	defer file.Close()

	if header.Size <= 0 || header.Size > maxUploadSize() {
		result.WriteError(c, exception.NewBusinessError(fmt.Sprintf("文件大小超过限制 (%d MB)", maxUploadSize()/(1<<20)), 400))
		return
	}
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext == "" || !allowedExtensions[ext] {
		result.WriteError(c, exception.NewBusinessError("不支持的文件类型: "+ext, 400))
		return
	}

	engineType := c.PostForm("engine")
	if engineType == "" {
		engineType = "LOCAL"
	}
	bucket := c.PostForm("bucket")
	if bucket == "" {
		bucket = "DEFAULT"
	}

	fileKey := utils.GenerateID() + ext

	eng := storage.GetStorage(engineType)
	if eng == nil {
		result.WriteError(c, exception.NewBusinessError("不支持的存储类型: "+engineType, 500))
		return
	}

	hr := newHashReader(file)
	storagePath, err := storeStreamWithContext(c.Request.Context(), eng, bucket, fileKey, hr, header.Size)
	if err != nil {
		result.WriteError(c, exception.NewBusinessError("保存文件失败: "+err.Error(), 500))
		return
	}

	checksum := hr.Sum()
	fileSizeKb, sizeInfo := formatFileSize(header.Size)

	thumbnail := ""
	if isImageExt(ext) {
		thumbnail = storage.GetURL(engineType, bucket, fileKey)
	}

	msgType := c.PostForm("msg_type")
	if msgType == "" {
		msgType = "FILE"
	}

	record := imModel.ImFile{
		ID:             utils.GenerateID(),
		Engine:         engineType,
		Bucket:         bucket,
		FileKey:        fileKey,
		Name:           header.Filename,
		Suffix:         ext,
		SizeKb:         fileSizeKb,
		SizeInfo:       sizeInfo,
		StoragePath:    storagePath,
		DownloadPath:   storage.GetURL(engineType, bucket, fileKey),
		Thumbnail:      thumbnail,
		Checksum:       checksum,
		ChecksumAlgo:   "sha256",
		ConversationID: c.PostForm("conversation_id"),
		SenderID:       senderID,
		SenderType:     senderType,
		MsgType:        msgType,
	}
	if err := CreateFile(c.Request.Context(), &record); err != nil {
		result.WriteError(c, exception.NewBusinessError("保存文件记录失败: "+err.Error(), 500))
		return
	}

	result.Success(c, &FileUploadResult{
		URL:          storage.GetURL(engineType, bucket, fileKey),
		FileKey:      fileKey,
		Bucket:       bucket,
		Engine:       engineType,
		OriginalName: header.Filename,
		FileSize:     header.Size,
		FileType:     ext,
	})
}

// ResolveFileURL constructs a full HTTP URL from message content and extra for IMAGE/FILE types.
func ResolveFileURL(content, extra string) string {
	if strings.HasPrefix(content, "http") {
		return content
	}
	if content == "" {
		return ""
	}

	engine := "LOCAL"
	bucket := "DEFAULT"

	if extra != "" {
		var meta struct {
			Engine string `json:"engine"`
			Bucket string `json:"bucket"`
		}
		if err := json.Unmarshal([]byte(extra), &meta); err == nil {
			if meta.Engine != "" {
				engine = meta.Engine
			}
			if meta.Bucket != "" {
				bucket = meta.Bucket
			}
		}
	}

	return storage.GetURL(engine, bucket, content)
}

func ServeUploadedFile(c *gin.Context, bucket, fileKey string) error {
	if bucket == "" || fileKey == "" || strings.Contains(bucket, "..") || strings.Contains(fileKey, "..") ||
		strings.Contains(bucket, "/") || strings.Contains(bucket, "\\") ||
		strings.Contains(fileKey, "/") || strings.Contains(fileKey, "\\") {
		return fmt.Errorf("文件不存在")
	}

	entity, err := FindFileByKey(c.Request.Context(), bucket, fileKey)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("文件不存在")
		}
		return fmt.Errorf("查询文件失败")
	}

	if strings.EqualFold(entity.Engine, "LOCAL") && entity.StoragePath != "" {
		c.File(entity.StoragePath)
		return nil
	}
	if entity.DownloadPath != "" {
		c.Redirect(302, entity.DownloadPath)
		return nil
	}
	return fmt.Errorf("文件路径为空")
}
