package fetch

import (
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestS3Fetcher_FilterMissingFiles(t *testing.T) {
	logger := logr.Discard()
	fetcher := NewS3Fetcher(logger, "test-bucket")

	// Create a temporary file
	existingFile := "existing.txt"
	err := os.WriteFile(existingFile, []byte("test"), 0o644)
	require.NoError(t, err)
	defer os.Remove(existingFile)

	tests := []struct {
		name     string
		objects  []types.Object
		expected int
	}{
		{
			name:     "Empty list returns empty",
			objects:  []types.Object{},
			expected: 0,
		},
		{
			name: "Filters only missing files",
			objects: []types.Object{
				{Key: aws.String("existing.txt")},
				{Key: aws.String("missing.txt")},
			},
			expected: 1, // Only missing.txt should be counted
		},
		{
			name: "Handles nested paths",
			objects: []types.Object{
				{Key: aws.String("data/nested/file.txt")},
			},
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fetcher.filterMissingFiles(tt.objects)
			assert.Equal(t, tt.expected, len(result))

			if tt.name == "Filters only missing files" {
				// Verify that the missing file is the one we expect
				require.Len(t, result, 1)
				assert.Equal(t, "missing.txt", *result[0].Key)
			}
		})
	}
}

// Rest of the file remains unchanged...
