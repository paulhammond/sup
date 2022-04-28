package remote

import (
	"context"
	"errors"
	"io"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/paulhammond/sup/internal/object"
)

var _ object.Object = s3Object{}

type s3Object struct {
	remote   s3Remote
	path     string
	hash     object.Hash
	metadata object.Metadata
}

func (o s3Object) Hash() (*object.Hash, error) {
	return &o.hash, nil
}

func (o s3Object) Metadata() (*object.Metadata, error) {
	return &o.metadata, nil
}

func (o s3Object) Open(func(io.Reader) error) error {
	panic("unimplemented")
}

type s3Remote struct {
	bucket string
	prefix string
	client *s3.Client
}

func openS3(ctx context.Context, spec string) (Remote, error) {

	u, err := url.Parse(spec)
	if err != nil {
		return nil, err
	}
	if u.Scheme != "s3" {
		return nil, errors.New("inconsistent scheme")
	}
	bucket := u.Host
	prefix := normalizePrefix(u.Path)

	region, err := getRegion(ctx, u.Host)
	if err != nil {
		return nil, err
	}

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)
	if err != nil {
		return nil, err
	}

	return &s3Remote{bucket: bucket, prefix: prefix, client: client}, nil
}

func normalizePrefix(p string) string {
	p = strings.TrimSuffix(p, "/") + "/"
	p = strings.TrimPrefix(p, "/")
	return p
}

func (r *s3Remote) Close() error {
	return nil
}

func (r s3Remote) Set(ctx context.Context) (object.Set, error) {

	input := &s3.ListObjectsV2Input{
		Bucket:  &r.bucket,
		Prefix:  &r.prefix,
		MaxKeys: 1000,
	}

	set := object.Set{}

	paginator := s3.NewListObjectsV2Paginator(r.client, input)
	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, item := range resp.Contents {
			path := *item.Key
			path = strings.TrimPrefix(path, r.prefix)

			hash, err := r.getHash(ctx, path, item)
			if err != nil {
				return nil, err
			}

			set[path] = s3Object{
				remote: r,
				path:   path,
				hash:   *hash,
			}
		}

	}

	return set, nil
}

func (r *s3Remote) getHash(ctx context.Context, path string, item types.Object) (*object.Hash, error) {
	etag := strings.Trim(*item.ETag, "\"")

	// simple upload
	if !strings.Contains(etag, "-") {
		return &object.Hash{
			Size:     *&item.Size,
			PartSize: 0,
			Hash:     etag,
		}, nil
	}

	// multipart upload
	input := &s3.HeadObjectInput{
		Bucket:     aws.String(r.bucket),
		Key:        aws.String(path),
		PartNumber: 1,
	}

	result, err := r.client.HeadObject(ctx, input)
	if err != nil {
		return nil, err
	}
	partSize := result.ContentLength

	return &object.Hash{
		Size:     *&item.Size,
		PartSize: partSize,
		Hash:     etag,
	}, nil
}

func (r *s3Remote) Upload(ctx context.Context, set object.Set, f func(Event)) error {
	mgr := manager.NewUploader(r.client)

	for _, p := range set.Paths() {
		err := r.uploadFile(ctx, mgr, set, p, f)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *s3Remote) uploadFile(ctx context.Context, mgr *manager.Uploader, set object.Set, path string, f func(Event)) error {
	timer := Timer()
	err := set[path].Open(func(file io.Reader) error {

		key := r.prefix + path
		meta, err := set[path].Metadata()
		if err != nil {
			return err
		}

		_, err = mgr.Upload(ctx, &s3.PutObjectInput{
			Bucket:                  &r.bucket,
			Key:                     &key,
			Body:                    file,
			ContentType:             meta.ContentType,
			WebsiteRedirectLocation: meta.WebsiteRedirectLocation,
		})

		return err
	})

	f(Event{Upload, path, set[path], timer()})
	return err

}

func getRegion(ctx context.Context, bucket string) (string, error) {

	tmpCfg, err := config.LoadDefaultConfig(ctx, config.WithDefaultRegion("us-east-1"))
	if err != nil {
		return "", err
	}

	return manager.GetBucketRegion(ctx, s3.NewFromConfig(tmpCfg), "supdev")

}
