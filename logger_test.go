package logger

import (
	"fmt"
	"testing"
)

func TestLogger(t *testing.T) {
	SetConsole(false)
	//根据配置文件设置日志等级

	SetLevel(DEBUG)
	fmt.Println(INFO)
	//根据配置文件，设置日志路径，日志名，日志切割大小限制
	SetRollingFile("log", "a.log", 10, 1, MB)
	//设置远程告警地址及服务id
	SetRemoteUrl("http://121.41.118.120:8092/api/alert/report")
	SetRemoteServerId("8")
	for i := 0; i < 100000; i++ {
		Log("log start")
		Debug("debug log")
		Info("info log")
		Warn("warn log")
		Error("error log")
		Fatal("fatal log")
		Info("remote", "你好")
	}
}
