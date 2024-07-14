package imres

import (
	"encoding/binary"
	"errors"
)

// GetBmpDimensions extracts dimensions from a BMP file.
func GetBmpDimensions(header []byte) (width, height int, err error) {
	if string(header[0:2]) != "BM" {
		return 0, 0, errors.New("not a valid BMP file")
	}

	// DIB header
	// https://docs.fileformat.com/image/dib/
	width = int(binary.LittleEndian.Uint32(header[18:22]))
	height = int(binary.LittleEndian.Uint32(header[22:26]))

	if width <= 0 || height <= 0 {
		return 0, 0, errors.New("invalid BMP dimensions")
	}

	return width, height, nil
}
