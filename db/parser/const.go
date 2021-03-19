package parser

const defaultBufSize = 4096
const maxAllowedPacket = defaultMaxAllowedPacket
const digits10 = "0000000000111111111122222222223333333333444444444455555555556666666666777777777788888888889999999999"
const digits01 = "0123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789"

const (
	// defaultAuthPlugin       = "mysql_native_password"
	defaultMaxAllowedPacket = 4 << 20 // 4 MiB
	// minProtocolVersion      = 10
/* 	maxPacketSize           = 1<<24 - 1
timeFormat              = "2006-01-02 15:04:05.999999" */
)

// http://dev.mysql.com/doc/internals/en/status-flags.html
type statusFlag uint16

const (
	statusInTrans statusFlag = 1 << iota
	statusInAutocommit
	statusReserved // Not in documentation
	statusMoreResultsExists
	statusNoGoodIndexUsed
	statusNoIndexUsed
	statusCursorExists
	statusLastRowSent
	statusDbDropped
	statusNoBackslashEscapes
	statusMetadataChanged
	statusQueryWasSlow
	statusPsOutParams
	statusInTransReadonly
	statusSessionStateChanged
)
