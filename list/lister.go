package list

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-logr/logr"
)

type S3Lister struct {
	logger logr.Logger
	bucket string
}

func NewS3Lister(logger logr.Logger, bucket string) *S3Lister {
	return &S3Lister{
		logger: logger,
		bucket: bucket,
	}
}

func (l *S3Lister) List(ctx context.Context) error {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		return fmt.Errorf("unable to load SDK config: %w", err)
	}

	client := s3.NewFromConfig(cfg)

	input := &s3.ListObjectsV2Input{
		Bucket: &l.bucket,
	}

	paginator := s3.NewListObjectsV2Paginator(client, input)

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("error listing objects: %w", err)
		}

		for _, obj := range page.Contents {
			if strings.HasSuffix(*obj.Key, ".txt") {
				sizeMB := float64(*obj.Size) / 1024 / 1024
				fmt.Printf("%-60s  %8.2f MB  %v\n",
					*obj.Key,
					sizeMB,
					obj.LastModified.Format("2006-01-02 15:04:05"),
				)
			}
		}
	}

	return nil
}
