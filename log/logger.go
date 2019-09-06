package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

//全局自定义logger
var root *logrus.Logger

//写入文件拿到文件writer

func logToFile(fileName string) io.Writer{
    f, err := os.OpenFile(fileName, os.O_CREATE | os.O_APPEND | os.O_WRONLY, 0777)

    //如果有error则直接终止整个程序
    if err != nil{
    	panic(err)
	}

    return f
}

//初始化logger

func LogTo(target string, level string){
	//初始化
	root = logrus.New()

	//配置level
	switch level {
	case "DEBUG":
		root.SetLevel(logrus.DebugLevel)
	case "INFO":
		root.SetLevel(logrus.InfoLevel)
	case "WARN":
		root.SetLevel(logrus.WarnLevel)
	case "ERROR":
		root.SetLevel(logrus.ErrorLevel)
	default:
		Warn("Expected level [DEBUG, INFO, WARN, ERROR] ")
	}

    //输出日志到目标文件

	switch target {
	case "stdout":
		root.SetOutput(os.Stdout)
	default:
		// output to file target
        root.SetOutput(logToFile(target))
	}
}


//默认输出到控制台
func LogToStdout(level string){
    LogTo("stdout", level)
}



type Logger interface {

	//增加日志前缀
	AddPrefix(prefixes ...string)

	//清除日志前缀
	//
	ClearPrefix()

	//提供四种日志级别
	Debug(format string, args ... interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})

}



//带有自定义前缀的日志
type PrefixedLogger struct {
	prefix string
}

func (pfl *PrefixedLogger) AddPrefix(prefixes ...string){
	for _, prefix := range prefixes{
        pfl.prefix += fmt.Sprintf("[%s]", prefix)
	}
}

func (pfl *PrefixedLogger) ClearPrefix(){
	pfl.prefix = ""
}

func (pfl *PrefixedLogger) withPrefix(format string) string{
	return fmt.Sprintf("%s%s", pfl.prefix, format)
}

func (pfl *PrefixedLogger) Debug(format string, args ...interface{}){
	root.Debugf(pfl.withPrefix(format), args...)
}

func (pfl *PrefixedLogger) Info(format string, args ...interface{}){
	root.Infof(pfl.withPrefix(format), args...)
}

func (pfl *PrefixedLogger) Warn(format string, args ...interface{}){
	root.Warnf(pfl.withPrefix(format), args...)
}

func (pfl *PrefixedLogger) Error(format string, args ...interface{}){
	root.Errorf(pfl.withPrefix(format), args...)
}


//

func NewPrefixedLogger(prefixes ...string) *PrefixedLogger{
	pfl := &PrefixedLogger{}
	pfl.AddPrefix(prefixes...)
	return pfl
}

//不带前缀的全局日志函数
func Debug(format string, args ...interface{}){
	root.Debugf(format, args...)
}

func Info(format string, args ...interface{}){
	root.Infof(format, args...)
}

func Warn(format string, args ...interface{}){
	root.Warnf(format, args...)
}

func Error(format string, args ...interface{}){
	root.Errorf(format, args...)
}