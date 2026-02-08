package baidupanSDK

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/natefinch/lumberjack.v2"
)

func logInit() {
	// 创建日志目录
	if err := ensureLogDir(config.LogPath); err != nil {
		fmt.Printf("failed to create log directory: %v\n", err)
		return
	}
	log.SetFlags(log.LstdFlags) // 使用标准时间戳 + 文件名行号

	// 使用 lumberjack 实现日志轮转
	lumberjackLogger := &lumberjack.Logger{
		Filename:   config.LogPath,
		MaxSize:    100, // MB
		MaxBackups: 3,
		MaxAge:     28, // days
		Compress:   true,
	}

	// 同时输出到控制台和文件轮转
	mw := io.MultiWriter(os.Stdout, lumberjackLogger)
	log.SetOutput(mw)
	log.SetPrefix("[INFO] ") // 设置默认日志前缀
	log.Println("日志系统初始化成功")
}

// ensureLogDir 创建日志文件所需的目录
func ensureLogDir(filePath string) error {
	dir := filepath.Dir(filePath)
	return os.MkdirAll(dir, 0755)
}

func logWithCaller(level string, format string, args ...interface{}) {
	// 2 表示跳过：
	// runtime.Caller -> logWithCaller -> Info/Warn/Error -> 业务代码
	_, file, line, ok := runtime.Caller(2)

	location := "unknown"
	if ok {
		location = fmt.Sprintf("%s:%d", filepath.Base(file), line)
	}

	log.SetPrefix(fmt.Sprintf("[%s] ", level))
	log.Println(location, "-", fmt.Sprintf(format, args...))
}
func Info(message string, args ...interface{}) {
	logWithCaller("INFO", message, args...)
}

func Warn(message string, args ...interface{}) {
	logWithCaller("WARN", message, args...)
}

func Error(message string, args ...interface{}) {
	logWithCaller("ERROR", message, args...)
}

func Debug(message string, args ...interface{}) {
	logWithCaller("DEBUG", message, args...)
}

func Fatal(message string, args ...interface{}) {
	logWithCaller("FATAL", message, args...)
	os.Exit(1)
}
