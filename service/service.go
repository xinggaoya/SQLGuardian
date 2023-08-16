package service

import (
	"SQLGuardian/app"
	"SQLGuardian/router"
	"fmt"
	"github.com/kardianos/service"
	"log"
	"net/http"
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
	var port = "8000"
	// 打印端口
	fmt.Println("server is running at port " + port)
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

// RegisterService 注册服务
func RegisterService() {
	svcConfig := &service.Config{
		Name:        "SQLGuardian",
		DisplayName: "SQLGuardian",
		Description: "一款简单的MySQL数据库备份工具",
	}

	prg := &program{}

	s, err := service.New(prg, svcConfig)
	if err != nil {
		fmt.Println(err)
	}

	// 获取用户输入
	var choice int
	fmt.Println("Select an operation:")
	fmt.Println("1. install")
	fmt.Println("2. Start")
	fmt.Println("3. stop")
	fmt.Println("4. uninstall")

	fmt.Print("Enter your choice (1/2/3/4): ")
	fmt.Scan(&choice)

	switch choice {
	case 1:
		err = s.Install()
		if err != nil {
			return
		}
	case 2:
		err = s.Start()
		if err != nil {
			return
		}
	case 3:
		err = s.Stop()
		if err != nil {
			return
		}
	case 4:
		err = s.Stop()
		err = s.Uninstall()
		if err != nil {
			return
		}
	default:
		// 运行
		s.Run()
	}

}
