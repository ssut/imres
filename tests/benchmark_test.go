package tests

import (
	"bytes"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"path/filepath"
	"testing"

	"os"

	"github.com/ssut/imres"
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

		expectedWidth, expectedHeight, _ := parseFilename(file.Name())
		b.Run(file.Name(), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				r := bytes.NewReader(data)
				width, height, err := imres.GetImageDimensions(r)
				if err != nil {
					b.Errorf("failed to get dimensions for file %s: %v", file.Name(), err)
				}

				if width != expectedWidth || height != expectedHeight {
					b.Errorf("expected dimensions for file %s: %dx%d, got: %dx%d", file.Name(), expectedWidth, expectedHeight, width, height)
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

		expectedWidth, expectedHeight, _ := parseFilename(file.Name())
		b.Run(file.Name(), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				r := bytes.NewReader(data)
				config, _, err := image.DecodeConfig(r)
				if err != nil {
					b.Errorf("failed to decode config for file %s: %v", file.Name(), err)
				}

				if config.Width != expectedWidth || config.Height != expectedHeight {
					b.Errorf("expected dimensions for file %s: %dx%d, got: %dx%d", file.Name(), expectedWidth, expectedHeight, config.Width, config.Height)
				}
			}
		})
	}
}
