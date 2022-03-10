package utils_v1

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
)

// 结构体1:
type uTime struct {
}

// 对外函数1：
func Time() *uTime {
	return &uTime{}
}

/**-------------------------
// 名称：获取当前时间字串
// 参数: 无
// 返回："2021-08-25 11:16:20"
***-----------------------*/
func (t *uTime) DateTime(T ...time.Time) string {
	timeObj := time.Now()
	if len(T) > 0 {
		timeObj = T[0]
	}
	return timeObj.Format("2006-01-02 15:04:05")
}

/**-------------------------
// 名称：获取当前日期字串
// 参数: 无
// 返回："2021-08-25"
***-----------------------*/
func (t *uTime) Date(T ...time.Time) string {
	timeObj := time.Now()
	if len(T) > 0 {
		timeObj = T[0]
	}
	return timeObj.Format("2006-01-02")
}

/**-------------------------
// 名称：日期调整
// 参数: date:"2021-08-25"
//       dayNum:1/-1
// 返回："2021-08-25"
***-----------------------*/
func (t *uTime) DateAdd(date string, dayNum int) string {
	dateObj := t.Date2Time(date)
	// 时间调整
	d, _ := time.ParseDuration(fmt.Sprintf("%dh", 24*dayNum))
	timeObj := dateObj.Add(d)
	return timeObj.Format("2006-01-02")
}

/**-------------------------
// 名称：获取当前日期字串
// 参数: "2021-08-25"
// 返回：Time
***-----------------------*/
func (t *uTime) Date2Time(date string) time.Time {
	// 日期参数转成时间
	at, _ := time.ParseInLocation("2006-01-02", date, time.Local)
	return at
}

// 简易即时定时封装
// SimpleMsgCron ( &ListenEvent, 2000,  func(IsInterval bool) bool{return true})
// ListenEvent: 监听的通道;
// cronMill:定时触发时间(毫秒);
// fEvents:触发事件回调函数；参数为触发类型，一般用不上的；函数返回true则继续监听，否则退出监听
// -------------------------------------------
// 示例：
// func ExampleSimpleMsgCron() {
// 	var ListenCh = make(chan bool)
// 	// 1、触发一次执行
// 	go func() { ListenCh <- true }()
// 	// 2、启动即时和定时模块
// 	utils_v1.Time().SimpleMsgCron(ListenCh, 1000*60, func(IsInterval bool) bool {
// 		// 2.1、处理
// 		fmt.Print("SimpleMsgCron run event!")
// 		return true
// 	})
// }
func (t *uTime) SimpleMsgCron(ExListenCh chan bool, cronMill int, fEvent func(IsInterval bool) bool) {

	var (
		ChCache      = make(chan bool, 5) // 增加5次的消息缓冲归并功能【理论上只需要1个的】
		ChInnerGoEnd = make(chan bool)    // 退出前通知内部goroutine退出用
		hasEvent     = false              // 有消息标志变量
	)
	defer func() { close(ChInnerGoEnd); fmt.Println("SimpleMsgCron defer") }()

	// 循环接收通知消息goroutine
	go func() {
		defer fmt.Println("SimpleMsgCron Inner func defer")
		for {
			select {
			case <-ChInnerGoEnd: // 退出内部goroutine
				return
			case _, ok := <-ExListenCh: // 外部传来的消息
				if !ok {
					log.Error("SimpleMsgCron ExListenCh通道被关闭，停止接收ExListenCh事件")
					return
				}
				// 只有当消息都处理了时才通知处理
				if !hasEvent {
					hasEvent = true
					// 发给下一个通道
					ChCache <- true
				}
			}
		}
	}()

	// 即时和定时处理响应
	timerA := time.NewTicker(time.Millisecond * time.Duration(cronMill))
	for {
		select {
		case <-ChCache:
			// 只有消息标记为true的时候才处理，
			if hasEvent {
				hasEvent = false
				// 事件处理回调，回调函数返回为true，则继续等待下一条消息，否则退出消息监听
				if !fEvent(true) {
					return
				}
			} else {
				// fmt.Println("ignore")
			}
		case <-timerA.C:
			hasEvent = false
			// 定时时间到了执行回调函数，回调函数返回为true，则继续等待下一条消息，否则退出消息监听
			if !fEvent(true) {
				return
			}
		}
	}
}
