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

	switch header[0] {
	case 0x89:
		if header[1] == 0x50 && header[2] == 0x4E && header[3] == 0x47 {
			return GetPngDimensions(header)
		}

	case 0xFF:
		if header[1] == 0xD8 && header[2] == 0xFF {
			return GetJpegDimensions(r, header)
		}

	case 'G':
		if header[1] == 'I' && header[2] == 'F' {
			if header[3] == '8' && (header[4] == '7' || header[4] == '9') && header[5] == 'a' {
				return GetGifDimensions(header)
			}
		}

	case 'R':
		if header[1] == 'I' && header[2] == 'F' && header[3] == 'F' {
			if header[8] == 'W' && header[9] == 'E' && header[10] == 'B' && header[11] == 'P' {
				return GetWebPDimensions(r, header)
			}
		}

	case 'B':
		if header[1] == 'M' {
			return GetBmpDimensions(header)
		}

	case 'I', 'M':
		if (header[0] == 'I' && header[1] == 'I') || (header[0] == 'M' && header[1] == 'M') {
			return GetTiffDimensions(r)
		}
	}

	// Check for file type box
	if string(header[4:8]) == "ftyp" {
		nextBytes := string(header[8:12])
		switch nextBytes {
		case "mif1", "msf1", "heic", "heix", "hevc", "hevx":
			return GetHeifDimensions(r)

		case "avif", "avis":
			return GetAvifDimensions(r, header)
		}
	}

	return 0, 0, errors.New("unsupported file format")
}
