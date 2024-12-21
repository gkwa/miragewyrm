package fetch

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/go-logr/logr"
)

type S3Fetcher struct {
	logger logr.Logger
	client *s3.Client
	bucket string
	outDir string
}

func NewS3Fetcher(logger logr.Logger, bucket string) *S3Fetcher {
	return &S3Fetcher{
		logger: logger,
		bucket: bucket,
		outDir: ".",
	}
}

func (f *S3Fetcher) SetOutputDir(dir string) {
	f.outDir = dir
}

func (f *S3Fetcher) getLocalPath(key string) string {
	return filepath.Join(f.outDir, key)
}

func (f *S3Fetcher) initClient(ctx context.Context) error {
	if f.client != nil {
		return nil
	}

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		return fmt.Errorf("unable to load SDK config: %w", err)
	}

	f.client = s3.NewFromConfig(cfg)
	return nil
}

func (f *S3Fetcher) FetchRandom(ctx context.Context, count int) error {
	if err := f.initClient(ctx); err != nil {
		return err
	}

	f.logger.Info("Fetching manifest from S3...")
	objects, err := f.listAllObjects(ctx)
	if err != nil {
		return err
	}
	f.logger.Info("Found objects in bucket", "total", len(objects))

	missingObjects := f.filterMissingFiles(objects)
	f.logger.Info("Filtered for missing files", "missing", len(missingObjects), "total", len(objects))

	f.logger.Info("Selecting random files", "requested", count, "available", len(missingObjects))
	selectedObjects := f.selectRandomFiles(missingObjects, count)
	f.logger.Info("Selected files for download", "count", len(selectedObjects))

	for _, obj := range selectedObjects {
		if err := f.downloadFile(ctx, obj); err != nil {
			f.logger.Error(err, "Failed to download file", "key", aws.ToString(obj.Key))
			continue
		}
	}

	return nil
}

func (f *S3Fetcher) listAllObjects(ctx context.Context) ([]types.Object, error) {
	var allObjects []types.Object

	input := &s3.ListObjectsV2Input{
		Bucket: &f.bucket,
	}

	paginator := s3.NewListObjectsV2Paginator(f.client, input)

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("error listing objects: %w", err)
		}

		allObjects = append(allObjects, page.Contents...)
	}

	return allObjects, nil
}

func (f *S3Fetcher) filterMissingFiles(objects []types.Object) []types.Object {
	var missing []types.Object

	for _, obj := range objects {
		localPath := f.getLocalPath(aws.ToString(obj.Key))
		if _, err := os.Stat(localPath); os.IsNotExist(err) {
			missing = append(missing, obj)
		}
	}

	return missing
}

func (f *S3Fetcher) selectRandomFiles(objects []types.Object, count int) []types.Object {
	if len(objects) <= count {
		return objects
	}

	indices := make([]int, len(objects))
	for i := range indices {
		indices[i] = i
	}

	rand.Shuffle(len(indices), func(i, j int) {
		indices[i], indices[j] = indices[j], indices[i]
	})

	selected := make([]types.Object, count)
	for i := 0; i < count; i++ {
		selected[i] = objects[indices[i]]
	}

	return selected
}

func (f *S3Fetcher) downloadFile(ctx context.Context, obj types.Object) error {
	key := aws.ToString(obj.Key)
	localPath := f.getLocalPath(key)
	f.logger.Info("Downloading file", "key", key, "to", localPath)

	if err := os.MkdirAll(filepath.Dir(localPath), 0o755); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	result, err := f.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(f.bucket),
		Key:    obj.Key,
	})
	if err != nil {
		return fmt.Errorf("failed to get object: %w", err)
	}
	defer result.Body.Close()

	file, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	written, err := io.Copy(file, result.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	f.logger.Info("Downloaded file",
		"key", key,
		"to", localPath,
		"size_bytes", written,
	)

	return nil
}
