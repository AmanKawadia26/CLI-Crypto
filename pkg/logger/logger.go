package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

var Logger *zap.SugaredLogger

func init() {
	// Create or open the log file
	logFile, err := os.OpenFile("server.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	// Configure zap to write to both file and console
	writers := zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(logFile))

	// Custom time encoder in RFC3339 format
	timeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(time.RFC3339)) // Example: "2024-09-20T14:55:02Z"
	}

	// Configure encoder with custom time and level encoder
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = timeEncoder                  // Timestamp format
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder // Log level formatting without color
	encoderConfig.TimeKey = "timestamp"                     // Field name for timestamp
	encoderConfig.LevelKey = "level"                        // Field name for log level
	encoderConfig.MessageKey = "message"                    // Field name for the log message
	encoderConfig.CallerKey = "caller"                      // Include the caller file and line number
	encoderConfig.StacktraceKey = "stacktrace"              // Include stacktrace for errors

	// Create the core for logging with JSON encoder
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig), // JSON formatted output
		writers,                               // Log to both console and file
		zapcore.InfoLevel,                     // Log level threshold
	)

	// Initialize zap logger
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	Logger = zapLogger.Sugar()

	// Ensure logs are flushed to file
	defer zapLogger.Sync()
}
