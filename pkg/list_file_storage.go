package pkg

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func ListFiles(client *s3.Client, bucket string) ([]types.Object, error) {
	objects := []types.Object{}
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	}

	paginator := s3.NewListObjectsV2Paginator(client, input)

	// Iterate through the pages of results
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			return objects, fmt.Errorf("failed to get page of results: %w", err)
		}

		objects = append(objects, page.Contents...)
	}

	return objects, nil
}
