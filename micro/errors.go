package micro

//https://blog.csdn.net/m0_37125796/article/details/89447627
//https://go.googlesource.com/proposal/+/master/design/29934-error-values.md
//https://github.com/golang/go/wiki/ErrorValueFAQ

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/yaitoo/sparrow/log"
)

// // Error returns an error representing c and msg.  If c is OK, returns nil.
// func Error(ctx context.Context, c codes.Code, msg ...interface{}) error {
// 	return status.Error(c, fmt.Sprint(msg...))
// }

// // Errorf returns Error(c, fmt.Sprintf(format, a...)).
// func Errorf(ctx context.Context, c codes.Code, format string, a ...interface{}) error {

// 	return status.Errorf(c, format, a...)

// }

var (
	//ThrowLogger 内部错误日志,用于串联内外错误信息
	ThrowLogger = log.NewLogger("throw")

	seededIDGen = rand.New(rand.NewSource(time.Now().UnixNano()))
)

//Error 自定義錯誤對象
type Error struct {
	message string
	err     error
}

func (e *Error) Error() string {
	if len(e.message) > 0 {
		return e.err.Error() + ": " + e.message
	}
	return e.err.Error()
}

//Unwrap 實現errors.Wrapper
func (e *Error) Unwrap() error { return e.err }

func createTraceID() string {
	return strings.ToUpper(strconv.FormatInt(time.Now().Unix(), 36) + "-" + strconv.Itoa(seededIDGen.Intn(1000)))
}

//LogThrow 记录内部错误，返回错误编码
func LogThrow(ctx context.Context, l log.Level, message string) string {
	traceID := createTraceID()
	switch l {
	case log.Info:
		ThrowLogger.Printf("%s: %s\n", traceID, message)
	case log.Warn:
		ThrowLogger.Warnf("%s: %s\n", traceID, message)
	case log.Error:
		ThrowLogger.Errorf("%s: %s\n", traceID, message)
	case log.Fatal:
		ThrowLogger.Fatalf("%s: %s\n", traceID, message)
	default:
		ThrowLogger.Errorf("%s: %s\n", traceID, message)
	}

	return traceID
}

//LogThrowf 记录内部错误，返回错误编码
func LogThrowf(ctx context.Context, l log.Level, format string, args ...interface{}) string {
	traceID := createTraceID()

	switch l {
	case log.Info:
		ThrowLogger.Printf("%s: %s\n", traceID, fmt.Sprintf(format, args...))
	case log.Warn:
		ThrowLogger.Warnf("%s: %s\n", traceID, fmt.Sprintf(format, args...))
	case log.Error:
		ThrowLogger.Errorf("%s: %s\n", traceID, fmt.Sprintf(format, args...))
	case log.Fatal:
		ThrowLogger.Fatalf("%s: %s\n", traceID, fmt.Sprintf(format, args...))
	default:
		ThrowLogger.Errorf("%s: %s\n", traceID, fmt.Sprintf(format, args...))
	}

	return traceID
}

//Throw 創建一個可帶參數和基礎錯誤類型的錯誤對象，支持go2的errors設計規範
//https://github.com/golang/go/wiki/ErrorValueFAQ
func Throw(ctx context.Context, err error, message string) error {
	return &Error{
		message: message,
		err:     err,
	}
}

//Throwf 創建一個可帶參數和基礎錯誤類型的錯誤對象，支持go2的errors設計規範
func Throwf(ctx context.Context, err error, format string, args ...interface{}) error {
	return &Error{
		message: fmt.Sprintf(format, args...),
		err:     err,
	}
}

//IsError 檢測是否自定義錯誤
func IsError(err error) bool {
	_, ok := err.(*Error)
	return ok
}

var (
	//ErrUnknown 无法识别的错误
	ErrUnknown = errors.New("ErrUnknown")
	//ErrBadRequest 客戶端請求參數錯誤
	ErrBadRequest = errors.New("ErrBadRequest")
	//ErrInvalidArgument 無效變數,服務器內部調用，參數不合法
	ErrInvalidArgument = errors.New("ErrInvalidArgument")
	//ErrUnauthorized 未驗證登錄操作
	ErrUnauthorized = errors.New("ErrUnauthorized")
	//ErrForbidden 服務端未授權操作
	ErrForbidden = errors.New("ErrForbidden")
	//ErrMySQL MySQL操作错误
	ErrMySQL = errors.New("ErrMySQL")
	//ErrInternalServerError 服务器内部通用错误
	ErrInternalServerError = errors.New("ErrInternalServerError")
	//ErrNotFound 资源不存在
	ErrNotFound = errors.New("ErrNotFound")
)
