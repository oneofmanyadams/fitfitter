package fitprotocol

import "errors"

const (
	MSG_RESERVED_SIZE        = 1
	MSG_ARCHITECTURE_SIZE    = 1
	MSG_GLOBAL_NUM_SIZE      = 2
	MSG_FIELD_COUNT_SIZE     = 1
	MSG_FIELD_DEF_SIZE       = 3
	MSG_FIELD_COUNT_SIZE_DEV = 1
	MSG_FIELD_DEF_SIZE_DEV   = 3
)

func MsgFixedContentSize() int {
	return MSG_RESERVED_SIZE +
		MSG_ARCHITECTURE_SIZE +
		MSG_GLOBAL_NUM_SIZE +
		MSG_FIELD_COUNT_SIZE
}

type DefinitionMessage struct {
	Architecture        uint8 // 0 is little endian.
	GlobalMessageNumber uint16
	NumberOfFields      uint8
	FieldDefinitions    []FieldDefinition
	DevFlag             bool
	NumberOfDevFields   uint8
	DevFieldDefinitions []FieldDefinition
}

func NewDefinitionMessage(b []byte) DefinitionMessage {
	var dm DefinitionMessage
	dm.Architecture = b[1]
	dm.GlobalMessageNumber = uint16(b[2])<<8 + uint16(b[3])
	dm.NumberOfFields = b[4]
	return dm
}

func (s *DefinitionMessage) AddFieldDef(bytes []byte) error {
	var fd FieldDefinition
	if len(bytes) != MSG_FIELD_DEF_SIZE {
		return errors.New("Incorrect amount of bytes provided.")
	}
	fd.Number = bytes[0]
	fd.Size = bytes[1]
	fd.BaseType = bytes[2]
	fd.Bytes = bytes
	s.FieldDefinitions = append(s.FieldDefinitions, fd)
	return nil
}

type FieldDefinition struct {
	Number   uint8 //Defined in the global fit profile for the specified message.
	Size     uint8 // size in bytes of specifieed fit message's field.
	BaseType uint8 // (unsigned char, etc...) defined in fit.h in SDK
	Bytes    []byte
}