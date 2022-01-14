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
}
