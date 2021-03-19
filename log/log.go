package log

import (
	"fmt"
	"log"
	"runtime"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/yaitoo/sparrow/config"
)

//Level log level
type Level int

var (

	//Info info level
	Info Level
	//Warn warn level
	Warn Level = 1
	//Error error level
	Error Level = 2
	//Fatal fatal level
	Fatal Level = 3
	//None disable log output
	None Level = 4
)

var (
	infoString  = color.WhiteString
	warnString  = color.YellowString
	errorString = color.MagentaString
	fatalString = color.RedString
)

const (
	infoPrefix  = "[INFO ]"
	warnPrefix  = "[WARN ]"
	errorPrefix = "[ERROR]"
	fatalPrefix = "[FATAL]"

	timeLayout = "2006/01/02 15:04:05"
)

var (
	//ENV 环境变量前缀，默认空
	ENV = ""
)

//Logger provides a convenient interface for logging
type Logger interface {
	// Printf calls l.Output to print to the logger.
	// Arguments are handled in the manner of fmt.Printf.
	Printf(format string, v ...interface{})

	// Print calls l.Output to print to the logger.
	// Arguments are handled in the manner of fmt.Print.
	Print(v ...interface{})

	// Println calls l.Output to print to the logger.
	// Arguments are handled in the manner of fmt.Println.
	Println(v ...interface{})

	// Warnf calls l.Output to print to the logger.
	// Arguments are handled in the manner of fmt.Printf.
	Warnf(format string, v ...interface{})

	// Warn calls l.Output to print to the logger.
	// Arguments are handled in the manner of fmt.Print.
	Warn(v ...interface{})

	// Warnln calls l.Output to print to the logger.
	// Arguments are handled in the manner of fmt.Println.
	Warnln(v ...interface{})

	// Errorf calls l.Output to print to the logger.
	// Arguments are handled in the manner of fmt.Printf.
	Errorf(format string, v ...interface{})

	// Error calls l.Output to print to the logger.
	// Arguments are handled in the manner of fmt.Print.
	Error(v ...interface{})

	// Errorln calls l.Output to print to the logger.
	// Arguments are handled in the manner of fmt.Println.
	Errorln(v ...interface{})

	// Fatalf calls l.Output to print to the logger.
	// Arguments are handled in the manner of fmt.Printf.
	Fatalf(format string, v ...interface{})

	// Fatal calls l.Output to print to the logger.
	// Arguments are handled in the manner of fmt.Print.
	Fatal(v ...interface{})

	// Fatalln calls l.Output to print to the logger.
	// Arguments are handled in the manner of fmt.Println.
	Fatalln(v ...interface{})

	SetLevel(l Level)
	//OnMessageOutput 消息输出回调
	OnMessageOutput(func(msg string))
}

var cfg *config.Configuration
var err error

func init() {
	cfg, err = config.OpenConfiguration("./conf.d/log.conf")
	if err != nil {
		log.Println("log: ", err.Error())
	}

}

//NewLogger create new named logger
func NewLogger(name string) Logger {

	logger := &logger{
		name:        name,
		level:       Warn,
		infoLogger:  log.New(color.Output, "", 0),
		warnLogger:  log.New(color.Output, "", 0),
		errorLogger: log.New(color.Error, "", 0),
		fatalLogger: log.New(color.Error, "", 0),
	}
	logger.eventMessageOutput = make([]func(msg string), 0, 1)

	logger.SetLevel(parseLevel(cfg.Value(name, "level", "warn")))

	cfg.OnFileChanged(func(c *config.Configuration) {
		logger.SetLevel(parseLevel(c.Value(name, "level", "warn")))
	})

	return logger
}

type logger struct {
	name               string
	level              Level
	infoLogger         *log.Logger
	warnLogger         *log.Logger
	errorLogger        *log.Logger
	fatalLogger        *log.Logger
	eventMessageOutput []func(msg string)
}

func parseLevel(l string) Level {
	switch strings.ToLower(l) {
	case "info":
		return Info
	case "error":
		return Error
	case "warn":
		return Warn
	case "fatal":
		return Fatal
	case "none":
		return None
	default:
		return Error
	}
}

func slog(prefix string, topic string, v ...interface{}) string {
	now := time.Now()
	pc, _, line, _ := runtime.Caller(2)
	if ENV == "" {
		return fmt.Sprintf("%s%s [%s][%s:%d] %s", prefix, now.Format(timeLayout), topic, chopPath(runtime.FuncForPC(pc).Name()), line, fmt.Sprint(v...))
	}

	return fmt.Sprintf("%s%s [%s %s][%s:%d] %s", prefix, now.Format(timeLayout), ENV, topic, chopPath(runtime.FuncForPC(pc).Name()), line, fmt.Sprint(v...))

}

