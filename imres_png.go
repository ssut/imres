package imres

import (
	"encoding/binary"
	"errors"
)

// GetPngDimensions extracts dimensions from a PNG header.
func GetPngDimensions(header []byte) (width, height int, err error) {
	if len(header) < 24 {
		return 0, 0, errors.New("invalid PNG header")
	}

	width = int(binary.BigEndian.Uint32(header[16:20]))
	height = int(binary.BigEndian.Uint32(header[20:24]))

	return width, height, nil
}
