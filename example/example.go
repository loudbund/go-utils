package main

import (
	"fmt"
	_ "github.com/loudbund/go-json/json_v1"
	_ "github.com/loudbund/go-mysql/mysql_v1"
	_ "github.com/loudbund/go-pool/pool_v1"
	_ "github.com/loudbund/go-progress/progress_v1"
	_ "github.com/loudbund/go-request/request_v1"
	_ "github.com/loudbund/go-socket/socket_v1"
	"github.com/loudbund/go-utils/utils_v1"
	_ "github.com/loudbund/go-utils/utils_v1"
)

func main() {
	Env, _ := utils_v1.Config().GetCfgString("app.conf", "main", "env")
	fmt.Println("Env:", Env)

	Port, _ := utils_v1.Config().GetCfgInt("app.conf", "main", "port")
	fmt.Println("Port:", Port)
}
