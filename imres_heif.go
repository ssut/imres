package imres

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

// GetHeifDimensions extracts dimensions from a HEIF header.
func GetHeifDimensions(r io.ReadSeeker) (width, height int, err error) {
	_, err = r.Seek(0, io.SeekStart)
	if err != nil {
		return 0, 0, err
	}

	pos, err := findFtyp(r)
	if err != nil {
		return 0, 0, err
	}

	_, err = r.Seek(pos, io.SeekStart)
	if err != nil {
		return 0, 0, err
	}

	_, err = r.Seek(8, io.SeekCurrent)
	if err != nil {
		return 0, 0, err
	}

	buffer := make([]byte, 8)

	for {
		_, err := io.ReadFull(r, buffer)
		if err != nil {
			if err == io.EOF {
				return 0, 0, errors.New("ispe box not found")
			}
			return 0, 0, err
		}

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
			_, err = r.Seek(-7, io.SeekCurrent)
			if err != nil {
				return 0, 0, err
			}
		}
	}
}
