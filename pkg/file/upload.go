package file

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func GenerateUniqueFilename(filename string) string {
	ext := filepath.Ext(filename)
	name := strings.TrimSuffix(filename, ext)
	return fmt.Sprintf("%s_%d%s", name, GetTimestamp(), ext)
}

func GetTimestamp() int64 {
	return time.Now().Unix()
}

func SaveFile(file *multipart.FileHeader, uploadDir, filename string) error {
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return err
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(filepath.Join(uploadDir, filename))
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}
