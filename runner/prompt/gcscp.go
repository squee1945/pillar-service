// This package is a tiny GCS object downloader binary.
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"cloud.google.com/go/storage"
)

const (
	gcsDownloadTimeout = 30 * time.Second
)

var (
	gcsPath   = flag.String("gcs-path", "", "GCS path (e.g., gs://bucket/object)")
	localPath = flag.String("local-path", "", "Local file path")
)

func main() {
	ctx := context.Background()
	flag.Parse()

	if *gcsPath == "" || *localPath == "" {
		fail("Usage: gcscp -gcs-path gs://bucket/object -local-path /local/path/to/file")
	}

	gcsURL, err := url.Parse(*gcsPath)
	if err != nil {
		fail("Parsing GCS path: %v", err)
	}
	if gcsURL.Scheme != "gs" {
		fail("Invalid GCS path scheme, expected 'gs://'")
	}
	if gcsURL.Host == "" {
		fail("Invalid GCS path, missing bucket name")
	}
	if gcsURL.Path == "" {
		fail("Invalid GCS path, missing object name")
	}

	if err := copyGCSFileToLocal(ctx, gcsURL, *localPath); err != nil {
		fail("Error: %v\n", err)
	}
	fmt.Printf("Copied %s to %s\n", *gcsPath, *localPath)
}

func fail(f string, args ...any) {
	fmt.Fprintf(os.Stderr, f, args...)
	os.Exit(1)
}

func copyGCSFileToLocal(ctx context.Context, gcs *url.URL, localPath string) error {
	bucketName := gcs.Host
	objectName := gcs.Path
	if objectName[0] == '/' {
		objectName = objectName[1:]
	}

	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("creating GCS client: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, gcsDownloadTimeout)
	defer cancel()

	rc, err := client.Bucket(bucketName).Object(objectName).NewReader(ctx)
	if err != nil {
		return fmt.Errorf("creating GCS object reader: %w", err)
	}
	defer rc.Close()

	if err := os.MkdirAll(filepath.Dir(localPath), 0x777); err != nil {
		return fmt.Errorf("creating local directory: %w", err)
	}

	localFile, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("creating local file: %w", err)
	}
	defer localFile.Close()

	if _, err := io.Copy(localFile, rc); err != nil {
		return fmt.Errorf("copying data to local file: %w", err)
	}
	return nil
}
