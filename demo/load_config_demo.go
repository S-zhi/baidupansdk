package demo

import (
	BaiduPlus "github.com/S-zhi/baidupansdk/baidupanSDK"
)

// initByQueryCodeExample 通过代码初始化
func initByQueryCodeExample() BaiduPlus.QueryDirConfig {
	BaiduPlus.NewBasicConfig(ACCESS_TOKEN,
		true,
		"/Users/wenzhengfeng/code/go/baiduNetdisk/logs/baiduPanSDK.log")
	return BaiduPlus.NewQueryDirConfig("/project/luckyProject/weights/", 100)
}

// initByUploadFileExample 通过代码初始化
func initByUploadFileExample() BaiduPlus.UploadFileConfig {
	BaiduPlus.NewBasicConfig(ACCESS_TOKEN,
		true,
		"/Users/wenzhengfeng/code/go/baiduNetdisk/logs/baiduPanSDK.log")
	return BaiduPlus.NewUploadFileConfig("/Users/wenzhengfeng/code/go/baiduNetdisk/demo/test_2_8.txt",
		"/project/luckyProject/weights/test_2_8.txt")
}

// initByDownloadFileExample 通过代码初始化
func initByDownloadFileExample() BaiduPlus.DownloadFileConfig {
	BaiduPlus.NewBasicConfig(ACCESS_TOKEN,
		true,
		"/Users/wenzhengfeng/code/go/baiduNetdisk/logs/baiduPanSDK.log")
	return BaiduPlus.NewDownloadFileConfig("/Users/wenzhengfeng/code/go/baiduNetdisk/demo/test_2_8_1.txt",
		"/project/luckyProject/weights/test_2_8.txt")
}
