package fonts

import (
	_ "embed"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

//go:embed static/Inconsolata-Regular.ttf
var inconsolataRegular []byte

//go:embed static/Inconsolata-Bold.ttf
var inconsolataBold []byte

func getFontPathWithFCMatch(pattern string) (string, error) {
	fontPath, err := exec.Command("fc-match", "--format=%{file}", pattern).
		Output()
	if err != nil {
		return string(
				fontPath,
			), fmt.Errorf(
				"could not find font path for pattern %s: %w",
				pattern,
				err,
			)
	}
	return string(fontPath), nil
}

func loadTTFontFromBytes(bytes []byte, size float64) (font.Face, error) {
	font, err := truetype.Parse(bytes)
	if err != nil {
		return nil, fmt.Errorf("could not parse contents of bytearray: %w", err)
	}

	face := truetype.NewFace(font, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	})
	return face, nil
}

func LoadTTFontFromPath(path string, size float64) (font.Face, error) {
	var (
		file *os.File
		face font.Face
		err  error
	)
	file, err = os.Open(path)
	if err != nil {
		return nil, fmt.Errorf(
			"could not open font file at path %s: %w",
			path,
			err,
		)
	}
	defer file.Close()

	contentType, err := GetFileContentType(file)
	if err != nil {
		return nil, errors.New(
			fmt.Sprintf(
				"could not read content type of file %s: %w",
				path,
				err,
			),
		)
	}
	if !strings.HasSuffix(contentType, "ttf") {
		return nil, errors.New(
			fmt.Sprintf("file %s is not of type ttf", path),
		)
	}

	var bytes []byte
	bytes, err = io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf(
			"could not read contents of file %s: %w",
			path,
			err,
		)
	}
	face, err = loadTTFontFromBytes(bytes, size)
	if err != nil {
		return face, fmt.Errorf(
			"could not parse contents of file %s: %w",
			path,
			err,
		)
	}
	return face, nil
}

func LoadTTFontFromPattern(pattern string, size float64) (font.Face, error) {
	fontPath, err := getFontPathWithFCMatch(pattern)
	if err != nil {
		return nil, err
	}
	return LoadTTFontFromPath(fontPath, size)
}

func LoadTTFontFromFamily(
	family string,
	style string,
	size float64,
) (font.Face, error) {
	pattern := fmt.Sprintf("%s:style=%s", family, style)
	return LoadTTFontFromPattern(pattern, size)
}

type FontSet struct {
	Regular font.Face
	Bold    font.Face
}

func newFontSet(regular font.Face, bold font.Face) *FontSet {
	fs := FontSet{
		Regular: regular,
		Bold:    bold,
	}
	return &fs
}

func LoadTTFontSetFromFamily(family string, size float64) (*FontSet, error) {
	regular, err := LoadTTFontFromFamily(family, "Regular", size)
	if err != nil {
		return nil, err
	}
	bold, err := LoadTTFontFromFamily(family, "Bold", size)
	if err != nil {
		return nil, err
	}
	return newFontSet(regular, bold), nil
}

func LoadTTFontSetFromPaths(
	regularFontPath, boldFontPath string,
	size float64,
) (*FontSet, error) {
	regular, err := LoadTTFontFromPath(regularFontPath, size)
	if err != nil {
		return nil, err
	}
	bold, err := LoadTTFontFromPath(boldFontPath, size)
	if err != nil {
		return nil, err
	}
	return newFontSet(regular, bold), nil
}

func LoadTTFontSetDefault(size float64) (*FontSet, error) {
	regular, err := loadTTFontFromBytes(inconsolataRegular, size)
	if err != nil {
		return nil, err
	}
	bold, err := loadTTFontFromBytes(inconsolataBold, size)
	if err != nil {
		return nil, err
	}
	return newFontSet(regular, bold), nil
}
