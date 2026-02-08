package demo

import (
	baidupanSDK "baiduNetdisk/baidupan_SDK"
	"fmt"
)

func DownloadExample() {

	// 创建下载配置
	downloadFileConfig := initByDownloadFileExample()

	// 使用配置进行下载
	err := baidupanSDK.DownloadFileWithConfig(downloadFileConfig)

	if err != nil {
		fmt.Printf("下载失败: %v\n", err)
	} else {
		fmt.Println("下载成功!")
	}
}
