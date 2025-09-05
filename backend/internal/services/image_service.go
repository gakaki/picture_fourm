package services

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"nano-bana-qwen/internal/config"
	"nano-bana-qwen/internal/models"

	"github.com/nfnt/resize"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ImageService struct {
}

// NewImageService 创建图片服务实例
func NewImageService() *ImageService {
	return &ImageService{}
}

// SaveImageFromURL 从URL下载并保存图片
func (s *ImageService) SaveImageFromURL(imageURL, generationID string) (localPath, thumbnailPath string, err error) {
	// 确保目录存在
	if err := s.ensureDirectories(); err != nil {
		return "", "", fmt.Errorf("创建目录失败: %v", err)
	}

	// 下载图片
	resp, err := http.Get(imageURL)
	if err != nil {
		return "", "", fmt.Errorf("下载图片失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("下载图片失败，状态码: %d", resp.StatusCode)
	}

	// 读取图片数据
	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("读取图片数据失败: %v", err)
	}

	// 生成文件名
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("generated_%s_%s.png", timestamp, generationID[:8])
	
	// 保存原图
	localPath = filepath.Join(config.AppConfig.GeneratedPath, filename)
	if err := s.saveImageFile(localPath, imageData); err != nil {
		return "", "", fmt.Errorf("保存原图失败: %v", err)
	}

	// 生成缩略图
	thumbnailFilename := fmt.Sprintf("thumb_%s", filename)
	thumbnailPath = filepath.Join(config.AppConfig.ThumbnailPath, thumbnailFilename)
	if err := s.generateThumbnail(imageData, thumbnailPath); err != nil {
		return "", "", fmt.Errorf("生成缩略图失败: %v", err)
	}

	// 保存图片元数据到数据库
	imageInfo := models.Image{
		ID:               primitive.NewObjectID(),
		Filename:         filename,
		OriginalFilename: filename,
		FilePath:         localPath,
		ThumbnailPath:    thumbnailPath,
		FileSize:         int64(len(imageData)),
		Format:           "PNG",
		CreatedAt:        time.Now(),
		Deleted:          false,
	}

	// 获取图片尺寸
	if width, height, err := s.getImageDimensions(imageData); err == nil {
		imageInfo.Width = width
		imageInfo.Height = height
	}

	// 保存到数据库
	_, err = MongoDB.Collection("images").InsertOne(context.Background(), imageInfo)
	if err != nil {
		return "", "", fmt.Errorf("保存图片信息到数据库失败: %v", err)
	}

	return localPath, thumbnailPath, nil
}

// SaveImageFromBase64 从base64数据保存图片
func (s *ImageService) SaveImageFromBase64(base64Data, generationID string) (localPath, thumbnailPath string, err error) {
	// 处理base64数据（去除data:image/xxx;base64,前缀）
	if strings.Contains(base64Data, ",") {
		parts := strings.Split(base64Data, ",")
		if len(parts) > 1 {
			base64Data = parts[1]
		}
	}

	// 解码base64
	imageData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", "", fmt.Errorf("base64解码失败: %v", err)
	}

	// 确保目录存在
	if err := s.ensureDirectories(); err != nil {
		return "", "", fmt.Errorf("创建目录失败: %v", err)
	}

	// 生成文件名
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("generated_%s_%s.png", timestamp, generationID[:8])
	
	// 保存原图
	localPath = filepath.Join(config.AppConfig.GeneratedPath, filename)
	if err := s.saveImageFile(localPath, imageData); err != nil {
		return "", "", fmt.Errorf("保存原图失败: %v", err)
	}

	// 生成缩略图
	thumbnailFilename := fmt.Sprintf("thumb_%s", filename)
	thumbnailPath = filepath.Join(config.AppConfig.ThumbnailPath, thumbnailFilename)
	if err := s.generateThumbnail(imageData, thumbnailPath); err != nil {
		return "", "", fmt.Errorf("生成缩略图失败: %v", err)
	}

	return localPath, thumbnailPath, nil
}

