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

/**
  @author: XingGao
  @date: 2023/8/13
**/

func Run(host string, port string, user string, password string, database string) {
	// MySQL数据库连接信息
	dbConfig := mysql.Config{
		User:                 user,
		Passwd:               password,
		Addr:                 "localhost:" + port,
		Net:                  "tcp",
		DBName:               database,
		AllowNativePasswords: true,
	}
	var c = cron.New()

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
			log.Fatal(err)
		}
	}
	// 清理历史任务
	entre := c.Entries()
	for _, en := range entre {
		c.Remove(en.ID)
	}

	// 每天凌晨执行备份任务
	_, err = c.AddFunc("@every 1m", func() {
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
		backupCmd := fmt.Sprintf("mysqldump -P%s -u%s -p%s %s > %s",
			port,
			dbConfig.User,
			dbConfig.Passwd,
			dbName,
			backupFileName)

		// 执行备份命令
		backupErr := execSystemCommand(backupCmd)
		if backupErr != nil {
			log.Println("备份失败:", backupErr)
			return
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
