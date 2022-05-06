# sup

sup is a tool for uploading files to S3 in a repeatable manner. Like `aws s3
sync` it only uploads files that have changed. Unlike that tool it uses a
configuration file that allows you to specify per-object metadata and
transformations. It is intended for use in deployment pipelines, especially for
uploading static web content.

## Installing

Prepackaged binaries may be provided in the future. Until then:

1. [Install Go][go-install].
2. Clone the repository
3. Inside the checkout run `go generate ./... && go build ./cmd/sup`

This will generate a static executable called `sup`.

[go-install]: https://go.dev/doc/install

## Usage

Assuming you want to upload files from a directory called 'www', start by
creating  a config file:

```
echo 'source = "www"' > sup.hcl
```

Then run sup with the path to the config file and the URL of an S3 bucket:

```
sup sup.hcl s3://bucket/prefix
```

For more documentation, including all supported configuration options, run `sup --help`

## License

sup is licensed under the MIT license. For information on the licenses of all included modules run `sup --credits`.