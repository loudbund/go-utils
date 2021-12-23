# go-utils
整合几个包，引入go-utils模块后，除了可以使用utils的time和config模块，还间接下载了以下几个模块
1. github.com/loudbund/go-json/json_v1
2. github.com/loudbund/go-mysql/mysql_v1
3. github.com/loudbund/go-pool/pool_v1
4. github.com/loudbund/go-progress/progress_v1
5. github.com/loudbund/go-request/request_v1
6. github.com/loudbund/go-socket/socket_v1

## 安装
go get github.com/loudbund/go-utils

## 引入
```golang
import "github.com/loudbund/go-utils/utils_v1"
```

## 使用
```golang

// time模块使用
utils_v1.Time().xxx

// config模块使用
utils_v1.Config().xxx

```