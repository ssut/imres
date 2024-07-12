package imres

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

// GetHeifDimensions extracts dimensions from a HEIF header.
func GetHeifDimensions(r io.ReadSeeker) (width, height int, err error) {
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

	// Skip the "ftyp" box content
	_, err = r.Seek(8, io.SeekCurrent)
	if err != nil {
		return 0, 0, err
	}

	// Define the ispe signature
	ispeSignature := []byte{0x00, 0x00, 0x00, 0x14, 0x69, 0x73, 0x70, 0x65}

	// Buffer to read the stream
	buffer := make([]byte, 8)

	for {
		_, err := io.ReadFull(r, buffer)
		if err != nil {
			if err == io.EOF {
				return 0, 0, errors.New("ispe box not found")
			}
			return 0, 0, err
		}

		// Check if the buffer matches the ispe signature
		if bytes.Equal(buffer, ispeSignature) {
			ispeData := make([]byte, 12)
			_, err := io.ReadFull(r, ispeData)
			if err != nil {
				return 0, 0, err
			}

			width = int(binary.BigEndian.Uint32(ispeData[4:8]))
			height = int(binary.BigEndian.Uint32(ispeData[8:12]))
			return width, height, nil
		} else {
			// Move back 7 bytes to ensure overlapping sequences are handled
			_, err = r.Seek(-7, io.SeekCurrent)
			if err != nil {
				return 0, 0, err
			}
		}
	}
}
