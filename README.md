# BaiduPan SDK 使用文档

本 SDK 提供了百度网盘文件的上传、下载和目录查询功能。以下是主要功能函数的使用说明。

## 1. 初始化配置

在使用 SDK 进行任何操作之前，建议先初始化基础配置。

### `NewBasicConfig`

初始化全局基础配置，包括 AccessToken、是否为 SVIP 以及日志路径。

**函数签名:**
```go
func NewBasicConfig(accessToken string, isSVIP bool, logPath string)
```

**参数说明:**
*   `accessToken`: 百度网盘接口调用凭证。
*   `isSVIP`: 是否为超级会员（影响分片上传大小）。
*   `logPath`: 日志文件存储路径。

**示例:**
```go
baidupanSDK.NewBasicConfig("your-access-token", true, "./sdk.log")
```

---

## 2. 文件上传

文件上传分为配置初始化和执行上传两步。

### `NewUploadFileConfig`

创建文件上传配置对象。

**函数签名:**
```go
func NewUploadFileConfig(localPath, remotePath string) UploadFileConfig
```

**参数说明:**
*   `localPath`: 本地文件绝对路径。
*   `remotePath`: 网盘中的目标路径（包含文件名）。

**返回值:**
*   `UploadFileConfig`: 上传配置对象。

### `UploadFileWithConfig`

根据配置执行文件上传操作。自动处理预上传、分片上传和文件合并。

**函数签名:**
```go
func UploadFileWithConfig(uploadFileConfig UploadFileConfig) error
```

**示例:**
```go
// 1. 创建配置
uploadConfig := baidupanSDK.NewUploadFileConfig("/local/file.txt", "/apps/myapp/file.txt")

// 2. 执行上传
err := baidupanSDK.UploadFileWithConfig(uploadConfig)
if err != nil {
    fmt.Printf("Upload failed: %v\n", err)
} else {
    fmt.Println("Upload success!")
}
```

---

## 3. 文件下载

文件下载同样分为配置初始化和执行下载两步。

### `NewDownloadFileConfig`

创建文件下载配置对象。

**函数签名:**
```go
func NewDownloadFileConfig(localPath, remotePath string) DownloadFileConfig
```

**参数说明:**
*   `localPath`: 本地保存路径（包含文件名）。
*   `remotePath`: 网盘中的源文件路径。

**返回值:**
*   `DownloadFileConfig`: 下载配置对象。

### `DownloadFileWithConfig`

根据配置执行文件下载操作。自动获取文件 dlink 并完成下载。

**函数签名:**
```go
func DownloadFileWithConfig(config DownloadFileConfig) error
```

**示例:**
```go
// 1. 创建配置
downloadConfig := baidupanSDK.NewDownloadFileConfig("/local/save/path.txt", "/apps/myapp/remote.txt")

// 2. 执行下载
err := baidupanSDK.DownloadFileWithConfig(downloadConfig)
if err != nil {
    fmt.Printf("Download failed: %v\n", err)
} else {
    fmt.Println("Download success!")
}
```

---

## 4. 目录查询

查询网盘指定目录下的文件列表。

### `NewQueryDirConfig`

创建目录查询配置对象。

**函数签名:**
```go
func NewQueryDirConfig(dir string, limit int) QueryDirConfig
```

**参数说明:**
*   `dir`: 网盘目录路径。
*   `limit`: 返回的文件数量限制。

**返回值:**
*   `QueryDirConfig`: 查询配置对象。

**注意:** 虽然 `QueryDir` 函数在内部使用，但通常可以通过配置对象配合其他逻辑使用，或者直接调用 `QueryDir` (如果已公开)。

**示例:**
```go
// 创建配置
queryConfig := baidupanSDK.NewQueryDirConfig("/apps/myapp", 100)
// 调用查询 (假设 QueryDir 是公开的)
// resp, err := baidupansdk.QueryDir(&queryConfig)
```

---

## 完整示例

```go
package main

import (
	"baiduNetdisk/baidupan_SDK"
	"fmt"
)

func main() {
	accessToken := "your-access-token"
	
	// 1. 初始化
	baidupanSDK.NewBasicConfig(accessToken, true, "./app.log")

	// 2. 上传文件
	fmt.Println("--- Uploading ---")
	upConfig := baidupanSDK.NewUploadFileConfig("/tmp/test.jpg", "/apps/test_app/test.jpg")
	if err := baidupanSDK.UploadFileWithConfig(upConfig); err != nil {
		fmt.Printf("Upload error: %v\n", err)
	}

	// 3. 下载文件
	fmt.Println("--- Downloading ---")
	downConfig := baidupanSDK.NewDownloadFileConfig("/tmp/downloaded_test.jpg", "/apps/test_app/test.jpg")
	if err := baidupanSDK.DownloadFileWithConfig(downConfig); err != nil {
		fmt.Printf("Download error: %v\n", err)
	}
}
```
