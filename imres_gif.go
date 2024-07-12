package imres

import (
	"encoding/binary"
	"errors"
)

// GetGifDimensions extracts dimensions from a GIF header.
func GetGifDimensions(header []byte) (width, height int, err error) {
	if len(header) < 10 {
		return 0, 0, errors.New("invalid GIF header")
	}

	width = int(binary.LittleEndian.Uint16(header[6:8]))
	height = int(binary.LittleEndian.Uint16(header[8:10]))

	return width, height, nil
}
