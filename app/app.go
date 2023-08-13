package app

import (
	"SQLGuardian/cache/appCache"
	"SQLGuardian/consts"
	"SQLGuardian/job"
	"net/http"
	"os"
)

/**
  @author: XingGao
  @date: 2023/8/11
**/

// GetBackupDir 返回备份文件夹
func GetBackupDir(writer http.ResponseWriter, request *http.Request) {
	// 获取工作目录
	dir, _ := os.Getwd()
	// 获取备份文件夹文件
	files, _ := os.ReadDir(dir + "/" + consts.BackupDir)
	// 返回文件夹文件
	html := "<html><body><h1>Backup Files</h1><ul>"
	for _, file := range files {

		a := "<a href='/delete?name=" + file.Name() + "'>" + file.Name() + "</a>"
		// li
		html += "<li>" + a + "</li>"
	}
	html += "</ul></body></html>"
	// 写入
	_, _ = writer.Write([]byte(html))
}

// DeleteFile 删除文件
func DeleteFile(writer http.ResponseWriter, request *http.Request) {
	//  获取参数
	name := request.URL.Query().Get("name")
	// 获取工作目录
	dir, _ := os.Getwd()
	// 删除文件
	_ = os.Remove(dir + "/" + consts.BackupDir + "/" + name)
	// 重定向
	http.Redirect(writer, request, "/tmpfiles", http.StatusFound)
}

// SetConfig 配置数据库
func SetConfig(writer http.ResponseWriter, request *http.Request) {
	// Get
	if request.Method == "GET" {
		// 写出三个输入框
		html := "<html><body><h1>Config Your Database</h1><form action='/config' method='post'>"
		html += "<input type='text' name='host' placeholder='host' /><br/>"
		html += "<input type='text' name='port' placeholder='port' /><br/>"
		html += "<input type='text' name='user' placeholder='user' /><br/>"
		html += "<input type='text' name='password' placeholder='password' /><br/>"
		html += "<input type='text' name='database' placeholder='database' /><br/>"
		html += "<input type='submit' value='submit' />"
		html += "</form></body></html>"

		host, _ := appCache.Get("host")
		port, _ := appCache.Get("port")
		user, _ := appCache.Get("user")
		password, _ := appCache.Get("password")
		database, _ := appCache.Get("database")
		if host != nil && port != nil && user != nil && password != nil && database != nil {
			html = ""
			html = "<html><body><h1>Config Your Database</h1><form action='/config' method='post'>"
			html += "<input type='text' name='host' placeholder='host' value='" + host.(string) + "' /><br/>"
			html += "<input type='text' name='port' placeholder='port' value='" + port.(string) + "' /><br/>"
			html += "<input type='text' name='user' placeholder='user' value='" + user.(string) + "' /><br/>"
			html += "<input type='text' name='password' placeholder='password' value='" + password.(string) + "' /><br/>"
			html += "<input type='text' name='database' placeholder='database' value='" + database.(string) + "' /><br/>"
			html += "<input type='submit' value='submit' />"
			html += "</form></body></html>"
		}
		// css
		html += "<style>input{margin: 5px;}</style>"
		// 写出
		_, _ = writer.Write([]byte(html))
	}
	// Post
	if request.Method == "POST" {
		// 获取参数
		host := request.FormValue("host")
		port := request.FormValue("port")
		user := request.FormValue("user")
		password := request.FormValue("password")
		database := request.FormValue("database")
		// 设置缓存
		appCache.Set("host", host)
		appCache.Set("port", port)
		appCache.Set("user", user)
		appCache.Set("password", password)
		appCache.Set("database", database)

		job.Run(host, port, user, password, database)
		// 重定向
		http.Redirect(writer, request, "/tmpfiles", http.StatusFound)
	}
}
