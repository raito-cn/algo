package util

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"sync"
	"time"
)

var log *zap.Logger
var once sync.Once

func initLogger() {
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:       "time",  // 时间字段名
		LevelKey:      "level", // 日志等级字段名
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalColorLevelEncoder, // 彩色等级
		EncodeCaller:  zapcore.FullCallerEncoder,        // 文件:行号
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			//  自定义时间格式
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderCfg),
		zapcore.AddSync(os.Stdout),
		zapcore.DebugLevel,
	)

	log = zap.New(core, zap.AddCaller())
}

func GetLog() *zap.Logger {
	once.Do(initLogger)
	return log
}
