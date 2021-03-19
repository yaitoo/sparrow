package pkg

import (
	"errors"
)

//OpCode 操作代码
type OpCode uint8

type TrackEnabled uint8

type EHPLength uint16

const (
	// Request Request/Repsonse
	REQUEST_RESPONSE OpCode = 0x1
	// Publish Publish/Pubback
	PUBLISH_PUBACK OpCode = 0x2
	// Subscribe Subscribe/Subback
	SUBSCRIBE_SUBACK OpCode = 0x3
	// Unsubscribe Unsubscribe/Unsubback
	UNSUBSCRIBE_UNSUBACK OpCode = 0x4
	// Connect  Connect/Connback
	CONNECT_CONNACK OpCode = 0x6
	// SYNCROUNTING
	SYNCROUNTING OpCode = 0x7
	// Disconnect Disconnect/Disconnback
	DISCONNECT_DISCONNACK OpCode = 0x8
	// PING
	PING OpCode = 0x9
	// PONG
	PONG OpCode = 0xA
)

const (
	// 啟動追蹤
	DISABLED TrackEnabled = 0x0
	// 不啟動追蹤
	ENABLED TrackEnabled = 0x1
)

//SBHeader Standard Binanry Header 标准二进制包头
type SHeader struct {
	//OpCode 操作代码
	OpCode       OpCode
	TrackEnabled TrackEnabled
	EHPLength    EHPLength
}

type ExternalHeader struct {
	Status        uint32
	PacketID      string
	TrackID       string
	PayloadLength uint32
	Routing       string
	Version       string
}

//Packet 数据包
type Packet struct {
	IPacket IPacket
	//SBHeader 标准包头
	rawBytes    []byte
	lenEHeader  uint16
	lenEPayload uint16
}

type SyncRountingNames struct {
	RountingNames string
	// Parameters
	Parameters string
	// QueryStrings
	QueryStrings string
}

var (
	ErrIllegalParseMessage = errors.New("illegal parse message")
)

type IPacket interface {
	Init(data []byte)
	ReadOpCode() uint8
	ReadTrackEnabled() ([]byte, error)
	ReadEHeaderLength() ([]byte, error)
	ReadEHeaderBuffer() ([]byte, error)
	ReadEPayloadBuffer() ([]byte, error)
	ReadEHeaderPB() ([]byte, error)
}
