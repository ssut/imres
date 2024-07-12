package imres

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

// GetJpegDimensions extracts dimensions from a JPEG header.
func GetJpegDimensions(r io.Reader, header []byte) (width, height int, err error) {
	// FFD8
	buf := header[2:]

	for {
		if len(buf) < 2 {
			tempBuf := make([]byte, 2-len(buf))
			n, err := r.Read(tempBuf)
			if err != nil {
				return 0, 0, err
			}
			buf = append(buf, tempBuf[:n]...)
			if len(buf) < 2 {
				return 0, 0, errors.New("invalid JPEG data (short marker)")
			}
		}

		if buf[0] != 0xFF {
			return 0, 0, fmt.Errorf("invalid JPEG data (0xFF marker not found, got 0x%02X)", buf[0])
		}

		marker := buf[1]
		buf = buf[2:]

		if marker >= 0xC0 && marker <= 0xCF && marker != 0xC4 && marker != 0xC8 && marker != 0xCC {
			if len(buf) < 7 {
				tempBuf := make([]byte, 7-len(buf))
				n, err := r.Read(tempBuf)
				if err != nil {
					return 0, 0, err
				}

				buf = append(buf, tempBuf[:n]...)
				if len(buf) < 7 {
					return 0, 0, errors.New("invalid JPEG data (short segment)")
				}
			}

			height = int(binary.BigEndian.Uint16(buf[3:5]))
			width = int(binary.BigEndian.Uint16(buf[5:7]))

			return width, height, nil
		}

		if len(buf) < 2 {
			tempBuf := make([]byte, 2-len(buf))
			n, err := r.Read(tempBuf)
			if err != nil {
				return 0, 0, err
			}

			buf = append(buf, tempBuf[:n]...)
			if len(buf) < 2 {
				return 0, 0, errors.New("invalid JPEG data (short segment length)")
			}
		}

		segmentLength := int(binary.BigEndian.Uint16(buf[0:2]))
		buf = buf[2:]

		if segmentLength < 2 {
			return 0, 0, errors.New("invalid JPEG data (segment length is too short)")
		}

		skipLength := segmentLength - 2
		if len(buf) < skipLength {
			tempBuf := make([]byte, skipLength-len(buf))
			n, err := r.Read(tempBuf)
			if err != nil {
				return 0, 0, err
			}

			buf = append(buf, tempBuf[:n]...)
			if len(buf) < skipLength {
				return 0, 0, errors.New("invalid JPEG data (short segment skip)")
			}
		}

		buf = buf[skipLength:]
	}
}
