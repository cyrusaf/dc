package source

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
)

type Parser[T any] interface {
	Parse(b []byte) (T, error)
}

type S3[T any] struct {
	BucketName   string
	Key          string
	S3Downloader s3manageriface.DownloaderAPI
	Parser       Parser[T]
}

func (s S3[T]) Read(ctx context.Context) (T, error) {
	b := bytes.Buffer{}
	_, err := s.S3Downloader.DownloadWithContext(ctx, fakeWriterAt{
		W: &b,
	}, &s3.GetObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(s.Key),
	})
	if err != nil {
		var t T
		return t, fmt.Errorf("downloading object from s3: %w", err)
	}

	data, err := s.Parser.Parse(b.Bytes())
	if err != nil {
		var t T
		return t, fmt.Errorf("parsing bytes: %w", err)
	}

	return data, nil
}

type fakeWriterAt struct {
	W io.Writer
}

func (fw fakeWriterAt) WriteAt(p []byte, offset int64) (n int, err error) {
	// ignore 'offset' because we forced sequential downloads
	return fw.W.Write(p)
}
