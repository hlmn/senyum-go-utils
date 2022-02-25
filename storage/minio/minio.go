package minio

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Config struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	ProxyUrl        string
	Token           string
	Bucket          string
	ProxyPublic     string
	UrlPublic       string
	UrlProxyPublic  string
}

func New(config *Config) (*Config, error) {
	var err error
	minio := Config{}
	minio = *config
	return &minio, err
}

func (config Config) initMinio() (*minio.Client, error) {

	proxyUrl, _ := url.Parse(config.ProxyUrl)

	// Initialize minio client object.
	minioClient, errInit := minio.New(config.Endpoint, &minio.Options{
		Creds:     credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, config.Token),
		Secure:    config.UseSSL,
		Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)},
	})

	_, err := minioClient.BucketExists(context.Background(), config.Bucket)
	if err != nil {
		return minioClient, err

	}

	return minioClient, errInit

}

func (config Config) FileBase64(fileName string, fileContent string) (*string, error) {
	idx := strings.Index(fileContent, ";base64,")

	if idx < 0 {
		return nil, errors.New("format not match")
	}

	imageType := fileContent[11:idx]

	unbased, err := base64.StdEncoding.DecodeString(fileContent[idx+8:])
	if err != nil {
		return nil, err
	}

	minioClient, errs := config.initMinio()
	if errs != nil {
		return nil, errs
	}

	file := bytes.NewReader(unbased)

	// fileName := fmt.Sprint(uuid.New(), ".", imageType)
	cacheControl := "max-age=31536000"
	userMetaData := map[string]string{"x-amz-acl": "public-read"}

	info, err := minioClient.PutObject(context.Background(), config.Bucket, fileName, file, -1, minio.PutObjectOptions{
		ContentType:  "image/" + imageType,
		CacheControl: cacheControl,
		UserMetadata: userMetaData,
	})

	return &info.Key, err
}

func (config Config) File(fileContent *multipart.FileHeader) (*string, error) {

	// Get Buffer from file
	buffer, err := fileContent.Open()

	if err != nil {
		return nil, err
	}
	defer buffer.Close()

	minioClient, errs := config.initMinio()
	if errs != nil {
		return nil, errs
	}
	ctx := context.Background()

	objectName := fileContent.Filename
	fileBuffer := buffer
	contentType := fileContent.Header["Content-Type"][0]
	fileSize := fileContent.Size
	cacheControl := "max-age=31536000"
	userMetaData := map[string]string{"x-amz-acl": "public-read"}

	info, err := minioClient.PutObject(ctx, "umi", objectName, fileBuffer, fileSize, minio.PutObjectOptions{ContentType: contentType, CacheControl: cacheControl, UserMetadata: userMetaData})

	return &info.Key, err
}

func (config Config) GetFile(fileName string) (*http.Response, error) {
	client := &http.Client{}

	if config.ProxyPublic == "true" {
		proxyUrl, err := url.Parse(config.UrlProxyPublic)
		if err != nil {
			return nil, err
		}
		client = &http.Client{
			Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)},
		}
	}

	urls := config.UrlPublic + "/" + config.Bucket + "/" + fileName

	req, err := http.NewRequestWithContext(context.Background(), "GET", urls, nil)
	if err != nil {
		return nil, err
	}

	resultHttp, err := client.Do(req)

	return resultHttp, err
}
