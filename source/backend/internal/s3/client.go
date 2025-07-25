package s3

import (
	"context"
	"log"

	"github.com/EquipQR/equipqr/backend/internal/utils"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var Client *minio.Client

func Init() {
	config := utils.LoadConfigFromEnv()
	endpoint := config.MinioEndpoint
	accessKey := config.MinioAccessKey
	secretKey := config.MinioSecretKey
	useSSL := config.MinioUseSSL
	bucket := config.MinioBucket

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalf("Failed to create MinIO client: %v", err)
	}

	exists, err := minioClient.BucketExists(context.Background(), bucket)
	if err != nil {
		log.Fatalf("Error checking bucket: %v", err)
	}
	if !exists {
		err = minioClient.MakeBucket(context.Background(), bucket, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatalf("Error creating bucket: %v", err)
		}
		log.Printf("Created bucket %q\n", bucket)
	}

	Client = minioClient
}
