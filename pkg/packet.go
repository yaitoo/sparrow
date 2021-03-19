package pkg

import (
	"encoding/binary"

	proto "github.com/golang/protobuf/proto"
)

func createSHeaderPacket(opCode OpCode, trackEnabled TrackEnabled, lenEHeader uint16) []byte {
	buffer := make([]byte, 4)
	buffer[0] = byte(opCode)
	buffer[1] = byte(trackEnabled)
	buffer[2] = byte(lenEHeader)
	buffer[3] = byte(lenEHeader >> 8)
	return buffer
}

func CreateSyncRounting(data *SyncRountingNames) ([]byte, error) {
	syncRountingNames := &SyncRounting{
		RountingNames: data.RountingNames,
		Parameters:    data.Parameters,
		QueryStrings:  data.QueryStrings,
	}

	return proto.Marshal(syncRountingNames)
}

func createEHeader(data *ExternalHeader) ([]byte, error) {
	externalHeader := &EHeader{
		Status:        data.Status,
		PacketID:      data.PacketID,
		TrackID:       data.TrackID,
		PayloadLength: data.PayloadLength,
		Routing:       data.Routing,
		Version:       data.Version,
	}

	return proto.Marshal(externalHeader)
}

func createPacket(sHeader []byte, eHeader []byte, ePayload []byte) []byte {
	lenSHeader := len(sHeader)
	lenEHeader := len(eHeader)
	lenEPayload := len(ePayload)
	lenAll := lenSHeader + lenEHeader + lenEPayload
	buffer := make([]byte, lenAll)
	copy(buffer[:lenSHeader], sHeader)
	copy(buffer[lenSHeader:lenSHeader+lenEHeader], eHeader)
	copy(buffer[lenSHeader+lenEHeader:], ePayload)
	return buffer
}

func (packet *Packet) Init(data []byte) {
	packet.rawBytes = data
	packet.lenEHeader = binary.LittleEndian.Uint16(data[2:4])
}

func (packet *Packet) ReadOpCode() uint8 {
	return packet.rawBytes[0]
}

func (packet *Packet) ReadTrackEnabled() bool {
	return packet.rawBytes[1] == 0x1
}

func (packet *Packet) ReadEHeaderLength() uint16 {
	return packet.lenEHeader
}

func (packet *Packet) ReadEHeaderBuffer() []byte {
	return packet.rawBytes[4 : 4+packet.lenEHeader]
}

func (packet *Packet) ReadEPayloadBuffer() []byte {
	return packet.rawBytes[4+packet.lenEHeader:]

}

func (packet *Packet) ReadEHeaderPB() (*EHeader, error) {
	externalHeader := &EHeader{}
	if err := proto.Unmarshal(packet.rawBytes[4:4+packet.lenEHeader], externalHeader); err != nil {
		return nil, err
	}
	return externalHeader, nil
}
