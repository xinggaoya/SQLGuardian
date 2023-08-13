package router

import (
	"SQLGuardian/app"
	"net/http"
)

/**
  @author: XingGao
  @date: 2023/8/13
**/

func InitRouters() {
	var err error
	// 备份文件列表
	http.HandleFunc("/tmpfiles", app.GetBackupDir)
	// 删除文件
	http.HandleFunc("/delete", app.DeleteFile)
	// 配置数据库
	http.HandleFunc("/config", app.SetConfig)
	// 异常
	if err != nil {
		panic(err)
	}
}
