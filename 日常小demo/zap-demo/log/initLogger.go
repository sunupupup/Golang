package log

import (
	"io"
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type MyLogger struct {
	*zap.Logger
}

var Logger MyLogger

//构建存放日志的目录
func CreatDir() error {
	_, err := os.Stat("log")
	if os.IsNotExist(err) {
		err = os.Mkdir("log", 0600)
		if err != nil {
			return err
		}
	}

	_, err = os.Stat("log/info")
	if os.IsNotExist(err) {
		err = os.Mkdir("log/info", os.ModePerm)
		if err != nil {
			return err
		}
	}

	_, err = os.Stat("log/warning")
	if os.IsNotExist(err) {
		err = os.Mkdir("log/warning", os.ModePerm)
		if err != nil {
			return err
		}
	}

	_, err = os.Stat("log/error")
	if os.IsNotExist(err) {
		err = os.Mkdir("log/error", os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

func InitLogger(level string) {
	//根据传进来的 level ， 来定logger记录日志的等级
	encoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		MessageKey:    "key",
		CallerKey:     "caller",
		TimeKey:       "timestamp",
		StacktraceKey: "stack",
		EncodeCaller:  zapcore.ShortCallerEncoder,
		EncodeTime: func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendString(t.Format("2006-01-02 15:04:05"))
		},
	})

	if level == "info" {
		hook := getWriter("log/info/info.log")
		core := zapcore.NewTee(
			zapcore.NewCore(encoder, zapcore.AddSync(hook), zapcore.InfoLevel),
			zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.InfoLevel),
		)
		Logger = MyLogger{zap.New(core, zap.AddCaller())}
	}
	if level == "warning" {
		hook := getWriter("log/warning/warning.log")
		core := zapcore.NewTee(
			zapcore.NewCore(encoder, zapcore.AddSync(hook), zapcore.WarnLevel),
			zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.WarnLevel),
		)
		Logger = MyLogger{zap.New(core, zap.AddCaller())}
	}
	if level == "error" {
		hook := getWriter("log/error/error.log")
		core := zapcore.NewTee(
			zapcore.NewCore(encoder, zapcore.AddSync(hook), zapcore.ErrorLevel),
			zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.ErrorLevel),
		)
		Logger = MyLogger{zap.New(core, zap.AddCaller())}
	}

}

func getWriter(filename string) io.Writer {
	return &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    1,     //megabytes
		MaxBackups: 10,    //最多保留的日志文件数
		MaxAge:     7,     //旧log文件最多保存多少天
		Compress:   false, //是否使用 gzip压缩
	}
}

func init() {

	//zap.AddCaller()
	//会在日志中增加这么一条记录 "caller":"caller/main.go:9"   显示函数在哪和行数
	/*
		logger, err := zap.NewProduction(zap.AddCaller())
		if err != nil {
			log.Fatalln(err.Error())
		}
		InfoLogger = MyInfoLogger{logger}

		logger, err = zap.NewProduction(zap.AddCaller())
		if err != nil {
			log.Fatalln(err.Error())
		}
		WarningLogger = MyWarningLogger{logger}

		logger, err = zap.NewProduction(zap.AddCaller())
		if err != nil {
			log.Fatalln(err.Error())
		}
		ErrorLogger = MyErrorLogger{logger}
	*/
}
