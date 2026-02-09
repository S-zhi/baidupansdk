package main

import (
	"github.com/S-zhi/baidupansdk/baidupanplus"
	demo "github.com/S-zhi/baidupansdk/demo"
)

func main() {
	//demo.QueryDirExample()
	//demo.UploadFileExample()
	demo.DownloadExample()
	baidupanplus.NewQueryDirConfig("/apps/myapp", 100)
}
