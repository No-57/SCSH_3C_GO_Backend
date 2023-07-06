package Utils

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

// var MainLogger *zap.Logger
// var GatewayLogger *zap.Logger
var LogUtil *zap.Logger

const (
	logTmFmtWithMS = "2006-01-02 15:04:05.000"
)

func init() {

	//MainLogger = NewLogger("logs/main.Utils", zapcore.InfoLevel, 128, 30, 7, true)
	LogUtil = NewLogger("logs/practice.log", zapcore.DebugLevel, 128, 30, 365, true)
	LogUtil.Info("===log initialize success===")
}

/**
 * 獲取日誌
 * filePath 日誌檔案路徑
 * level 日誌級別
 * maxSize 每個日誌檔案儲存的最大尺寸 單位：M
 * maxBackups 日誌檔案最多儲存多少個備份
 * maxAge 檔案最多儲存多少天
 * compress 是否壓縮
 * serviceName 服務名
 */
func NewLogger(filePath string, level zapcore.Level, maxSize int, maxBackups int, maxAge int, compress bool) *zap.Logger {
	core := newCore(filePath, level, maxSize, maxBackups, maxAge, compress)
	return zap.New(core, zap.AddCaller(), zap.Development())
}

/**
 * zapcore構造
 */
func newCore(filePath string, level zapcore.Level, maxSize int, maxBackups int, maxAge int, compress bool) zapcore.Core {

	//日誌檔案路徑配置2
	hook := lumberjack.Logger{
		Filename:   filePath,   // 日誌檔案路徑
		MaxSize:    maxSize,    // 每個日誌檔案儲存的最大尺寸 單位：M
		MaxBackups: maxBackups, // 日誌檔案最多儲存多少個備份
		MaxAge:     maxAge,     // 檔案最多儲存多少天
		Compress:   compress,   // 是否壓縮
	}
	// 設定日誌級別
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(level)
	//公用編碼器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "linenum",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    cEncodeLevel,                   // 小寫編碼器
		EncodeTime:     cEncodeTime,                    // ISO8601 UTC 時間格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   cEncodeCaller,                  // 全路徑編碼器
		EncodeName:     zapcore.FullNameEncoder,
	}
	return zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig), // 編碼器配置
		//zapcore.NewJSONEncoder(encoderConfig),                                           // 編碼器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // 列印到控制檯和檔案
		atomicLevel, // 日誌級別
	)

}

// cEncodeLevel 自定义日志级别显示
func cEncodeLevel(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + level.CapitalString() + "]")
}

// cEncodeTime 自定义时间格式显示
func cEncodeTime(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + t.Format(logTmFmtWithMS) + "]")
}

// cEncodeCaller 自定义行号显示
func cEncodeCaller(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + caller.TrimmedPath() + "]")
}
