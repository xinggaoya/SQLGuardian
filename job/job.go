package job

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/robfig/cron/v3"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"
)

/**
  @author: XingGao
  @date: 2023/8/13
**/

func Run(host string, port string, user string, password string, database string) {
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
	backupDir := fmt.Sprintf("%s/backup", dir)
	// 创建备份文件夹
	if _, err = os.Stat(backupDir); os.IsNotExist(err) {
		err = os.Mkdir(backupDir, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}
	// 创建定时任务
	c := cron.New()

	// 每天凌晨执行备份任务
	_, err = c.AddFunc("@every 1m", func() {
		// 格式化当前日期作为备份文件名
		backupFileName := fmt.Sprintf("%s/backup_%s.sql", backupDir, time.Now().Format("20060102150405"))

		// 执行备份命令
		backupCmd := fmt.Sprintf("mysqldump -u%s -h%s -p%s --all-databases > %s",
			dbConfig.User,
			dbConfig.Addr,
			dbConfig.Passwd,
			backupFileName)

		// 执行备份命令
		backupErr := execSystemCommand(backupCmd)
		if backupErr != nil {
			log.Println("备份失败:", backupErr)
			return
		}

		log.Println("备份成功:", backupFileName)
	})

	if err != nil {
		log.Fatal(err)
	}

	// 启动定时任务
	c.Start()
}

// 根据操作系统执行Shell命令的辅助函数
func execSystemCommand(cmd string) error {
	var command *exec.Cmd

	if runtime.GOOS == "windows" {
		command = exec.Command("cmd", "/C", cmd)
	} else {
		command = exec.Command("bash", "-c", cmd)
	}

	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command.Run()
}
