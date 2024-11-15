// Protocol package implements a simple protocol for
// sending and receiving data.
//
// Protocol format: ABBBBDDD...DDD
// Where:
// A: command byte
// C: payload length (32 bits, big endian)
// D: payload (array of bytes)

package angus

import (
	"encoding/binary"
	"errors"
)

const (
	MaxPackageSize = BufferSize + 5
)

var (
	ErrInvalidSize = errors.New("invalid size")
)

// Encode encodes the source data into the destination buffer
// using the specified command.
// It returns the number of bytes written and an error, if any.
func Encode(dest, src []byte, cmd byte) (int, error) {
	lenData := len(src)
	if lenData > MaxPackageSize {
		return 0, ErrInvalidSize
	}
	if len(dest) < lenData+5 {
		return 0, ErrInvalidSize
	}
	dest[0] = cmd
	binary.BigEndian.PutUint32(dest[1:], uint32(lenData))
	copy(dest[5:], src)
	n := lenData + 5
	return n, nil
}

// Decode decodes the source buffer into the destination buffer.
// It returns the command byte, the number of bytes read, the
// counter value, and an error, if any.
// command byte + data length = 5 bytes
func Decode(dest, src []byte) (cmd byte, n int, err error) {
	if len(src) < 5 {
		return 0, 0, ErrInvalidSize
	}
	lenData := int(binary.BigEndian.Uint32(src[1:]))
	if lenData > BufferSize {
		return 0, 0, ErrInvalidSize
	}
	if len(src) < lenData+5 {
		return 0, 0, ErrInvalidSize
	}
	copy(dest, src[5:5+lenData])
	return src[0], lenData, nil
}
