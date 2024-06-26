package fitprotocol

import (
	"encoding/binary"
	"errors"
)

const (
	FIT_FILE_TYPE = ".FIT"
)

var INVALID_HEADER_LEN = errors.New("Provided header not correct byte length.")
var INVALID_HEADER_TYPE = errors.New("Provided header not .fit type.")

type Header struct {
	HeaderSize      int
	DataSize        uint32
	DataType        string
	ProtocolVersion int
	ProfileVersion  uint16
	HeaderBytes     []byte
	CRC             []byte
}

func DecodeHeader(b []byte) (Header, error) {
	var h Header
	size, proto, prof, dsize, dtype, hbytes, crc, err := headerParts(b)
	if err != nil {
		return Header{}, err
	}
	h.HeaderSize = int(size[0])
	h.DataSize = binary.LittleEndian.Uint32(dsize)
	h.ProtocolVersion = int(proto[0])
	// ProfileVersion in the SDK does some math as such:
	// if >2199 /1000 else /100
	// Which would turn 2205 to 2.205 or 220 to 2.20
	h.ProfileVersion = binary.LittleEndian.Uint16(prof)
	h.DataType = string(dtype)
	if h.DataType != FIT_FILE_TYPE {
		return Header{}, INVALID_HEADER_TYPE
	}
	h.HeaderBytes = hbytes
	h.CRC = crc
	return h, nil
}

func (s *Header) TotalFileSize() int {
	return s.HeaderSize + int(s.DataSize)
}

func headerParts(h []byte) (size, proto, prof, dsize, dtype, hbytes, crc []byte, err error) {
	if len(h) != 12 && len(h) != 14 {
		err = INVALID_HEADER_LEN
		return
	}
	size = h[0:1]
	proto = h[1:2]
	prof = h[2:4]
	dsize = h[4:8]
	dtype = h[8:12]
	hbytes = h[0:12]
	if len(h) == 14 {
		crc = h[12:14]
	}
	return
}
