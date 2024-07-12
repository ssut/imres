package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

var ispeSignature = []byte{0x00, 0x00, 0x00, 0x14, 0x69, 0x73, 0x70, 0x65}

// GetImageDimensions reads the dimensions from a readable stream.
func GetImageDimensions(r io.ReadSeeker) (width, height int, err error) {
	// 30-byte header should be enough to determine *most* image types
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

	case string(header[4:8]) == "ftyp" && (string(header[8:12]) == "avif" || string(header[8:12]) == "avis"):
		return GetAvifDimensions(r, header)

	default:
		return 0, 0, errors.New("unsupported file format")
	}
}

// GetWebPDimensions extracts dimensions from a WebP header.
func GetWebPDimensions(r io.Reader, header []byte) (width, height int, err error) {
	if string(header[0:4]) != "RIFF" || string(header[8:12]) != "WEBP" {
		return 0, 0, errors.New("not a valid WebP file")
	}

	switch string(header[12:16]) {
	case "VP8":
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

// GetPngDimensions extracts dimensions from a PNG header.
func GetPngDimensions(header []byte) (width, height int, err error) {
	if len(header) < 24 {
		return 0, 0, errors.New("invalid PNG header")
	}

	width = int(binary.BigEndian.Uint32(header[16:20]))
	height = int(binary.BigEndian.Uint32(header[20:24]))

	return width, height, nil
}

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

// GetGifDimensions extracts dimensions from a GIF header.
func GetGifDimensions(header []byte) (width, height int, err error) {
	if len(header) < 10 {
		return 0, 0, errors.New("invalid GIF header")
	}

	width = int(binary.LittleEndian.Uint16(header[6:8]))
	height = int(binary.LittleEndian.Uint16(header[8:10]))

	return width, height, nil
}

// GetAvifDimensions extracts dimensions from an AVIF header.
func GetAvifDimensions(r io.ReadSeeker, header []byte) (width, height int, err error) {
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
