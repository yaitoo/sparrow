package pkg

//PackPacket 根据EPHeader， 使用protobuf格式组装业务数据
func PackPacket(opCode OpCode, trackEnabled TrackEnabled, data *ExternalHeader, bufferEPayload []byte) ([]byte, error) {
	bufferEHeader, err := createEHeader(data)
	if err != nil {
		return nil, err
	}
	lenEHeader := uint16(len(bufferEHeader))
	bufferSHeader := createSHeaderPacket(opCode, trackEnabled, lenEHeader)
	return createPacket(bufferSHeader, bufferEHeader, bufferEPayload), nil
}

func PackEHeader(data *ExternalHeader) ([]byte, error) {
	return createEHeader(data)
}

//Unpack
func Unpack(buf []byte) *Packet {
	packet := &Packet{}
	packet.Init(buf)
	return packet
}
