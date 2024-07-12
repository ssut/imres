package imres

import (
	"errors"
	"io"
)

// GetImageDimensions reads the dimensions from a readable stream.
func GetImageDimensions(r io.ReadSeeker) (width, height int, err error) {
	header := make([]byte, 30)
	_, err = r.Read(header)
	if err != nil {
		return 0, 0, err
	}

	switch {
	case header[0] == 0x89 && header[1] == 0x50 && header[2] == 0x4E && header[3] == 0x47:
		return GetPngDimensions(header)

	case header[0] == 0xFF && header[1] == 0xD8 && header[2] == 0xFF:
		return GetJpegDimensions(r, header)

	case header[0] == 'G' && header[1] == 'I' && header[2] == 'F' && (header[3] == '8' && (header[4] == '7' || header[4] == '9') && header[5] == 'a'):
		return GetGifDimensions(header)

	case header[0] == 'R' && header[1] == 'I' && header[2] == 'F' && header[3] == 'F' && header[8] == 'W' && header[9] == 'E' && header[10] == 'B' && header[11] == 'P':
		return GetWebPDimensions(r, header)

	case string(header[4:8]) == "ftyp":
		nextBytes := string(header[8:12])
		switch nextBytes {
		case "mif1", "msf1", "heic", "heix", "hevc", "hevx":
			return GetHeifDimensions(r)

		case "avif", "avis":
			return GetAvifDimensions(r, header)
		}

		fallthrough

	default:
		return 0, 0, errors.New("unsupported file format")
	}
}
