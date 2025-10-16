package runner

import (
	"context"
	"fmt"

	"cloud.google.com/go/storage"
)

func uploadToGCS(ctx context.Context, bucket, object string, data []byte) error {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create GCS client: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, gcsUploadTimeout)
	defer cancel()

	obj := client.Bucket(bucket).Object(object)
	wc := obj.NewWriter(ctx)

	if _, err := wc.Write(data); err != nil {
		return fmt.Errorf("writing data to GCS: %w", err)
	}

	if err := wc.Close(); err != nil {
		return fmt.Errorf("closing GCS writer: %w", err)
	}

	return nil
}
