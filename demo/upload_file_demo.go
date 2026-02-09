package demo

import (
	baidupanPlus "github.com/S-zhi/baidupansdk/baidupanplus"
	"log"
)

func UploadFileExample() { // 分片大小，建议 4MB
	uploadFileConfig := initByUploadFileExample()
	err := baidupanPlus.UploadFileWithConfig(uploadFileConfig)
	if err != nil {
		log.Printf("上传失败: %v", err)
		return
	}
	log.Println("上传成功！")
}
