package tools

import "os"

// GetFileSizeByPath 获取文件的大小(B)
func GetFileSizeByPath(filePath string) (int64, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}

	return fileInfo.Size(), nil
}
