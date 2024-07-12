package imres

import (
	"bytes"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"path/filepath"
	"testing"

	"os"

	_ "golang.org/x/image/webp"
)

func BenchmarkCustomImageDimensions(b *testing.B) {
	files, err := getTestFiles()
	if err != nil {
		b.Fatalf("failed to read testdata directory: %v", err)
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filepath := filepath.Join("testdata", file.Name())
		data, err := os.ReadFile(filepath)
		if err != nil {
			b.Errorf("failed to read file %s: %v", file.Name(), err)
			continue
		}

		b.Run(file.Name(), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				r := bytes.NewReader(data)
				_, _, err := GetImageDimensions(r)
				if err != nil {
					b.Errorf("failed to get dimensions for file %s: %v", file.Name(), err)
				}
			}
		})
	}
}

func BenchmarkStandardImageDimensions(b *testing.B) {
	files, err := getTestFiles()
	if err != nil {
		b.Fatalf("failed to read testdata directory: %v", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filepath := filepath.Join("testdata", file.Name())
		data, err := os.ReadFile(filepath)
		if err != nil {
			b.Errorf("failed to read file %s: %v", file.Name(), err)
			continue
		}

		b.Run(file.Name(), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				r := bytes.NewReader(data)
				_, _, err := image.DecodeConfig(r)
				if err != nil {
					b.Errorf("failed to decode config for file %s: %v", file.Name(), err)
				}
			}
		})
	}
}