func slogf(prefix string, topic string, format string, v ...interface{}) string {
	now := time.Now()
	pc, _, line, _ := runtime.Caller(2)
	if ENV == "" {
		return fmt.Sprintf("%s%s [%s][%s:%d] %s", prefix, now.Format(timeLayout), topic, chopPath(runtime.FuncForPC(pc).Name()), line, fmt.Sprintf(format, v...))
	}
	return fmt.Sprintf("%s%s [%s %s][%s:%d] %s", prefix, now.Format(timeLayout), ENV, topic, chopPath(runtime.FuncForPC(pc).Name()), line, fmt.Sprintf(format, v...))
}

// return the source filename after the last slash
func chopPath(original string) string {
	i := strings.LastIndex(original, "/")
	if i == -1 {
		return original
	}
	return original[i+1:]
}

func (l *logger) OnMessageOutput(callback func(msg string)) {
	if l == nil {
		return
	}

	if callback != nil {
		l.eventMessageOutput = append(l.eventMessageOutput, callback)
	}
}

func (l *logger) fireMessageOutput(msg string) {
	if l != nil && l.eventMessageOutput != nil {
		for _, callback := range l.eventMessageOutput {
			callback(msg)
		}
	}
}

func (l *logger) Printf(format string, v ...interface{}) {
	if l.level <= Info {
		output := slogf(infoPrefix, l.name, format, v...)
		l.infoLogger.Print(infoString(output))

		go l.fireMessageOutput(output)
	}
}

func (l *logger) Print(v ...interface{}) {
	if l.level <= Info {
		output := slog(infoPrefix, l.name, v...)

		l.infoLogger.Print(infoString(output))

		go l.fireMessageOutput(output)
	}
}

func (l *logger) Println(v ...interface{}) {
	if l.level <= Info {
		output := slog(infoPrefix, l.name, v...)
		l.infoLogger.Println(infoString(output))

		go l.fireMessageOutput(output)
	}
}

func (l *logger) Warnf(format string, v ...interface{}) {
	if l.level <= Warn {
		output := slogf(warnPrefix, l.name, format, v...)

		l.warnLogger.Print(warnString(output))

		go l.fireMessageOutput(output)
	}
}

func (l *logger) Warn(v ...interface{}) {
	if l.level <= Warn {
		output := slog(warnPrefix, l.name, v...)

		l.warnLogger.Print(warnString(output))

		go l.fireMessageOutput(output)
	}
}

func (l *logger) Warnln(v ...interface{}) {
	if l.level <= Warn {
		output := slog(warnPrefix, l.name, v...)
		l.warnLogger.Println(warnString(output))
		go l.fireMessageOutput(output)
	}
}

func (l *logger) Errorf(format string, v ...interface{}) {
	if l.level <= Error {
		output := slogf(errorPrefix, l.name, format, v...)
		l.errorLogger.Print(errorString(output))
		go l.fireMessageOutput(output)
	}
}

func (l *logger) Error(v ...interface{}) {
	if l.level <= Error {
		output := slog(errorPrefix, l.name, v...)
		l.errorLogger.Print(errorString(output))
		go l.fireMessageOutput(output)
	}
}

func (l *logger) Errorln(v ...interface{}) {
	if l.level <= Error {
		output := slog(errorPrefix, l.name, v...)
		l.errorLogger.Println(errorString(output))
		go l.fireMessageOutput(output)
	}
}

func (l *logger) Fatalf(format string, v ...interface{}) {
	if l.level <= Fatal {
		output := slogf(fatalPrefix, l.name, format, v...)
		l.fatalLogger.Print(fatalString(output))
		go l.fireMessageOutput(output)
	}
}

func (l *logger) Fatal(v ...interface{}) {
	if l.level <= Fatal {

		output := slog(fatalPrefix, l.name, v...)

		l.fatalLogger.Print(fatalString(output))
		go l.fireMessageOutput(output)
	}
}

func (l *logger) Fatalln(v ...interface{}) {
	if l.level <= Fatal {
		output := slog(fatalPrefix, l.name, v...)
		l.fatalLogger.Println(fatalString(output))
		go l.fireMessageOutput(output)
	}
}

func (l *logger) SetLevel(minLevel Level) {
	l.level = minLevel
}
