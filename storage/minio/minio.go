package minio

import (
	"net/http"
	"net/url"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Config struct {
	endpoint        string
	accessKeyID     string
	secretAccessKey string
	useSSL          bool
	proxyUrl        string
	token           string
}

func New(config Config) (*minio.Client, error) {

	proxyUrl, _ := url.Parse(config.proxyUrl)

	// Initialize minio client object.
	minioClient, err := minio.New(config.endpoint, &minio.Options{
		Creds:     credentials.NewStaticV4(config.accessKeyID, config.secretAccessKey, config.token),
		Secure:    config.useSSL,
		Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)},
	})
	return minioClient, err
}
