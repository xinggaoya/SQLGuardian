package job

import (
	"SQLGuardian/consts"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/robfig/cron/v3"
	"log"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"time"
)

/*
*

	@author: XingGao
	@date: 2023/8/13

*
*/
var c = cron.New()

func Run(cronStr string, host string, port string, user string, password string, database string) {
	if host == "" {
		host = "localhost"
	}
	// MySQL数据库连接信息
	dbConfig := mysql.Config{
		User:                 user,
		Passwd:               password,
		Addr:                 host + ":" + port,
		Net:                  "tcp",
		DBName:               database,
		AllowNativePasswords: true,
	}

	// 创建MySQL数据库连接  用于测试
	db, err := sql.Open("mysql", dbConfig.FormatDSN())

	if err != nil {
		log.Fatal(err)
	}
	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)
	// 获取当前程序运行的目录
	dir, err := os.Getwd()
	//备份文件夹
	backupDir := fmt.Sprintf("%s/"+consts.BackupDir, dir)
	// 创建备份文件夹
	if _, err = os.Stat(backupDir); os.IsNotExist(err) {
		err = os.Mkdir(backupDir, os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
	}
	// 清理历史任务
	entre := c.Entries()
	for _, en := range entre {
		c.Remove(en.ID)
	}

	// 每天凌晨执行备份任务
	_, err = c.AddFunc(cronStr, func() {
		if database == "" {
			database = "all"
		}
		// 备份文件名 使用时间戳
		backupFileName := fmt.Sprintf("%s/%s_%s.sql", backupDir, database, time.Now().Format("20060102150405"))
		dbName := database
		if database == "all" {
			dbName = "--all-databases"
		} else {
			dbName = "--databases " + dbName
		}
		// 执行备份命令
		backupCmd := fmt.Sprintf(backupDir+"/static/mysqldump.exe -h %s -P%s -u%s -p%s %s > %s",
			host,
			port,
			dbConfig.User,
			dbConfig.Passwd,
			dbName,
			backupFileName)

		// 执行备份命令
		fmt.Printf("备份命令:%s \n", backupCmd)
		output, backupErr := execSystemCommand(backupCmd)
		if backupErr != nil {
			fmt.Println("备份失败:", backupErr)
			// 写入err.log日志
			errPath := fmt.Sprintf("err.log")
			// 内容
			content := fmt.Sprintf("%s 备份失败,error:%s \n", time.Now().Format("2006-01-02 15:04:05"), output)
			writeLog(errPath, content)
			return
		} else {
			// 写入成功日志
			successPath := fmt.Sprintf("success.log")
			// 内容
			content := fmt.Sprintf("%s 备份成功,备份路径:%s \n", time.Now().Format("2006-01-02 15:04:05"), output)
			writeLog(successPath, content)
		}

		// 保留最近的5个备份文件
		// 获取备份文件列表
		backupFiles, err := os.ReadDir(backupDir)
		if err != nil {
			fmt.Println(err)
		}
		// 文件名排序
		sort.Slice(backupFiles, func(i, j int) bool {
			return backupFiles[i].Name() > backupFiles[j].Name()
		})

		// 删除多余的备份文件
		for i := 5; i < len(backupFiles); i++ {
			_ = os.Remove(fmt.Sprintf("%s/%s", backupDir, backupFiles[i].Name()))
		}

		log.Println("备份成功,备份路径:", backupFileName)
	})

	if err != nil {
		fmt.Printf("cron failed, err:%v\n", err)
	}

	// 启动定时任务
	c.Start()
}

// 根据操作系统执行Shell命令的辅助函数
func execSystemCommand(cmd string) ([]byte, error) {
	var command *exec.Cmd

	// 使用管理员权限执行命令
	if runtime.GOOS == "windows" {
		command = exec.Command("cmd", "/C", cmd)
	} else {
		command = exec.Command("bash", "-c", cmd)
	}
	return command.CombinedOutput()
}

// 日志文件写入
func writeLog(logPath string, content string) {
	// 获取当前程序运行的目录
	dir, err := os.Getwd()
	path := fmt.Sprintf("%s/log", dir)
	// 检查文件夹是否存在
	if _, err = os.Stat(path); os.IsNotExist(err) {
		err = os.Mkdir(path, os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
	}
	// 打开文件
	file := fmt.Sprintf("%s/%s", path, logPath)
	f, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}
	// 写入文件
	_, err = f.WriteString(content)
	if err != nil {
		fmt.Println(err)
	}
	// 关闭文件
	err = f.Close()
	if err != nil {
		fmt.Println(err)
	}
}