// ensureDirectories 确保必要的目录存在
func (s *ImageService) ensureDirectories() error {
	dirs := []string{
		config.AppConfig.GeneratedPath,
		config.AppConfig.ThumbnailPath,
		config.AppConfig.TempPath,
		config.AppConfig.UploadPath,
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("创建目录 %s 失败: %v", dir, err)
		}
	}

	return nil
}

// saveImageFile 保存图片文件
func (s *ImageService) saveImageFile(filePath string, imageData []byte) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(imageData)
	return err
}

// generateThumbnail 生成缩略图
func (s *ImageService) generateThumbnail(imageData []byte, thumbnailPath string) error {
	// 解码图片
	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return fmt.Errorf("解码图片失败: %v", err)
	}

	// 调整尺寸为200x200
	thumbnail := resize.Resize(200, 200, img, resize.Lanczos3)

	// 保存缩略图
	file, err := os.Create(thumbnailPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 根据文件扩展名决定编码格式
	ext := strings.ToLower(filepath.Ext(thumbnailPath))
	switch ext {
	case ".jpg", ".jpeg":
		err = jpeg.Encode(file, thumbnail, &jpeg.Options{Quality: 85})
	case ".png":
		err = png.Encode(file, thumbnail)
	default:
		err = png.Encode(file, thumbnail) // 默认PNG
	}

	return err
}

// getImageDimensions 获取图片尺寸
func (s *ImageService) getImageDimensions(imageData []byte) (width, height int, err error) {
	config, _, err := image.DecodeConfig(bytes.NewReader(imageData))
	if err != nil {
		return 0, 0, err
	}
	return config.Width, config.Height, nil
}

// GetImageByID 根据ID获取图片信息
func (s *ImageService) GetImageByID(id primitive.ObjectID) (*models.Image, error) {
	var image models.Image
	err := MongoDB.Collection("images").FindOne(context.Background(), bson.M{
		"_id":     id,
		"deleted": false,
	}).Decode(&image)

	if err != nil {
		return nil, err
	}

	return &image, nil
}

// DeleteImage 删除图片（软删除）
func (s *ImageService) DeleteImage(id primitive.ObjectID) error {
	update := bson.M{
		"$set": bson.M{
			"deleted":        true,
			"deleted_at":     time.Now(),
			"deleted_reason": "用户删除",
		},
	}

	result, err := MongoDB.Collection("images").UpdateOne(context.Background(), bson.M{
		"_id":     id,
		"deleted": false,
	}, update)

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("图片不存在")
	}

	return nil
}

// ListImages 获取图片列表
func (s *ImageService) ListImages(page, pageSize int, prompt string) ([]models.Image, int64, error) {
	// 构建查询条件
	filter := bson.M{"deleted": false}
	
	if prompt != "" {
		filter["prompt_text"] = bson.M{"$regex": prompt, "$options": "i"}
	}

	// 获取总数
	total, err := MongoDB.Collection("images").CountDocuments(context.Background(), filter)
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	skip := int64((page - 1) * pageSize)
	limit := int64(pageSize)
	cursor, err := MongoDB.Collection("images").Find(context.Background(), filter, &options.FindOptions{
		Skip:  &skip,
		Limit: &limit,
		Sort:  bson.D{{Key: "created_at", Value: -1}},
	})
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.Background())

	var images []models.Image
	if err := cursor.All(context.Background(), &images); err != nil {
		return nil, 0, err
	}

	return images, total, nil
}

// CleanupTempFiles 清理临时文件
func (s *ImageService) CleanupTempFiles(olderThan time.Duration) error {
	tempDir := config.AppConfig.TempPath
	
	files, err := os.ReadDir(tempDir)
	if err != nil {
		return err
	}

	cutoff := time.Now().Add(-olderThan)
	
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		info, err := file.Info()
		if err != nil {
			continue
		}

		if info.ModTime().Before(cutoff) {
			filePath := filepath.Join(tempDir, file.Name())
			if err := os.Remove(filePath); err != nil {
				fmt.Printf("删除临时文件失败 %s: %v\n", filePath, err)
			}
		}
	}

	return nil
}