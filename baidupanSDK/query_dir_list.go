package baidupanSDK

import (
	"encoding/json"
	"fmt"
	openapi "github.com/S-zhi/baidupansdk/openxpanapi"
)

// FileListResponse 文件列表响应结构
type FileListResponse struct {
	Errno     int32                        `json:"errno"`
	Guid      int32                        `json:"guid"`
	List      []openapi.Filecreateresponse `json:"list"`
	RequestId int64                        `json:"request_id"`
}

// QueryDir 获取文件列表
// 参数:
//
//	accessToken: 访问令牌
//	dir: 目录路径，如 "/apps/myapp"
//	start: 起始位置，通常为 "0"
//	limit: 返回条数，默认 100
//
// 返回值:
//
//	*FileListResponse: 文件列表数据
//	error: 错误信息
func QueryDirWithConfig(qConfig *QueryDirConfig) (*FileListResponse, error) {
	if qConfig == nil {
		Warn("QueryDir: 参数为空，调用 defaultQueryDirConfig")
		qConfig = &defaultQueryDirConfig
		if qConfig == nil {
			Error("QueryDir: defaultQueryDirConfig is nil")
			return nil, fmt.Errorf("QueryDir: defaultQueryDirConfig is nil")
		}
	}

	apiXpanfilelistRequest := client.FileinfoApi.Xpanfilelist(ctx).
		AccessToken(qConfig.AccessToken).
		Dir(qConfig.Dir).
		Start("0").
		Limit(qConfig.Limit)

	// SDK 返回的是原始 JSON 字符串
	jsonStr, _, err := client.FileinfoApi.XpanfilelistExecute(apiXpanfilelistRequest)
	if err != nil {
		Error(fmt.Sprintf("Failed to execute Xpanfilelist: %v", err))
		return nil, err
	}

	var fileListResp FileListResponse
	err = json.Unmarshal([]byte(jsonStr), &fileListResp)
	if err != nil {
		Error(fmt.Sprintf("Failed to unmarshal file list response: %v", err))
		return nil, err
	}

	if fileListResp.Errno != 0 {
		return nil, fmt.Errorf("get file list failed with errno: %d", fileListResp.Errno)
	}

	Info(fmt.Sprintf("Successfully retrieved file list for: %s, count: %d", qConfig.Dir, len(fileListResp.List)))
	return &fileListResp, nil
}
