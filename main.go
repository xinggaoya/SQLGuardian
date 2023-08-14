package main

import (
	"SQLGuardian/app"
	"SQLGuardian/router"
	"fmt"
	"log"
	"net/http"
)

func main() {
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
