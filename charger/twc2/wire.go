package twc2

import (
	"bytes"
	"errors"
)

const (
	delimiter = 0xC0
	escape    = 0xDB
)

var (
	delimiterSequence = []byte{0xDB, 0xDC}
	escapeSequence    = []byte{0xDB, 0xDD}
)

// Checksum calculates message checksum starting from 2nd byte
func Checksum(msg []byte) byte {
	var acc byte
	for _, b := range msg[1:] {
		acc += b
	}
	return acc
}

// Encode adds checksum, escapes special characters and adds delimiters
func Encode(msg []byte) []byte {
	msg = append(msg, Checksum(msg))

	// substitute escape and delimiter characters
	msg = bytes.ReplaceAll(msg, []byte{escape}, escapeSequence)
	msg = bytes.ReplaceAll(msg, []byte{delimiter}, delimiterSequence)

	// add delimiters
	buf := bytes.NewBuffer([]byte{delimiter})
	buf.Write(msg)
	buf.Write([]byte{delimiter})

	return buf.Bytes()
}

// Decode reverses Encode
func Decode(msg []byte) ([]byte, error) {
	// must be at least 2 delimiters, payload and checksum
	if len(msg) < 4 {
		return []byte{}, errors.New("invalid message")
	}

	// strip delimiters
	msg = msg[1 : len(msg)-1]

	// substitute escape and delimiter sequences
	msg = bytes.ReplaceAll(msg, delimiterSequence, []byte{delimiter})
	msg = bytes.ReplaceAll(msg, escapeSequence, []byte{escape})

	// validate checksum
	cks := Checksum(msg[0 : len(msg)-1])
	if cks != msg[len(msg)-1] {
		return []byte{}, errors.New("invalid checksum")
	}

	return msg[0 : len(msg)-1], nil
}
