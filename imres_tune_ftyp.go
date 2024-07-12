package imres

import (
	"errors"
	"io"
)

// findFtyp searches for the "ftyp" box in the byte stream and returns its position.
func findFtyp(r io.ReadSeeker) (int64, error) {
	buffer := make([]byte, 8)
	for {
		_, err := io.ReadFull(r, buffer)
		if err != nil {
			if err == io.EOF {
				return 0, errors.New("ftyp not found")
			}
			return 0, err
		}
		if string(buffer[4:8]) == "ftyp" {
			currentPos, _ := r.Seek(0, io.SeekCurrent)
			return currentPos - 8, nil
		}
		_, err = r.Seek(-7, io.SeekCurrent)
		if err != nil {
			return 0, err
		}
	}
}
