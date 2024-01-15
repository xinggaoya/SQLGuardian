package app

import (
	"SQLGuardian/cache/boltdb"
	"SQLGuardian/consts"
	"SQLGuardian/job"
	"SQLGuardian/utils"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
)

/**
  @author: XingGao
  @date: 2023/8/11
**/

type Config struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Cron     string `json:"cron"`
	Database string `json:"database"`
}

// InitJob 初始化
func InitJob(systemPort string) {
	var list []Config

	boltDb := boltdb.NewBoltDB()
	data, err := boltDb.QueryData(consts.CronListKey)
	if err != nil {
		return
	}
	_ = json.Unmarshal(data, &list)

	if len(list) == 0 {
		log.Printf("please config your database in %s\n", "http://localhost:"+systemPort+"/config")
		return
	}

	for _, item := range list {
		if item.Host != "" && item.Port != "" && item.User != "" && item.Password != "" && item.Cron != "" {
			job.Run(item.Cron, item.Host, item.Port, item.User, item.Password, item.Database)
		}
	}
	log.Printf("please visit %s\n", "http://localhost:"+systemPort)
}

// GetBackupFilesList 获取备份文件信息列表
func GetBackupFilesList(c *gin.Context) {
	// 获取可执行文件所在目录
	dir := utils.GetExeDir()
	// 获取备份文件夹文件
	files, _ := os.ReadDir(dir + "/" + consts.BackupDir)
	type fileInfo struct {
		Name string `json:"name"`
		Size int64  `json:"size"`
		Time string `json:"time"`
	}
	var fileInfoList []fileInfo
	for _, file := range files {
		info, _ := file.Info()
		fileInfoList = append(fileInfoList, fileInfo{
			Name: info.Name(),
			Size: info.Size(),
			Time: info.ModTime().Format("2006-01-02 15:04:05"),
		})
	}

	// 写入
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": fileInfoList,
	})
}

// DownloadFile 下载文件
func DownloadFile(c *gin.Context) {
	// 获取参数
	name := c.Query("name")
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
	c.Writer.Header().Add("Content-Disposition", "attachment; filename="+name)
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	// 写出
	_, _ = io.Copy(c.Writer, file)
}

// DeleteFile 删除文件
func DeleteFile(c *gin.Context) {
	//  获取参数
	name := c.Query("name")
	// 获取工作目录
	dir := utils.GetExeDir()
	// 删除文件
	_ = os.Remove(dir + "/" + consts.BackupDir + "/" + name)
	// 重定向
	http.Redirect(c.Writer, c.Request, "/", http.StatusFound)
}

// SetConfig 配置数据库
func SetConfig(c *gin.Context) {
	// Get
	//if c.Request.Method == "GET" {
	//	// 写出三个输入框
	//	html := "<html><body><h1>Config Your Database</h1><form action='/config' method='post'>"
	//	html += "<input type='text' name='host' placeholder='host' /><br/>"
	//	html += "<input type='text' name='port' placeholder='port' /><br/>"
	//	html += "<input type='text' name='user' placeholder='user' /><br/>"
	//	html += "<input type='text' name='password' placeholder='password' /><br/>"
	//	html += "<input type='text' name='cron' placeholder='cron' /><br/>"
	//	html += "<div>备份全部<input type='checkbox' name='switch' value='on' /></div>"
	//	html += "<input type='text' name='database' placeholder='database' /><br/>"
	//	html += "<input type='submit' value='submit' />"
	//	html += "</form></body></html>"
	//
	//	host, _ := db.Get([]byte("host"))
	//	port, _ := db.Get([]byte("port"))
	//	user, _ := db.Get([]byte("user"))
	//	password, _ := db.Get([]byte("password"))
	//	cron, _ := db.Get([]byte("cron"))
	//	database, _ := db.Get([]byte("database"))
	//	switchOn, _ := db.Get([]byte("switch"))
	//	if port != nil && user != nil && password != nil {
	//		if string(database) == "" {
	//			switchOn = []byte("on")
	//		}
	//		html = ""
	//		html = "<html><body><h1>Config Your Database</h1><form action='/config' method='post'>"
	//		html += "<input type='text' name='host' placeholder='host' value='" + string(host) + "' /><br/>"
	//		html += "<input type='text' name='port' placeholder='port' value='" + string(port) + "' /><br/>"
	//		html += "<input type='text' name='user' placeholder='user' value='" + string(user) + "' /><br/>"
	//		html += "<input type='text' name='password' placeholder='password' value='" + string(password) + "' /><br/>"
	//		html += "<input type='text' name='cron' placeholder='cron' value='" + string(cron) + "' /><br/>"
	//		html += "<input type='text' name='database' placeholder='database' value='" + string(database) + "' /><br/>"
	//		html += "<div>备份全部<input type='checkbox' name='switch' value='" + string(switchOn) + "' /><div/>"
	//		html += "<input type='submit' value='submit' />"
	//		html += "</form></body></html>"
	//	}
	//	// css
	//	html += "<style>input{margin: 5px;}</style>"
	//	// 写出
	//	_, _ = writer.Write([]byte(html))
	//}
	//// Post
	//if c.Request.Method == "POST" {
	//	// 获取参数
	//	host := request.FormValue("host")
	//	port := request.FormValue("port")
	//	user := request.FormValue("user")
	//	password := request.FormValue("password")
	//	cron := request.FormValue("cron")
	//	database := request.FormValue("database")
	//	switchOn := request.FormValue("switch")
	//	if switchOn == "on" {
	//		database = ""
	//	}
	//	// 设置缓存
	//	db.Set([]byte("host"), []byte(host))
	//	db.Set([]byte("port"), []byte(port))
	//	db.Set([]byte("user"), []byte(user))
	//	db.Set([]byte("password"), []byte(password))
	//	db.Set([]byte("cron"), []byte(cron))
	//	db.Set([]byte("database"), []byte(database))
	//	db.Set([]byte("switch"), []byte(switchOn))
	//
	//	job.Run(cron, host, port, user, password, database)
	//	// 重定向
	//	http.Redirect(writer, request, "/", http.StatusFound)
	//}
}
