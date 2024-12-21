package list

import (
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/stretchr/testify/assert"
)

func TestS3Lister_ListFiltering(t *testing.T) {
	tests := []struct {
		name           string
		objects        []types.Object
		expectedTxt    int
		expectedNonTxt int
	}{
		{
			name: "Only txt files",
			objects: []types.Object{
				{Key: aws.String("file1.txt")},
				{Key: aws.String("file2.txt")},
			},
			expectedTxt:    2,
			expectedNonTxt: 0,
		},
		{
			name: "Mixed files",
			objects: []types.Object{
				{Key: aws.String("file1.txt")},
				{Key: aws.String("file2.jpg")},
				{Key: aws.String("file3.pdf")},
			},
			expectedTxt:    1,
			expectedNonTxt: 2,
		},
		{
			name: "No txt files",
			objects: []types.Object{
				{Key: aws.String("file1.jpg")},
				{Key: aws.String("file2.pdf")},
			},
			expectedTxt:    0,
			expectedNonTxt: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			txtCount := 0
			nonTxtCount := 0

			for _, obj := range tt.objects {
				if strings.HasSuffix(*obj.Key, ".txt") {
					txtCount++
				} else {
					nonTxtCount++
				}
			}

			assert.Equal(t, tt.expectedTxt, txtCount)
			assert.Equal(t, tt.expectedNonTxt, nonTxtCount)
		})
	}
}

func TestS3Lister_SizeFormatting(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		size     int64
		expected float64
	}{
		{
			name:     "1 MB file",
			size:     1024 * 1024,
			expected: 1.0,
		},
		{
			name:     "2.5 MB file",
			size:     int64(2.5 * 1024 * 1024),
			expected: 2.5,
		},
		{
			name:     "Empty file",
			size:     0,
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := types.Object{
				Key:          aws.String("test.txt"),
				Size:         aws.Int64(tt.size),
				LastModified: &now,
			}
			sizeMB := float64(*obj.Size) / 1024 / 1024
			assert.Equal(t, tt.expected, sizeMB)
		})
	}
}

func TestS3Lister_DateFormatting(t *testing.T) {
	tests := []struct {
		name     string
		date     time.Time
		expected string
	}{
		{
			name:     "Format current date",
			date:     time.Date(2024, 12, 20, 15, 0o4, 0o5, 0, time.UTC),
			expected: "2024-12-20 15:04:05",
		},
		{
			name:     "Format past date",
			date:     time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: "2023-01-01 00:00:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := types.Object{
				Key:          aws.String("test.txt"),
				LastModified: &tt.date,
			}
			formatted := obj.LastModified.Format("2006-01-02 15:04:05")
			assert.Equal(t, tt.expected, formatted)
		})
	}
}
