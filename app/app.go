package app

import (
	service2 "SQLGuardian/app/service"
	"SQLGuardian/consts"
	"SQLGuardian/middleware"
	"SQLGuardian/router"
	"SQLGuardian/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kardianos/service"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type YmlConfig struct {
	Server ServerConfig `yaml:"server"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
	Mode string `yaml:"mode"`
}

type program struct{}

func (p *program) Start(s service.Service) error {
	go p.Run()
	return nil
}

func (p *program) Stop(s service.Service) error {
	return nil
}

func (p *program) Run() {
	config, err := loadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	r := gin.Default()
	r.Use(middleware.Cors())
	// 映射静态文件
	r.Static("/static", consts.StaticDir)
	// 加载模板文件
	r.LoadHTMLGlob(consts.ViewDir + "/**")
	// 打印端口
	log.Println("server is running at port " + config.Server.Port)
	// 提示配置数据库
	if err != nil {
		log.Fatal(err)
	}
	router.InitRouters(r)
	service2.InitJob(config.Server.Port)
	err = r.Run(":" + config.Server.Port)
	if err != nil {
		log.Fatal(err)
	}
}

// RegisterService 注册服务 service will install / un-install, start / stop
func RegisterService() {
	// 获取可执行文件所在目录
	dir := utils.GetExeDir()
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

func loadConfig(filename string) (*YmlConfig, error) {
	// 读取文件内容
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// 解析YAML
	var config YmlConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
