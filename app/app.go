package app

import (
	"SQLGuardian/cache/db"
	"SQLGuardian/consts"
	"SQLGuardian/job"
	"SQLGuardian/utils"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

/**
  @author: XingGao
  @date: 2023/8/11
**/

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Cron     string
	Database string
}

// InitJob 初始化
func InitJob(systemPort string) {
	host, _ := db.Get([]byte("host"))
	port, _ := db.Get([]byte("port"))
	user, _ := db.Get([]byte("user"))
	cron, _ := db.Get([]byte("cron"))
	password, _ := db.Get([]byte("password"))
	database, _ := db.Get([]byte("database"))
	if cron != nil && port != nil && user != nil && password != nil && database != nil {
		job.Run(string(cron), string(host), string(port), string(user), string(password), string(database))
		log.Printf("please visit %s\n", "http://localhost:"+systemPort)
	} else {
		log.Printf("please config your database in %s\n", "http://localhost:"+systemPort+"/config")
	}

}

// GetBackupDir 返回备份文件夹
func GetBackupDir(writer http.ResponseWriter, request *http.Request) {
	// 获取可执行文件所在目录
	dir := utils.GetExeDir()
	// 获取备份文件夹文件
	files, _ := os.ReadDir(dir + "/" + consts.BackupDir)
	// 返回文件夹文件
	html := `
<html>
<head>
	<style>
		table {
			border-collapse: collapse;
			width: 80%;
			margin: 20px auto;
		}
		th, td {
			border: 1px solid #ddd;
			padding: 8px;
			text-align: left;
		}
		th {
			background-color: #f2f2f2;
		}
		a {
			text-decoration: none;
		}
	</style>
</head>
<body>
	<h1>Backup Files</h1>
`
	html += "<div><a href='/config'>Config Your Database</a></div>"
	html += `
	<table>
		<thead>
			<tr>
				<th>File Name</th>
				<th>File Size</th>
				<th>File Path</th>
				<th>File Time</th>
				<th>Actions</th>
			</tr>
		</thead>
		<tbody>
	`

	for _, file := range files {
		info, _ := file.Info()
		timeStr := info.ModTime().Format("2006-01-02 15:04:05")
		a := "<li>" + file.Name() + "</li>"
		html += `
		<tr>
			<td>` + a + `</td>
			<td>` + strconv.FormatInt(info.Size(), 10) + `</td>
			<td>` + dir + "/" + consts.BackupDir + "/" + file.Name() + `</td>
			<td>` + timeStr + `</td>
			<td><a href="/download?name=` + file.Name() + `">Download</a> | <a href="/delete?name=` + file.Name() + `">Delete</a></td>
		</tr>`
	}

	html += `
		</tbody>
	</table>
</body>
</html>`

	// Now 'html' contains the formatted and styled HTML table with the file list and actions.

	// 写入
	_, _ = writer.Write([]byte(html))
}

// DownloadFile 下载文件
func DownloadFile(writer http.ResponseWriter, request *http.Request) {
	// 获取参数
	name := request.URL.Query().Get("name")
	// 获取工作目录
	dir := utils.GetExeDir()
	// 打开文件
	file, _ := os.Open(dir + "/" + consts.BackupDir + "/" + name)
	// 关闭文件
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)
	// 设置响应头
	writer.Header().Set("Content-Type", "application/octet-stream")
	writer.Header().Set("Content-Disposition", "attachment; filename="+name)
	// 写出
	_, _ = io.Copy(writer, file)
}

// DeleteFile 删除文件
func DeleteFile(writer http.ResponseWriter, request *http.Request) {
	//  获取参数
	name := request.URL.Query().Get("name")
	// 获取工作目录
	dir := utils.GetExeDir()
	// 删除文件
	_ = os.Remove(dir + "/" + consts.BackupDir + "/" + name)
	// 重定向
	http.Redirect(writer, request, "/", http.StatusFound)
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
		html += "<input type='text' name='cron' placeholder='cron' /><br/>"
		html += "<div>备份全部<input type='checkbox' name='switch' value='on' /></div>"
		html += "<input type='text' name='database' placeholder='database' /><br/>"
		html += "<input type='submit' value='submit' />"
		html += "</form></body></html>"

		host, _ := db.Get([]byte("host"))
		port, _ := db.Get([]byte("port"))
		user, _ := db.Get([]byte("user"))
		password, _ := db.Get([]byte("password"))
		cron, _ := db.Get([]byte("cron"))
		database, _ := db.Get([]byte("database"))
		switchOn, _ := db.Get([]byte("switch"))
		if port != nil && user != nil && password != nil {
			if string(database) == "" {
				switchOn = []byte("on")
			}
			html = ""
			html = "<html><body><h1>Config Your Database</h1><form action='/config' method='post'>"
			html += "<input type='text' name='host' placeholder='host' value='" + string(host) + "' /><br/>"
			html += "<input type='text' name='port' placeholder='port' value='" + string(port) + "' /><br/>"
			html += "<input type='text' name='user' placeholder='user' value='" + string(user) + "' /><br/>"
			html += "<input type='text' name='password' placeholder='password' value='" + string(password) + "' /><br/>"
			html += "<input type='text' name='cron' placeholder='cron' value='" + string(cron) + "' /><br/>"
			html += "<input type='text' name='database' placeholder='database' value='" + string(database) + "' /><br/>"
			html += "<div>备份全部<input type='checkbox' name='switch' value='" + string(switchOn) + "' /><div/>"
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
		cron := request.FormValue("cron")
		database := request.FormValue("database")
		switchOn := request.FormValue("switch")
		if switchOn == "on" {
			database = ""
		}
		// 设置缓存
		db.Set([]byte("host"), []byte(host))
		db.Set([]byte("port"), []byte(port))
		db.Set([]byte("user"), []byte(user))
		db.Set([]byte("password"), []byte(password))
		db.Set([]byte("cron"), []byte(cron))
		db.Set([]byte("database"), []byte(database))
		db.Set([]byte("switch"), []byte(switchOn))

		job.Run(cron, host, port, user, password, database)
		// 重定向
		http.Redirect(writer, request, "/", http.StatusFound)
	}
}
