package imres

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

// GetAvifDimensions extracts dimensions from an AVIF header.
func GetAvifDimensions(r io.ReadSeeker, header []byte) (width, height int, err error) {
	// Ensure we start from the beginning
	_, err = r.Seek(0, io.SeekStart)
	if err != nil {
		return 0, 0, err
	}

	// Find "ftyp" in the byte stream
	pos, err := findFtyp(r)
	if err != nil {
		return 0, 0, err
	}

	// Seek to the position of "ftyp"
	_, err = r.Seek(pos, io.SeekStart)
	if err != nil {
		return 0, 0, err
	}

	buffer := make([]byte, 8)
	ispeSignature := []byte{0x69, 0x73, 0x70, 0x65} // "ispe" in bytes

	for {
		_, err := io.ReadFull(r, buffer)
		if err != nil {
			if err == io.EOF {
				return 0, 0, errors.New("ispe box not found")
			}
			return 0, 0, err
		}

		if bytes.Equal(buffer[4:], ispeSignature) {
			ispeData := make([]byte, 12)
			_, err := io.ReadFull(r, ispeData)
			if err != nil {
				return 0, 0, err
			}

			width = int(binary.BigEndian.Uint32(ispeData[4:8]))
			height = int(binary.BigEndian.Uint32(ispeData[8:12]))
			return width, height, nil
		} else {
			_, err = r.Seek(-7, io.SeekCurrent)
			if err != nil {
				return 0, 0, err
			}
		}
	}
}
