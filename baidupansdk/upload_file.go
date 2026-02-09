package baidupansdk

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/S-zhi/baidupansdk/baidupansdk/tool"
	openapi "github.com/S-zhi/baidupansdk/openxpanapi"
)

var client = *(openapi.NewAPIClient(openapi.NewConfiguration()))
var ctx = context.Background()

const (
	isdir    = 0
	autoinit = 1
)

// 预上传文档 : https://pan.baidu.com/union/doc/3ksg0s9r7?from=open-sdk-go

// PrecreateFile 预创建文件，用于在远程服务器上预先创建文件结构
// 返回值:
//
//	uploadid: 上传任务ID
//	md5List: 分片MD5列表
//	error: 错误信息
func PrecreateFile(accessToken string, remotePath string, localPath string, shardSize int64) (string, []string, error) {
	fileSize, err := tools.GetFileSizeByPath(localPath)
	Info(fmt.Sprintf("File size: %d", fileSize))
	if err != nil {
		Error(fmt.Sprintf("Failed to get file size: %v", err))
		return "", nil, err
	}

	// 计算所有分片的MD5
	var md5List []string
	err = ProcessFileInShards(localPath, shardSize, func(index int, data []byte, isLast bool) error {
		md5Code := md5.New()
		md5Code.Write(data)
		md5Str := hex.EncodeToString(md5Code.Sum(nil))
		md5List = append(md5List, md5Str)
		return nil
	})
	if err != nil {
		Error(fmt.Sprintf("Failed to calculate shard MD5s: %v", err))
		return "", nil, err
	}

	md5ListByte, _ := json.Marshal(md5List)
	md5ListStr := string(md5ListByte)
	apiXpanfileprecreateRequest := client.FileuploadApi.Xpanfileprecreate(ctx).
		AccessToken(accessToken).
		Path(remotePath).
		Autoinit(autoinit).
		Size(int32(fileSize)).
		Isdir(isdir).
		BlockList(md5ListStr)

	fileprecreateresponse, _, err := client.FileuploadApi.XpanfileprecreateExecute(apiXpanfileprecreateRequest)
	if err != nil {
		Error(fmt.Sprintf("Failed to execute Xpanfileprecreate: %v", err))
		return "", nil, err
	}

	if fileprecreateresponse.GetErrno() != 0 {
		return "", nil, fmt.Errorf("precreate failed with errno: %d", fileprecreateresponse.GetErrno())
	}

	return fileprecreateresponse.GetUploadid(), md5List, nil
}

// UploadPart 分片上传
func UploadPart(accessToken string, remotePath string, uploadID string, partSeq int, partData []byte) error {

	tmpFile, err := os.CreateTemp("", "part-*")
	if err != nil {
		return err
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			Error(fmt.Sprintf("Failed to remove temporary file: %v", err))
		}
	}(tmpFile.Name())
	defer func(tmpFile *os.File) {
		err := tmpFile.Close()
		if err != nil {
			Error(fmt.Sprintf("Failed to close temporary file: %v", err))
		}
	}(tmpFile)

	if _, err := tmpFile.Write(partData); err != nil {
		return err
	}
	if _, err := tmpFile.Seek(0, 0); err != nil {
		return err
	}

	apiXpanfileuploadRequest := client.FileuploadApi.Pcssuperfile2(ctx).
		AccessToken(accessToken).
		Path(remotePath).
		Uploadid(uploadID).
		Type_("tmpfile").
		Partseq(fmt.Sprintf("%d", partSeq)).
		File(tmpFile)

	_, response, err := client.FileuploadApi.Pcssuperfile2Execute(apiXpanfileuploadRequest)
	if err != nil {
		Error(fmt.Sprintf("Failed to upload part %d: %v, status: %d", partSeq, err, response.StatusCode))
		return err
	}

	Info(fmt.Sprintf("Successfully uploaded part %d", partSeq))
	return nil
}

// CreateFile 合并分片创建文件
func CreateFile(accessToken string, remotePath string, uploadID string, fileSize int64, md5List []string) error {
	md5ListByte, _ := json.Marshal(md5List)
	md5ListStr := string(md5ListByte)

	apiXpanfilecreateRequest := client.FileuploadApi.Xpanfilecreate(ctx).
		AccessToken(accessToken).
		Path(remotePath).
		Isdir(isdir).
		Size(int32(fileSize)).
		Uploadid(uploadID).
		BlockList(md5ListStr)

	filecreateresponse, _, err := client.FileuploadApi.XpanfilecreateExecute(apiXpanfilecreateRequest)
	if err != nil {
		Error(fmt.Sprintf("Failed to execute Xpanfilecreate: %v", err))
		return err
	}

	if filecreateresponse.GetErrno() != 0 {
		return fmt.Errorf("create file failed with errno: %d", filecreateresponse.GetErrno())
	}

	Info(fmt.Sprintf("Successfully created file: %s", remotePath))
	return nil
}

// UploadFileWithConfig 完整上传流程封装
func UploadFileWithConfig(uploadFileConfig UploadFileConfig) error {
	if uploadFileConfig == (UploadFileConfig{}) {
		Warn("UploadFileConfig is empty, Use defaultUploadFileConfig")
		uploadFileConfig = defaultUploadFileConfig
	}
	accessToken := uploadFileConfig.AccessToken
	remotePath := uploadFileConfig.RemotePath
	localPath := uploadFileConfig.LocalPath
	Info("未实现普通会员的情况")
	var shardSize int64
	if uploadFileConfig.IsSVIP {
		shardSize = int64(4 * 1024 * 1024)
		Info("授权用户为普通用户时，单个分片大小固定为4MB，单文件总大小上限为4GB")
	} else {
		shardSize = int64(32 * 1024 * 1024)
		Info("授权用户为超级会员时，用户单个分片大小上限为32MB，单文件总大小上限为20GB")
	}
	// 1. 预上传
	uploadID, md5List, err := PrecreateFile(accessToken, remotePath, localPath, shardSize)
	if err != nil {
		return err
	}

	// 2. 分片上传
	err = ProcessFileInShards(localPath, shardSize, func(index int, data []byte, isLast bool) error {
		return UploadPart(accessToken, remotePath, uploadID, index, data)
	})
	if err != nil {
		return err
	}

	// 3. 创建文件
	fileSize, _ := tools.GetFileSizeByPath(localPath)
	return CreateFile(accessToken, remotePath, uploadID, fileSize, md5List)
}

// =================================== 分片处理器 ===================================

// ShardProcessor 分片处理器接口
type ShardProcessor func(index int, data []byte, isLast bool) error

// ProcessFileInShards 流式处理文件分片
func ProcessFileInShards(filePath string, shardSize int64, processor ShardProcessor) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("Failed to close file: %v", err)
		}
	}(file)

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	totalSize := fileInfo.Size()

	shardCount := (totalSize + shardSize - 1) / shardSize
	buffer := make([]byte, shardSize)
	currentIndex := 0

	for {
		n, err := file.Read(buffer)
		if n > 0 {
			isLast := currentIndex == int(shardCount-1)
			err := processor(currentIndex, buffer[:n], isLast)
			if err != nil {
				return err
			}
			currentIndex++
		}

		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
	}

	return nil
}
