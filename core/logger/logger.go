/*
 * @Description:日志系统
 * @Author: mali
 * @Date: 2022-09-08 13:37:31
 * @LastEditTime: 2022-09-13 09:19:36
 * @LastEditors: VSCode
 * @Reference:
 */
package logger

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/laravelGo/core/app"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

//Logger 全局的Logger对象
var Logger *zap.Logger

/**
 * @Author: mali
 * @Func:
 * @Description: 初始化日志
 * @Param:
 * @Return:
 * @Example:
 */
func InitLogger(filename string, maxSize, maxBackup, maxAge int, compress bool, logType string, level string) {
	// 获取日志写入介质
	writeSyncer := getLoggerWrite(filename, maxSize, maxBackup, maxAge, compress, logType)
	// 设置日志等级
	logLevel := new(zapcore.Level)
	//检查日志等级是否出错
	if err := logLevel.UnmarshalText([]byte(level)); err != nil {
		fmt.Println("日志初始化错误，日志级别设置有误。请修改 config/log.go 文件中的 log.level 配置项")
	}
	// 初始化 core
	core := zapcore.NewCore(getLoggerEncoder(), writeSyncer, logLevel)

	// 初始化 Logger
	Logger = zap.New(
		core,
		zap.AddCaller(),                   //将调用函数信息记录到日志中 调用文件和行号，必须设置配置对象Config中的CallerKey字段
		zap.AddCallerSkip(1),              // 封装了一层，调用文件去除一层(runtime.Caller(1))
		zap.AddStacktrace(zap.ErrorLevel), // Error 时才会显示 stacktrace
	)
	// 将自定义的 logger 替换为全局的 logger
	// zap.L().Fatal() 调用时，就会使用我们自定的 Logger
	zap.ReplaceGlobals(Logger)
}

/**
 * @Author: mali
 * @Func:
 * @Description: 日志输出的格式
 * @Param:
 * @Return:
 * @Example:
 */
func getLoggerEncoder() zapcore.Encoder {
	// 日志格式名称 指定的这些信息会出现在log 如果不指定对应key的name的话，对应key的信息不处理
	encoderConfig := zapcore.EncoderConfig{
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",                          //日志内容对应的key名，此参数必须不为空，否则日志主体不处理
		TimeKey:        "time",                         //输出时间的key名
		LevelKey:       "level",                        //输出日志级别的key名
		NameKey:        "logger",                       //日志名
		CallerKey:      "caller",                       //代码调用定位 如 paginator/paginator.go:148
		StacktraceKey:  "stacktrace",                   //调用栈追踪的key名
		LineEnding:     zapcore.DefaultLineEnding,      //每行日志的结尾添加 "\n"
		EncodeLevel:    zapcore.CapitalLevelEncoder,    //日志级别名称大写，如 ERROR、INFO
		EncodeTime:     customTimeEncoderfunc,          //时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //执行消耗时间，以秒为单位
		EncodeCaller:   zapcore.ShortCallerEncoder,     //记录调用路径短格式为 package/file:line，长格式为绝对路径
	}
	// 本地环境配置
	if app.IsLocal() {
		// 终端输出的关键词高亮
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		return zapcore.NewConsoleEncoder(encoderConfig)
	}
	//JSON 编码器
	return zapcore.NewJSONEncoder(encoderConfig)
}

/**
 * @Author: mali
 * @Func:
 * @Description: 自定义友好的日志输出时间格式
 * @Param:
 * @Return:
 * @Example:
 * @param {time.Time} t
 * @param {zapcore.PrimitiveArrayEncoder} enc
 */
func customTimeEncoderfunc(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

/**
 * @Author: mali
 * @Func:
 * @Description: 获取日志写入介质
 * @Param:
 * @Return:
 * @Example:
 * @param {string} pathname
 * @param {*} maxSize
 * @param {*} maxBackup
 * @param {int} maxAge
 * @param {bool} compress
 * @param {string} logType
 */
func getLoggerWrite(pathname string, maxSize, maxBackup, maxAge int, compress bool, logType string) zapcore.WriteSyncer {
	//获取文件名
	_, fileName := filepath.Split(pathname)
	//获取文件后缀
	ext := path.Ext(pathname)
	//日志文件记录分类
	if logType == "daily" {
		//按天分类
		logname := time.Now().Format("2006/01/02" + ext)
		pathname = strings.ReplaceAll(pathname, fileName, logname)
	}
	if logType == "month" {
		//按月分类
		logname := time.Now().Format("2006/01" + ext)
		pathname = strings.ReplaceAll(pathname, fileName, logname)
	}
	loggerWrite := &lumberjack.Logger{
		Filename:   pathname,  //文件名
		MaxSize:    maxSize,   //日志单文件的最大占用空间
		MaxAge:     maxAge,    //已经被分割存储的日志文件最大的留存时间，单位是天
		MaxBackups: maxBackup, //分割存储的日志文件最多的留存个数
		Compress:   compress,  //指定被分割之后的文件是否要压缩
	}

	// 配置输出介质
	if app.IsLocal() {
		// 本地开发终端打印和记录文件
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(loggerWrite))
	} else {
		//其他环境只记录文件
		return zapcore.AddSync(loggerWrite)
	}
}
