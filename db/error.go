package db

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/go-sql-driver/mysql"
)

var (
	// export variables
	ErrDefaultDatabaseMissing = errors.New("Default database is missing in config file")
	ErrInvalidTransaction     = errors.New("Invalid Transaction operation")
	ErrInvalidObject          = errors.New("Invalid Object")
	ErrTooManySelectedColumns = errors.New("Too Many selected columns")
	ErrClockMovedBackwards    = errors.New("Clock Moved Backwards")
	ErrInvalidCommand         = errors.New("Invalid Command")
	// ================================================= //
	errCrossDatabase = errors.New("cross database")
	//ErrNoDSN                                   = errors.New("no dsn to open database")
	errBadConn                                 = errors.New("bad connection")
	errTxNotFinish                             = errors.New("transaction has not commited or rolled back")
	errTxNotExist                              = errors.New("transaction does not exist")
	errOpenConnError func(string, error) error = func(conn string, err error) error {
		return errors.New(fmt.Sprintf("conn: %s, err: %s", conn, err))
	}
	errDbNotConfigurated  error = errors.New("db not configurated")
	errTableNotConfiugred       = errors.New("table not configurated")

	errLackVariable func(string) error = func(variable string) error {
		return errors.New(fmt.Sprintf("no variable %s", variable))
	}
	errLackTableName func() error = func() error {
		return errors.New("no tablename.")
	}
	errNotStruct                                           = errors.New("Expect a struct")
	errNotPointer                                          = errors.New("Expect a pointer")
	errNilPointer                                          = errors.New("nil pointer")
	errMissField                                           = errors.New("missing field")
	errScanDestOver1Col func(reflect.Kind, []string) error = func(kind reflect.Kind, cols []string) error {
		return fmt.Errorf("scannable dest type %s with >1 columns (%d) in result", kind, len(cols))
	}
	errMissDestName func(string, interface{}) error = func(col string, dest interface{}) error {
		return fmt.Errorf("missing destination name %s in %T", col, dest)
	}
	errNotStructFunc func(reflect.Kind, reflect.Kind) error = func(want, dest reflect.Kind) error {
		return fmt.Errorf("expected %s but got %s", want, dest)
	}
	errStructNoExportedFields func(string) error = func(name string) error {
		return fmt.Errorf("expected a struct, but struct %s has no exported fields", name)
	}
	errNotExpected func(reflect.Kind, reflect.Kind) error = func(want, dest reflect.Kind) error {
		return fmt.Errorf("expected %s but got %s", want, dest)
	}
)

type CommandError struct {
	Command string
	Err     error
}

func NewCommandError(command string, err error) CommandError {
	return CommandError{
		Command: command,
		Err:     err,
	}
}

//ErrKeyExits https://dev.mysql.com/doc/refman/8.0/en/server-error-reference.html#error_er_dup_entry
var ErrKeyExits = errors.New("db: Duplicate entry for key")

//ErrBadData 不符合mysql約束的無效數據
// 1062 : Duplicate entry for PRIMARY/UNIQUE INDEX
// 1048 : Symbol: ER_BAD_NULL_ERROR; SQLSTATE: 23000 Message: Column '%s' cannot be null
// 1406 : Symbol: ER_DATA_TOO_LONG; SQLSTATE: 22001 Message: Data too long for column '%s' at row %ld
var ErrBadData = errors.New("db: Bad data for column")

//IsErr 判斷是不是指定mysql錯誤
func IsErr(err, target error) bool {

	var dbErr *mysql.MySQLError
	if ok := errors.As(err, &dbErr); ok {
		switch target {
		case ErrKeyExits:
			return dbErr.Number == 1062
		case ErrBadData:
			return dbErr.Number == 1048 || dbErr.Number == 1046
		default:
			return errors.Is(err, target)
		}
	}

	return errors.Is(err, target)
}
