package demo

import (
	baidupanPlus "baiduNetdisk/baidupanSDK"
	"fmt"
	"log"
)

func QueryDirExample() {
	queryDirConfig := initByQueryCodeExample()
	resp, err := baidupanPlus.QueryDirWithConfig(&queryDirConfig)
	if err != nil {
		log.Fatalf("获取列表失败: %v", err)
	}
	fmt.Printf("请求成功！共找到 %d 个文件/目录\n", len(resp.List))
	fmt.Printf("%-20s %-10s %-10s\n", "文件名", "大小(字节)", "类型")
	fmt.Println("--------------------------------------------------")

	for _, file := range resp.List {
		fileName := "Unknown"
		if file.ServerFilename != nil {
			fileName = *file.ServerFilename
		}

		fileSize := int32(0)
		if file.Size != nil {
			fileSize = *file.Size
		}

		fileType := "文件"
		if file.Isdir != nil && *file.Isdir == 1 {
			fileType = "目录"
		}

		fmt.Printf("%-20s %-10d %-10s\n", fileName, fileSize, fileType)
	}
}
