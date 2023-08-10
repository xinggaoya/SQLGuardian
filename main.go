package main

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/robfig/cron/v3"
)

func main() {
	// MySQL数据库连接信息
	dbConfig := mysql.Config{
		User:                 "root",
		Passwd:               "xinggaocode",
		Addr:                 "139.159.194.35:3306", // 修改为你的数据库地址和端口
		Net:                  "tcp",
		DBName:               "db_broad",
		AllowNativePasswords: true,
	}

	// 创建MySQL数据库连接
	db, err := sql.Open("mysql", dbConfig.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 创建定时任务
	c := cron.New()

	// 每天凌晨执行备份任务
	_, err = c.AddFunc("@every 1m", func() {
		// 格式化当前日期作为备份文件名
		backupFileName := fmt.Sprintf("backup-%s.sql", time.Now().Format("2006-01-02"))

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

	// 保持程序运行，直到手动终止
	select {}
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
