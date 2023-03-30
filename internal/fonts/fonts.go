package fonts

import (
	_ "embed"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
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

func loadOpentypeFontFromBytes(bytes []byte, size float64) (font.Face, error) {
	f, err := opentype.Parse(bytes)
	if err != nil {
		return nil, fmt.Errorf("could not parse contents of bytearray: %w", err)
	}

	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, fmt.Errorf("could not create font face: %w", err)
	}
	return face, nil
}

func LoadOpentypeFontFromPath(path string, size float64) (font.Face, error) {
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
	if !strings.HasSuffix(contentType, "ttf") &&
		!strings.HasSuffix(contentType, "otf") {
		return nil, errors.New(
			fmt.Sprintf("file %s is not of type ttf or otf", path),
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
	face, err = loadOpentypeFontFromBytes(bytes, size)
	if err != nil {
		return face, fmt.Errorf(
			"could not parse contents of file %s: %w",
			path,
			err,
		)
	}
	return face, nil
}

func LoadOpentypeFontFromPattern(
	pattern string,
	size float64,
) (font.Face, error) {
	fontPath, err := getFontPathWithFCMatch(pattern)
	if err != nil {
		return nil, err
	}
	return LoadOpentypeFontFromPath(fontPath, size)
}

func LoadOpentypeFontFromFamily(
	family string,
	style string,
	size float64,
) (font.Face, error) {
	pattern := fmt.Sprintf("%s:style=%s", family, style)
	return LoadOpentypeFontFromPattern(pattern, size)
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

func LoadOpentypeFontSetFromFamily(
	family string,
	size float64,
) (*FontSet, error) {
	regular, err := LoadOpentypeFontFromFamily(family, "Regular", size)
	if err != nil {
		return nil, err
	}
	bold, err := LoadOpentypeFontFromFamily(family, "Bold", size)
	if err != nil {
		return nil, err
	}
	return newFontSet(regular, bold), nil
}

func LoadOpentypeFontSetFromPaths(
	regularFontPath, boldFontPath string,
	size float64,
) (*FontSet, error) {
	regular, err := LoadOpentypeFontFromPath(regularFontPath, size)
	if err != nil {
		return nil, err
	}
	bold, err := LoadOpentypeFontFromPath(boldFontPath, size)
	if err != nil {
		return nil, err
	}
	return newFontSet(regular, bold), nil
}

func LoadOpentypeFontSetDefault(size float64) (*FontSet, error) {
	regular, err := loadOpentypeFontFromBytes(inconsolataRegular, size)
	if err != nil {
		return nil, err
	}
	bold, err := loadOpentypeFontFromBytes(inconsolataBold, size)
	if err != nil {
		return nil, err
	}
	return newFontSet(regular, bold), nil
}
