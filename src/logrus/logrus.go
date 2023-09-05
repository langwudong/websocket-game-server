package logrus

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

type LogFormatter struct {
}

func (log *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	//定义缓冲区
	var byte *bytes.Buffer
	//entry.Buffer可以将日志消息存储载内存中
	if byte != nil {
		byte = entry.Buffer
	} else {
		byte = &bytes.Buffer{}
	}
	time := entry.Time.Format("2006-01-02 15:04:05")
	//获取日志报告的所在文件路径和行数
	fileVar := fmt.Sprintf("%s : %d", entry.Caller.File, entry.Caller.Line)
	//自定义格式输出到byte
	fmt.Fprintf(byte, "%s [%s] %s %s\n", entry.Level, time, fileVar, entry.Message)
	//使用Bytes方法获取到缓冲区的内容
	return byte.Bytes(), nil
}

func OutPutFile(path string) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Error("打开文件错误......")
		return
	}
	logrus.SetOutput(io.MultiWriter(file, os.Stdout))
}
