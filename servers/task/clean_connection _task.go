/**
* Created by GoLand.
* User: link1st
* Date: 2019-07-31
* Time: 15:17
 */

package task

import (
	"fmt"
	"gowebsocket/servers/websocket"
	"runtime/debug"
	"time"
)

func Init() {
	//注册定时清理连接任务
	Timer(3*time.Second, 30*time.Second, cleanConnection, "", nil, nil)

}

// 清理超时连接
func cleanConnection(param interface{}) (result bool) {
	result = true

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("ClearTimeoutConnections stop", r, string(debug.Stack()))
		}
	}()

	fmt.Println("定时任务，清理超时连接", param)

	websocket.ClearTimeoutConnections()

	return
}
