package helper

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
)

type Reader interface {
	Size() int64
	io.Reader
	io.ReaderAt
	io.Seeker
}

type File struct {
	Reader Reader
	// Reader    *bytes.Reader
	Name      string
	Extension string
	Mime      string
}

func NewFileFromBase64(fileName, base64Content string) (file *File, err error) {
	file = &File{}
	if fileName == "" {
		name, _ := uuid.NewRandom()
		fileName = name.String()
	}

	sDec, err := base64.StdEncoding.DecodeString(base64Content)
	if err != nil {
		return nil, err
	}
	fileBuffer := bytes.NewReader(sDec)

	mime, err := mimetype.DetectReader(fileBuffer)

	if err != nil {
		return nil, err
	}

	fileBuffer.Seek(0, 0)

	file.Reader = fileBuffer
	file.Extension = mime.Extension()
	file.Name = fmt.Sprintf("%s%s", fileName, mime.Extension())
	file.Mime = mime.String()

	return file, err
}
