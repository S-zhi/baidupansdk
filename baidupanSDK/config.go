package baidupanSDK

// config.go
import (
	"encoding/json"
	"fmt"
	"os"
)

// OperateType 操作类型枚举
type OperateType string

const (
	UploadFileOperate   OperateType = "uploadFile"
	DownloadFileOperate OperateType = "downloadFile"
	QueryDirOperate     OperateType = "queryDir"
	UnknownOperate      OperateType = "unknown"
)

var config Config
var defaultQueryDirConfig QueryDirConfig
var defaultUploadFileConfig UploadFileConfig
var defaultDownloadFileConfig DownloadFileConfig

// Config 配置结构体
type Config struct {
	AccessToken string      `json:"access_token"` // 访问令牌
	Operate     OperateType `json:"operate"`      // 操作类型
	IsSVIP      bool        `json:"is_svip"`      // 是否为超级会员
	LogPath     string      `json:"log_path"`
}

// UploadFileConfig 上传文件配置结构体
type UploadFileConfig struct {
	Config
	LocalPath  string `json:"local_path"`  // 本地文件路径
	RemotePath string `json:"remote_path"` // 远程文件路径
}

// DownloadFileConfig 下载文件配置结构体
type DownloadFileConfig struct {
	Config
	LocalPath  string `json:"local_path"`  // 本地文件路径
	RemotePath string `json:"remote_path"` // 远程文件路径
}

// QueryDirConfig 查询目录配置结构体
type QueryDirConfig struct {
	Config
	Dir   string `json:"dir"`   // 查询远程目录路径
	Limit int32  `json:"limit"` // 查询文件列表限制数量
}

// NewBasicConfig 实例化配置对象
func NewBasicConfig(accessToken string, isSVIP bool, logPath string) {
	config = Config{
		AccessToken: accessToken,
		Operate:     OperateType(UnknownOperate),
		IsSVIP:      isSVIP,
		LogPath:     logPath,
	}
	Info("NewBasicConfig 创建成功: %v", config)
	logInit()
}

// NewUploadFileConfig 实例化 UploadFileConfig
func NewUploadFileConfig(localPath, remotePath string) UploadFileConfig {
	if config.AccessToken == "" || config.Operate == "" {
		Error("BaiduPanPlus Config is not initialized,NewUploadFileConfig Failed")
		panic("BaiduPanPlus Config is not initialized,NewUploadFileConfig Failed")
	}
	defaultUploadFileConfig = UploadFileConfig{
		Config:     config, // 复用 Config 字段
		LocalPath:  localPath,
		RemotePath: remotePath,
	}
	defaultQueryDirConfig.Config.Operate = UploadFileOperate
	return defaultUploadFileConfig
}

// NewDownloadFileConfig 实例化 DownloadFileConfig
func NewDownloadFileConfig(localPath, remotePath string) DownloadFileConfig {
	if config.AccessToken == "" || config.Operate == "" {
		Error("BaiduPanPlus Config is not initialized,NewDownloadFileConfig Failed")
		panic("BaiduPanPlus Config is not initialized,NewDownloadFileConfig Failed")
	}
	defaultDownloadFileConfig = DownloadFileConfig{
		Config:     config, // 复用 Config 字段
		LocalPath:  localPath,
		RemotePath: remotePath,
	}
	defaultDownloadFileConfig.Config.Operate = DownloadFileOperate
	return defaultDownloadFileConfig
}

// NewQueryDirConfig 实例化 QueryDirConfig
func NewQueryDirConfig(dir string, limit int) QueryDirConfig {
	if config.AccessToken == "" || config.Operate == "" {
		Error("BaiduPanPlus Config is not initialized,NewQueryDirConfig Failed")
		panic("BaiduPanPlus Config is not initialized,NewQueryDirConfig Failed")
	}
	limit32 := int32(limit)
	defaultQueryDirConfig = QueryDirConfig{
		Config: config,
		Dir:    dir,
		Limit:  limit32,
	}
	defaultQueryDirConfig.Config.Operate = QueryDirOperate
	return defaultQueryDirConfig
}

// LoadUploadFileConfigFromFile 从配置文件加载 UploadFileConfig
func LoadUploadFileConfigFromFile(filePath string) error {
	err := LoadConfigFromFile(filePath)
	config.Operate = DownloadFileOperate
	return err
}

// LoadDownloadFileConfigFromFile 从配置文件加载 DownloadFileConfig
func LoadDownloadFileConfigFromFile(filePath string) error {
	err := LoadConfigFromFile(filePath)
	config.Operate = DownloadFileOperate
	return err
}

// LoadQueryDirConfigFromFile 从配置文件加载 QueryDirConfig
func LoadQueryDirConfigFromFile(filePath string) error {
	err := LoadConfigFromFile(filePath)
	config.Operate = QueryDirOperate
	return err
}

// LoadConfigFromFile 从配置文件加载配置
func LoadConfigFromFile(filePath string) error {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}
	if err := json.Unmarshal(file, &config); err != nil {
		return fmt.Errorf("failed to parse config file: %v", err)
	}

	// 验证操作类型是否合法
	switch config.Operate {
	case UploadFileOperate, DownloadFileOperate, QueryDirOperate:
	default:
		return fmt.Errorf("invalid operate type: %s", config.Operate)
	}

	return nil
}
