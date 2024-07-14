package imres

import (
	"encoding/binary"
	"io"
)

// GetHeifDimensions extracts dimensions.
// This can be used for both Heif and Avif: https://aomediacodec.github.io/av1-avif/v1.1.0.html#image-spatial-extents-property
func GetHeifDimensions(r io.ReadSeeker, header []byte) (width, height int, err error) {
	// Ensure we start from the beginning
	_, err = r.Seek(0, io.SeekStart)
	if err != nil {
		return 0, 0, err
	}

	buffer := make([]byte, 8)

	for {
		_, err := io.ReadFull(r, buffer)
		if err != nil {
			return 0, 0, err
		}

		if string(buffer) == ispeSignature {
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
