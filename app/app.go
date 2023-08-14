package app

import (
	"SQLGuardian/cache/db"
	"SQLGuardian/consts"
	"SQLGuardian/job"
	"fmt"
	"net/http"
	"os"
)

/**
  @author: XingGao
  @date: 2023/8/11
**/

// InitJob 初始化
func InitJob(systemPort string) {
	port, _ := db.Get([]byte("port"))
	user, _ := db.Get([]byte("user"))
	password, _ := db.Get([]byte("password"))
	database, _ := db.Get([]byte("database"))
	if port != nil && user != nil && password != nil && database != nil {
		job.Run("", string(port), string(user), string(password), string(database))
	} else {
		fmt.Printf("please config your database in %s\n", "http://localhost:"+systemPort+"/config")
	}
}

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
		//html += "<input type='text' name='host' placeholder='host' /><br/>"
		html += "<input type='text' name='port' placeholder='port' /><br/>"
		html += "<input type='text' name='user' placeholder='user' /><br/>"
		html += "<input type='text' name='password' placeholder='password' /><br/>"
		html += "<input type='text' name='cron' placeholder='cron' /><br/>"
		html += "<div>备份全部<input type='checkbox' name='switch' value='on' /></div>"
		html += "<input type='text' name='database' placeholder='database' /><br/>"
		html += "<input type='submit' value='submit' />"
		html += "</form></body></html>"

		//host, _ := db.Get([]byte("host"))
		port, _ := db.Get([]byte("port"))
		user, _ := db.Get([]byte("user"))
		password, _ := db.Get([]byte("password"))
		cron, _ := db.Get([]byte("cron"))
		database, _ := db.Get([]byte("database"))
		switchOn, _ := db.Get([]byte("switch"))
		if port != nil && user != nil && password != nil {
			html = ""
			html = "<html><body><h1>Config Your Database</h1><form action='/config' method='post'>"
			//html += "<input type='text' name='host' placeholder='host' value='" + string(host) + "' /><br/>"
			html += "<input type='text' name='port' placeholder='port' value='" + string(port) + "' /><br/>"
			html += "<input type='text' name='user' placeholder='user' value='" + string(user) + "' /><br/>"
			html += "<input type='text' name='password' placeholder='password' value='" + string(password) + "' /><br/>"
			html += "<input type='text' name='cron' placeholder='cron' value='" + string(cron) + "' /><br/>"
			html += "<input type='text' name='database' placeholder='database' value='" + string(database) + "' /><br/>"
			html += "<input type='checkbox' name='switch' value='on' checked='" + string(switchOn) + "' />"
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
		//host := request.FormValue("host")
		port := request.FormValue("port")
		user := request.FormValue("user")
		password := request.FormValue("password")
		cron := request.FormValue("cron")
		database := request.FormValue("database")
		switchOn := request.FormValue("switch")
		if switchOn == "on" {
			database = ""
		}
		// 设置缓存
		//db.Set([]byte("host"), []byte(host))
		db.Set([]byte("port"), []byte(port))
		db.Set([]byte("user"), []byte(user))
		db.Set([]byte("password"), []byte(password))
		db.Set([]byte("cron"), []byte(cron))
		db.Set([]byte("database"), []byte(database))
		db.Set([]byte("switch"), []byte(switchOn))

		job.Run("", port, user, password, database)
		// 重定向
		http.Redirect(writer, request, "/tmpfiles", http.StatusFound)
	}
}
