package cmn

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"time"

	"github.com/fatih/color"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogStruct struct {
	Writer    io.Writer
	File      *os.File
	Print_cfg bool // 此条打印到控制台
	Separator string
}

type LogFileld map[string]string

var (
	LOG_DEBUG   = "Debug"
	LOG_ERROR   = "Error"
	LOG_Info    = "Info"
	LOG_WARNING = "Warning"
)

// 日志颜色
var colors = map[string]func(a ...interface{}) string{
	"Warning": color.New(color.FgYellow).Add(color.Bold).SprintFunc(),
	"Panic":   color.New(color.BgRed).Add(color.Bold).SprintFunc(),
	"Error":   color.New(color.FgRed).Add(color.Bold).SprintFunc(),
	"Info":    color.New(color.FgCyan).Add(color.Bold).SprintFunc(),
	"Debug":   color.New(color.FgWhite).Add(color.Bold).SprintFunc(),
}

// 不同级别前缀与时间的间隔，保持宽度一致
var spaces = map[string]string{
	"Warning": "",
	"Panic":   "  ",
	"Error":   "  ",
	"Info":    "   ",
	"Debug":   "  ",
}

// 运行日志静态类
var runLogStatic = LogStruct{}

// 控制台打印
// Warning Panic Error Info Debug
func Pln(prefix string, msg string) {
	fmt.Printf(
		"%s%s %s %s\n",
		colors[prefix]("["+prefix+"]"),
		spaces[prefix],
		time.Now().Format(TimeFormatMode1),
		msg,
	)
}

// 控制台打印，支持颜色
func Print(color, key, msg string) {
	fmt.Printf(
		"%s%s %s\n",
		colors[color](key),
		time.Now().Format(TimeFormatMode1),
		msg,
	)
}

// func Debug(a ...interface{}) {
// 	fmt.Print("[Debug] ")
// 	fmt.Println(a...)
// }

// // 错误并退出
// func ErrorExit(err_title, err_msg string) {
// 	newLog := NewLog("err.log")

// 	Pln("Error", err_title+err_msg)
// 	newLog.Error(err_title, err_msg)
// 	os.Exit(1)
// }

// 写入日志的文件
func NewLog(log_file_name string) *LogStruct {
	logStruct := &LogStruct{}
	logStruct.Separator = ""
	logDir := path.Dir(log_file_name)
	ok, _ := PathExists(logDir)
	if !ok {
		if err := os.MkdirAll(logDir, 0666); err != nil {
			fmt.Println("创建日志文件错误", err.Error())
		}
	}
	_, err := os.Stat(log_file_name)
	if err != nil {
		f, _ := os.Create(log_file_name)
		logStruct.File = f
		logStruct.Writer = io.MultiWriter(f)
	} else {
		f, _ := os.OpenFile(log_file_name, os.O_APPEND|os.O_WRONLY, 0666)
		logStruct.File = f
		logStruct.Writer = io.MultiWriter(f)
	}
	return logStruct
}

// 运行日志直接静态
func RunLog() *LogStruct {
	// 按小时/日/月/年
	// 先判断文件（夹）是否存在。否多级创建
	log_file := "res/runtime/log/"
	ok, _ := PathExists(log_file)
	if !ok {
		os.MkdirAll(log_file, 0777)
	}
	log_file_name := log_file + time.Unix(time.Now().Unix(), 1).Format("2006-01-02") + ".log"
	_, err := os.Stat(log_file_name)
	runLogStatic.Separator = "|"
	if err != nil {
		f, _ := os.Create(log_file_name)
		runLogStatic.File = f
		runLogStatic.Writer = io.MultiWriter(f)
	} else {
		if runLogStatic.File == nil {
			f, _ := os.OpenFile(log_file_name, os.O_APPEND|os.O_WRONLY, 0666)
			runLogStatic.File = f
			runLogStatic.Writer = io.MultiWriter(f)
		}
	}
	return &runLogStatic
}

func (t *LogStruct) Write(content string) (n int, err error) {

	return io.WriteString(t.Writer, content)
}

func (t *LogStruct) Format(log_type string, content string) (n int, err error) {
	content = log_type + spaces[log_type] + " " + GetTime() + " " + content + "\n"
	return t.Write(content)
}

func (t *LogStruct) Info(content ...string) (n int, err error) {
	str := ""
	for i := 0; i < len(content); i++ {
		if i != 0 {
			str += t.Separator + content[i]
		} else {
			str += content[i]
		}
	}
	n, err = t.Format("Info", str)
	if t.Print_cfg == true {
		Pln("Info", str)
		t.Print_cfg = false
	}
	return
}

func (t *LogStruct) Debug(content string) {
	t.Format("Debug", content)
	if t.Print_cfg == true {
		Pln("Debug", content)
		t.Print_cfg = false
	}
}

func (t *LogStruct) Error(content ...string) {
	content_str := ""
	for i := 0; i < len(content); i++ {
		if i != 0 {
			content_str += t.Separator + content[i]
		} else {
			content_str += content[i]
		}
	}
	t.Format("Error", content_str)
	if t.Print_cfg == true {
		Pln("Error", content_str)
		t.Print_cfg = false
	}
}

// // 打印错误
// func (t *LogStruct) ErrorPrint(key, value string) {
// 	t.Print_cfg = true
// 	t.Error(key, value)
// }

// // 打印Debug
// func (t *LogStruct) DebugPrint(key, value string) {
// 	t.Print_cfg = true
// 	content := key + " " + value
// 	t.Debug(content)
// }

// func (t *LogStruct) Print() *LogStruct {
// 	t.Print_cfg = true
// 	return t
// }

// func (t *LogStruct) FormatFileld(field LogFileld) string {
// 	str := ""
// 	for k, v := range field {
// 		str += k + ":\"" + v + "\"" + t.Separator
// 	}
// 	if len(str) != 0 {
// 		str = str[0 : len(str)-1]
// 	}
// 	return str
// }

// TODO(GgoCoder) 日志轮转
func InitLogger(fileName string, level zapcore.LevelEnabler) *zap.SugaredLogger {
	fileWriteSyncer := getLogWriter(fileName)
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(fileWriteSyncer, zapcore.AddSync(os.Stdout)), level)
	logger := zap.New(core, zap.AddCaller())
	return logger.Sugar()
}

func getEncoder() zapcore.Encoder {
	logConf := zap.NewProductionEncoderConfig()
	logConf.EncodeTime = zapcore.ISO8601TimeEncoder
	logConf.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(logConf)
}

func getLogWriter(fileName string) zapcore.WriteSyncer {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		log.Panic("failed to create log file", fileName)
	}
	return zapcore.AddSync(file)
}
