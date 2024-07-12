package imres

import (
	"encoding/binary"
	"errors"
	"io"
)

// GetWebPDimensions extracts dimensions from a WebP header.
func GetWebPDimensions(r io.Reader, header []byte) (width, height int, err error) {
	if string(header[0:4]) != "RIFF" || string(header[8:12]) != "WEBP" {
		return 0, 0, errors.New("not a valid WebP file")
	}

	switch string(header[12:16]) {
	case "VP8 ":
		width = int(binary.LittleEndian.Uint16(header[26:28])) & 0x3FFF
		height = int(binary.LittleEndian.Uint16(header[28:30])) & 0x3FFF

	case "VP8L":
		bits := uint32(header[21]) | uint32(header[22])<<8 | uint32(header[23])<<16 | uint32(header[24])<<24
		width = int(bits&0x3FFF) + 1
		height = int((bits>>14)&0x3FFF) + 1

	case "VP8X":
		if len(header) < 30 {
			extraBytes := make([]byte, 30-len(header))
			_, err := io.ReadFull(r, extraBytes)
			if err != nil {
				return 0, 0, err
			}
			header = append(header, extraBytes...)
		}

		width = int(uint32(header[24]) | uint32(header[25])<<8 | uint32(header[26])<<16)
		height = int(uint32(header[27]) | uint32(header[28])<<8 | uint32(header[29])<<16)

		width++
		height++

	default:
		return 0, 0, errors.New("unsupported WebP format")
	}

	return width, height, nil
}
