package fonts

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

func GetFontPathWithFCMatch(pattern string) (string, error) {
	fontPath, err := exec.Command("fc-match", "--format=%{file}", pattern).Output()
	if err != nil {
		return string(fontPath), fmt.Errorf("could not find font path for pattern %s: %w", pattern, err)
	}
	return string(fontPath), nil
}

func LoadTTF(path string, size float64) (font.Face, error) {
	var (
		file *os.File
		face font.Face
		err  error
	)
	file, err = os.Open(path)
	if err != nil {
		return face, fmt.Errorf("could not open font file at path %s: %w", path, err)
	}
	defer file.Close()

	var bytes []byte
	bytes, err = io.ReadAll(file)
	if err != nil {
		return face, fmt.Errorf("could not read contents of file %s: %w", path, err)
	}

	var font *truetype.Font
	font, err = truetype.Parse(bytes)
	if err != nil {
		return face, fmt.Errorf("could not parse contents of file %s: %w", path, err)
	}

	face = truetype.NewFace(font, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	})
	return face, nil
}

func LoadTTFontFromPattern(pattern string, size float64) (font.Face, error) {
	fontPath, err := GetFontPathWithFCMatch(pattern)
	if err != nil {
		return nil, err
	}
	return LoadTTF(fontPath, size)
}
