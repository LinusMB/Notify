package font

import (
	"io"
	"net/http"
	"os"
)

func GetFileContentType(file *os.File) (string, error) {
	buf := make([]byte, 512)
	_, err := file.Read(buf)

	if err != nil {
		return "", err
	}

	contentType := http.DetectContentType(buf)

	file.Seek(0, io.SeekStart)

	return contentType, nil
}
