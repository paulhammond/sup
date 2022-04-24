package testutil

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
)

func S3Remote(t testing.TB) string {
	if testing.Short() {
		t.Skip("Skipping AWS tests in short mode")
		return ""
	}

	bucket := os.Getenv("SUP_BUCKET")
	if bucket == "" {
		t.Skip("Add SUP_BUCKET environment variable to run this test")
		return ""
	}

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err == nil {
		_, err = cfg.Credentials.Retrieve(ctx)
	}
	if err != nil {
		t.Fatal("error using AWS credentials:", err)
	}

	uniq := time.Now().Format(time.RFC3339Nano)
	return fmt.Sprintf("s3://%s/%s", bucket, uniq)
}
