package tasksAndEvents

import (
	"bytes"
	"edu_api/config"
	"net"
	"time"
)

/**
定义一个转换消息体的接口
*/
type Operate interface {
	ToBytes() ([]byte, error)
}

/**
数据库操作
*/
func operateDB(operate Operate) (int, error) {
	conn, err := net.DialTimeout("tcp", config.Config.Queue.Addr, 200*time.Millisecond)
	if err != nil {
		return 0, err
	}

	defer conn.Close()

	data, err := operate.ToBytes()
	if err != nil {
		return 0, err
	}

	var buffer bytes.Buffer
	buffer.Write(data)
	buffer.WriteString("\n") //这个必须，因为php swoole默认是以\n结束消息的

	return conn.Write(buffer.Bytes())
}
