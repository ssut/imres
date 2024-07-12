package imres

import (
	"encoding/binary"
	"errors"
	"io"
)

// GetTiffDimensions extracts dimensions from a TIFF file.
func GetTiffDimensions(r io.ReadSeeker) (width, height int, err error) {
	// Ensure we start from the beginning
	_, err = r.Seek(0, io.SeekStart)
	if err != nil {
		return 0, 0, err
	}

	header := make([]byte, 8)
	_, err = r.Read(header)
	if err != nil {
		return 0, 0, err
	}

	var order binary.ByteOrder
	switch string(header[:2]) {
	case "II":
		order = binary.LittleEndian

	case "MM":
		order = binary.BigEndian

	default:
		return 0, 0, errors.New("not a valid TIFF file")
	}

	// read the offset to the first IFD
	offset := order.Uint32(header[4:])
	if _, err = r.Seek(int64(offset), io.SeekStart); err != nil {
		return 0, 0, err
	}

	// read the number of IFD entries
	var numEntries uint16
	if err = binary.Read(r, order, &numEntries); err != nil {
		return 0, 0, err
	}

	for i := 0; i < int(numEntries); i++ {
		entry := make([]byte, 12)
		if _, err = r.Read(entry); err != nil {
			return 0, 0, err
		}

		tag := order.Uint16(entry[0:2])
		fieldType := order.Uint16(entry[2:4])
		count := order.Uint32(entry[4:8])
		valueOffset := entry[8:12]

		// width
		if tag == tiffTagImageWidth && fieldType == 3 && count == 1 {
			width = int(order.Uint16(valueOffset))
		} else if tag == tiffTagImageWidth && fieldType == 4 && count == 1 {
			width = int(order.Uint32(valueOffset))
		}

		// height
		if tag == tiffTagImageHeight && fieldType == 3 && count == 1 {
			height = int(order.Uint16(valueOffset))
		} else if tag == tiffTagImageHeight && fieldType == 4 && count == 1 {
			height = int(order.Uint32(valueOffset))
		}

		if width != 0 && height != 0 {
			break
		}
	}

	if width == 0 || height == 0 {
		return 0, 0, errors.New("could not find width or height in TIFF file")
	}

	return width, height, nil
}
