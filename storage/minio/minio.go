package minio

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"

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

func (m Minio) UploadMultiPartFile(fileName, bucket string, fileContent *multipart.FileHeader) (string, error) {

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

	info, err := minioClient.PutObject(ctx, bucket, objectName, fileBuffer, fileSize, minio.PutObjectOptions{})

	return info.Key, err
}

func (config Config) GetFile(url string, bucket string, fileName string) (*http.Response, error) {
	client := &http.Client{}

	// if config.UseProxyPublic {
	// 	proxyUrl, err := url.Parse(config.UrlProxyPublic)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	client = &http.Client{
	// 		Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)},
	// 	}
	// }

	urls := url + "/" + bucket + "/" + fileName

	req, err := http.NewRequestWithContext(context.Background(), "GET", urls, nil)
	if err != nil {
		return nil, err
	}

	resultHttp, err := client.Do(req)

	return resultHttp, err
}
