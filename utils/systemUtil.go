package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// GetExeDir 获取程序可执行文件所在目录
func GetExeDir() string {
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("获取可执行文件路径时出错:", err)
		return ""
	}

	// 获取可执行文件所在目录
	exeDir := filepath.Dir(exePath)
	return exeDir
}
