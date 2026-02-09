package baidupansdk

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
)

// FileMeta 百度网盘文件元数据结构
type FileMeta struct {
	FsId     int64  `json:"fs_id"`
	Path     string `json:"path"`
	Filename string `json:"server_filename"`
	Size     int64  `json:"size"`
	Dlink    string `json:"dlink"`
}

// FileMetasResponse 文件详情响应
type FileMetasResponse struct {
	Errno     int32       `json:"errno"`
	List      []FileMeta  `json:"list"`
	RequestId interface{} `json:"request_id"`
}

// GetFileMetas 获取文件详情（包含 dlink）
func GetFileMetas(accessToken string, fsids []int64) (*FileMetasResponse, error) {
	fsidsByte, _ := json.Marshal(fsids)
	fsidsStr := string(fsidsByte)

	apiXpanmetasRequest := client.MultimediafileApi.Xpanmultimediafilemetas(ctx).
		AccessToken(accessToken).
		Fsids(fsidsStr).
		Dlink("1") // 必须设置为 "1" 才会返回下载链接

	jsonStr, _, err := client.MultimediafileApi.XpanmultimediafilemetasExecute(apiXpanmetasRequest)
	if err != nil {
		Error(fmt.Sprintf("Failed to execute Xpanmultimediafilemetas: %v", err))
		return nil, err
	}

	var metasResp FileMetasResponse
	err = json.Unmarshal([]byte(jsonStr), &metasResp)
	if err != nil {
		return nil, err
	}

	if metasResp.Errno != 0 {
		return nil, fmt.Errorf("get file metas failed with errno: %d", metasResp.Errno)
	}

	return &metasResp, nil
}

// findFileFsIdByPath 根据路径查找文件的fs_id（支持分页查找）
func findFileFsIdByPath(accessToken string, dir string, filename string) (int64, error) {
	start := 0
	limit := 1000

	for {
		// 调用 SDK 的 list 接口
		apiReq := client.FileinfoApi.Xpanfilelist(ctx).
			AccessToken(accessToken).
			Dir(dir).
			Start(strconv.Itoa(start)).
			Limit(int32(limit))

		jsonStr, _, err := client.FileinfoApi.XpanfilelistExecute(apiReq)
		if err != nil {
			return 0, fmt.Errorf("execute list api failed: %v", err)
		}

		var fileListResp FileListResponse
		err = json.Unmarshal([]byte(jsonStr), &fileListResp)
		if err != nil {
			return 0, fmt.Errorf("unmarshal list response failed: %v", err)
		}

		if fileListResp.Errno != 0 {
			return 0, fmt.Errorf("get file list failed with errno: %d", fileListResp.Errno)
		}

		// 遍历查找
		for _, file := range fileListResp.List {
			if file.GetServerFilename() == filename && file.GetIsdir() == 0 {
				return file.GetFsId(), nil
			}
		}

		// 如果返回的数量小于限制，说明没有更多文件了
		if len(fileListResp.List) < limit {
			break
		}
		// 否则继续下一页
		start += limit
	}

	return 0, fmt.Errorf("file not found: %s/%s", dir, filename)
}

// DownloadFileWithConfig 使用DownloadFileConfig配置下载文件
func DownloadFileWithConfig(config DownloadFileConfig) error {
	// 验证配置参数
	if config.AccessToken == "" {
		Error("AccessToken不能为空")
		return fmt.Errorf("access token is required")
	}
	if config.RemotePath == "" {
		Error("RemotePath不能为空")
		return fmt.Errorf("remote path is required")
	}
	if config.LocalPath == "" {
		Error("LocalPath不能为空")
		return fmt.Errorf("local path is required")
	}

	Info("开始下载文件: remote=%s, local=%s", config.RemotePath, config.LocalPath)

	// 1. 获取文件所在目录和文件名
	dir := path.Dir(config.RemotePath)
	filename := path.Base(config.RemotePath)

	// 2. 查找文件获取 fs_id
	targetFsId, err := findFileFsIdByPath(config.AccessToken, dir, filename)
	if err != nil {
		Error("查找文件失败: %v", err)
		return err
	}

	// 3. 获取文件详情（获取dlink）
	metasResp, err := GetFileMetas(config.AccessToken, []int64{targetFsId})
	if err != nil {
		Error("获取文件详情失败: %v", err)
		return err
	}

	if len(metasResp.List) == 0 {
		Error("未获取到文件元数据")
		return fmt.Errorf("no file meta data found")
	}

	dlink := metasResp.List[0].Dlink
	if dlink == "" {
		Error("未获取到下载链接")
		return fmt.Errorf("dlink not found")
	}

	// 4. 下载文件
	err = DownloadFile(config.AccessToken, dlink, config.LocalPath)
	if err != nil {
		Error("下载文件失败: %v", err)
		return err
	}

	Info("下载流程完成")
	return nil
}

// DownloadFile 下载文件
func DownloadFile(accessToken string, dlink string, localPath string) error {
	// 解析 dlink URL
	u, err := url.Parse(dlink)
	if err != nil {
		Error("解析dlink失败: %v", err)
		return err
	}

	// 百度网盘下载必须携带 User-Agent: pan.baidu.com
	// 并且 access_token 需要作为 query 参数传递
	q := u.Query()
	q.Set("access_token", accessToken)
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		Error("创建HTTP请求失败: %v", err)
		return err
	}
	req.Header.Set("User-Agent", "pan.baidu.com")

	Info("发送下载请求到: %s", u.String())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		Error("HTTP请求失败: %v", err)
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			Error("关闭响应体失败: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		Error("下载失败，状态码: %s", resp.Status)
		// 尝试读取body看是否有错误信息
		bodyBytes, _ := io.ReadAll(resp.Body)
		Error("错误响应: %s", string(bodyBytes))
		return fmt.Errorf("download failed with status: %s", resp.Status)
	}

	out, err := os.Create(localPath)
	if err != nil {
		Error("创建本地文件失败: %v", err)
		return err
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			Error("关闭本地文件失败: %v", err)
		}
	}(out)

	Info("开始写入文件到: %s", localPath)
	// 使用 io.Copy 流式写入，避免内存溢出
	written, err := io.Copy(out, resp.Body)
	if err != nil {
		Error("写入文件失败: %v", err)
		return err
	}

	Info("文件下载成功: %s, 大小: %d bytes", localPath, written)
	return nil
}
