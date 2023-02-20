package minio

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/hlmn/senyum-go-utils/helper"
	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Config struct {
	UrlMinio        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          string
	ProxyUrl        string
	Token           string
}

type Minio struct {
	Config *Config
	client *minio.Client
}

func New(config *Config) (*Minio, error) {
	var err error
	var Minio Minio

	Minio.Config = config
	Minio.client, err = Minio.initMinio()

	return &Minio, err
}

func (m Minio) initMinio() (*minio.Client, error) {
	minioOpts := &minio.Options{
		Creds:  credentials.NewStaticV4(m.Config.AccessKeyID, m.Config.SecretAccessKey, m.Config.Token),
		Secure: false,
	}

	if m.Config.UseSSL == "true" {
		minioOpts.Secure = true
	}

	if m.Config.ProxyUrl != "" {
		proxyUrl, _ := url.Parse(m.Config.ProxyUrl)
		minioOpts.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	}

	minioClient, err := minio.New(m.Config.UrlMinio, minioOpts)

	// _, err := minioClient.BucketExists(context.Background(), m.Config.Bucket)
	// if err != nil {
	// 	return minioClient, err

	// }

	return minioClient, err

}

// func (config Config) FileBase64(fileName string, fileContent string) (*string, error) {
// 	idx := strings.Index(fileContent, ";base64,")

// 	if idx < 0 {
// 		return nil, errors.New("format not match")
// 	}

// 	imageType := fileContent[11:idx]

// 	unbased, err := base64.StdEncoding.DecodeString(fileContent[idx+8:])
// 	if err != nil {
// 		return nil, err
// 	}

// 	minioClient, errs := config.initMinio()
// 	if errs != nil {
// 		return nil, errs
// 	}

// 	file := bytes.NewReader(unbased)

// 	// fileName := fmt.Sprint(uuid.New(), ".", imageType)
// 	cacheControl := "max-age=31536000"
// 	userMetaData := map[string]string{"x-amz-acl": "public-read"}

// 	info, err := minioClient.PutObject(context.Background(), config.Bucket, fileName, file, -1, minio.PutObjectOptions{
// 		ContentType:  "image/" + imageType,
// 		CacheControl: cacheControl,
// 		UserMetadata: userMetaData,
// 	})

// 	return &info.Key, err
// }

func (m Minio) Upload(path, bucket string, file *helper.File) (key string, err error) {
	minioClient := m.client
	ctx := context.Background()

	err = minioClient.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(ctx, bucket)
		if !(errBucketExists == nil && exists) {
			return "", err
		}
	}

	if file.Reader == nil {
		return key, errors.New("Reader not found")
	}

	filePath := fmt.Sprintf("%s/%s", path, file.Name)

	info, err := minioClient.PutObject(ctx, bucket, filePath, file.Reader, file.Reader.Size(), minio.PutObjectOptions{})
	fmt.Println(err)
	return info.Key, err
}

func (m Minio) UploadMultiPartFile(path, bucket string, fileContent *multipart.FileHeader) (string, error) {

	// Get Buffer from file
	buffer, err := fileContent.Open()

	if err != nil {
		return "", err
	}

	defer buffer.Close()

	minioClient := m.client

	ctx := context.Background()
	objectName := fileContent.Filename
	fileBuffer := buffer
	fileSize := fileContent.Size

	err = minioClient.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(ctx, bucket)
		if !(errBucketExists == nil && exists) {
			return "", err
		}
	}
	// cacheControl := "max-age=31536000"
	// userMetaData := map[string]string{"x-amz-acl": "public-read"}

	filePath := fmt.Sprintf("%s/%s", path, objectName)
	info, err := minioClient.PutObject(ctx, bucket, filePath, fileBuffer, fileSize, minio.PutObjectOptions{})

	return info.Key, err
}

func (m Minio) GetFile(bucket, path string) (file *helper.File, err error) {
	minioClient := m.client

	reader, err := minioClient.GetObject(context.Background(), bucket, path, minio.GetObjectOptions{})
	if err != nil {
		return file, err
	}

	w := &bytes.Buffer{}
	if _, err := io.Copy(w, reader); err != nil {
		return file, err
	}

	base64 := base64.StdEncoding.EncodeToString(w.Bytes())
	mime := mimetype.Detect(w.Bytes())

	// if err != nil {
	// 	return file, err
	// }

	reader.Seek(0, 0)

	filename := strings.Split(path, "/")

	readerBytes := bytes.NewReader(w.Bytes())

	file = &helper.File{
		Base64:    base64,
		Extension: mime.Extension(),
		Name:      filename[len(filename)-1],
		Mime:      mime.String(),
		Reader:    readerBytes,
	}

	return file, err
}
