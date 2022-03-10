package main

import (
	"fmt"
	"github.com/loudbund/go-utils/utils_v1"
	_ "github.com/loudbund/go-utils/utils_v1"
)

func main() {
	Env, _ := utils_v1.Config().GetCfgString("app.conf", "main", "env")
	fmt.Println("Env:", Env)

	Port, _ := utils_v1.Config().GetCfgInt("app.conf", "main", "port")
	fmt.Println("Port:", Port)

	if true {
		if utils_v1.File().CheckFileExist("app.conf") {
			fmt.Println("app.conf 文件存在")
		} else {
			fmt.Println("app.conf 文件不存在")
		}
		if utils_v1.File().CheckFileExist("haha.conf") {
			fmt.Println("haha.conf 文件存在")
		} else {
			fmt.Println("haha.conf 文件不存在")
		}
	}

	// Time().SimpleMsgCron
	ExampleSimpleMsgCron()
}

// Time().SimpleMsgCron
func ExampleSimpleMsgCron() {
	var ListenCh = make(chan bool)

	// 1、触发一次执行
	go func() { ListenCh <- true }()

	// 2、启动即时和定时模块
	utils_v1.Time().SimpleMsgCron(ListenCh, 1000*60, func(IsInterval bool) bool {
		// 2.1、处理
		fmt.Print("SimpleMsgCron run event!")
		return true
	})
}
