package service

import (
	"SQLGuardian/app"
	"SQLGuardian/router"
	"fmt"
	"github.com/kardianos/service"
	"log"
	"net/http"
	"os"
)

type program struct{}

func (p *program) Start(s service.Service) error {
	go p.Run()
	return nil
}

func (p *program) Stop(s service.Service) error {
	return nil
}

func (p *program) Run() {
	var err error
	var port = "9210"
	// 打印端口
	log.Println("server is running at port " + port)
	// 提示配置数据库
	if err != nil {
		log.Fatal(err)
	}
	router.InitRouters()
	app.InitJob(port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

// RegisterService 注册服务 service will install / un-install, start / stop
func RegisterService() {
	// 获取工作目录
	dir, _ := os.Getwd()
	svcConfig := &service.Config{
		Name:             "SQLGuardian",
		DisplayName:      "一款简单的MySQL数据库备份工具",
		Description:      "一款简单的MySQL数据库备份工具",
		WorkingDirectory: dir,
	}

	prg := &program{}

	s, err := service.New(prg, svcConfig)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(os.Args) > 1 {
		err = service.Control(s, os.Args[1])
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("operation successful")
		return
	}

	err = s.Run()
	if err != nil {
		fmt.Println(err)
	}

}
