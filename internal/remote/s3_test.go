package remote_test

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/paulhammond/sup/internal/object"
	"github.com/paulhammond/sup/internal/remote"
	"github.com/paulhammond/sup/internal/remote/testutil"
)

var update = flag.Bool("update", false, "update golden files")

func TestS3(t *testing.T) {

	ctx := context.Background()
	url := testutil.S3Remote(t)

	createTestData(t, "testdata/content", url)
	r, err := remote.Open(ctx, url)
	ok(t, err, "Open")
	defer r.Close()

	set, err := r.Set(ctx)
	ok(t, err, "Set")

	compareSet(t, set, "s3-list")

	toUpload := object.Set{
		"c.txt": object.NewString("c"),
	}

	err = r.Upload(ctx, toUpload, func(remote.Event) {})
	ok(t, err, "Upload")

	r2, err := remote.Open(ctx, url)
	ok(t, err, "Open")
	defer r2.Close()

	set2, err := r.Set(ctx)
	ok(t, err, "Set")
	compareSet(t, set2, "s3-postupload")
}

func createTestData(t *testing.T, local, remote string) {
	cmd := exec.Command("aws", "s3", "sync", local, remote, "--delete")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("error uploading files: %s\n%s", err, out)
	}
	if testing.Verbose() {
		t.Logf("aws s3 sync:\n%s", out)
	}
}

func compareSet(t *testing.T, set object.Set, name string) {
	content := []byte(set.String())

	golden := fmt.Sprintf("testdata/%s.golden", name)
	if *update {
		os.WriteFile(golden, []byte(content), 0644)
	}
	expected, err := os.ReadFile(golden)
	ok(t, err, "Read Golden")
	if !bytes.Equal(content, expected) {
		t.Errorf("%s doesn't match\ngot %s\nexp %s", name, content, expected)
	}
}
