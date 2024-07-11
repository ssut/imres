package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
	"testing"
)

func parseFilename(filename string) (width, height int, err error) {
	base := filepath.Base(filename)
	parts := strings.Split(base, ".")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid file name format: %s", filename)
	}

	widthHeightParts := strings.Split(filename, "x")
	widthStr := widthHeightParts[0]
	heightStr := strings.Split(widthHeightParts[1], "_")[0]

	width, err = strconv.Atoi(widthStr)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid width: %s", widthStr)
	}

	height, err = strconv.Atoi(heightStr)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid height: %s", heightStr)
	}

	return width, height, nil
}

func getTestFiles() ([]fs.FileInfo, error) {
	allowedExtensions := [...]string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".avif"}

	f, err := os.Open("testdata")
	if err != nil {
		return nil, err
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	var filteredList []fs.FileInfo
	for _, file := range list {
		// check allowed exts
		if !file.IsDir() && slices.Contains(allowedExtensions[:], filepath.Ext(file.Name())) {
			filteredList = append(filteredList, file)
		}
	}

	sort.Slice(filteredList, func(i, j int) bool {
		return filteredList[i].Name() < filteredList[j].Name()
	})
	return filteredList, nil
}

func TestImageDimensions(t *testing.T) {
	files, err := getTestFiles()
	if err != nil {
		t.Fatalf("failed to read testdata directory: %v", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filepath := filepath.Join("testdata", file.Name())
		expectedWidth, expectedHeight, err := parseFilename(file.Name())
		if err != nil {
			t.Logf("failed to parse filename %s: %v", file.Name(), err)
			continue
		}

		f, err := os.Open(filepath)
		if err != nil {
			t.Errorf("failed to open file %s: %v", file.Name(), err)
			continue
		}

		width, height, err := GetImageDimensions(f)
		f.Close()
		if err != nil {
			t.Errorf("failed to get dimensions for file %s: %v", file.Name(), err)
			continue
		}

		if width != expectedWidth || height != expectedHeight {
			t.Errorf("dimensions for file %s do not match: expected %dx%d, got %dx%d", file.Name(), expectedWidth, expectedHeight, width, height)
			continue
		}

		t.Logf("dimensions for file %s match: %dx%d", file.Name(), width, height)
	}
}
