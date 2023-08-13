package main

import (
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
	fmt.Printf("please config your database in %s\n", "http://localhost:"+port+"/config")
	if err != nil {
		log.Fatal(err)
	}
	router.InitRouters()
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
