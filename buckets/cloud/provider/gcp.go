package provider

import (
	"time"

	pifc "github.com/hsuanshao/go-tools/buckets/cloud/provider/interface"
	e "github.com/hsuanshao/go-tools/buckets/entity"
	"github.com/hsuanshao/go-tools/ctx"
)

/**
offical document: https://cloud.google.com/go/docs/reference/cloud.google.com/go/storage/latest
*/

var (
	googleBlobDefaultPath = "https://storage.googleapis.com"
)

func NewGCS(ctx ctx.CTX, conf *e.Config) (gcp pifc.ObjectServiceProvider, err error) {
	return &gcsImpl{}, e.ErrNotImpl
}

type gcsImpl struct {
	bucket string
	region string
}

// GetObjectList to fetch object list
func (gim *gcsImpl) GetObjectList(ctx ctx.CTX, prefix, delim string) (objURLs []string, err error) {
	return nil, e.ErrNotImpl
}

// ReadObjectContent to read object content
func (gim *gcsImpl) ReadObjectContent(ctx ctx.CTX, objectPath string) (objRaw []byte, metadata map[string]string, err error) {
	return nil, nil, e.ErrNotImpl
}

// IsObjectExists to check object existence by given url
func (gim *gcsImpl) IsObjectExists(ctx ctx.CTX, objURL string) (existed bool, err error) {
	return false, e.ErrNotImpl
}

func (gim *gcsImpl) GenReadPresignedURL(ctx ctx.CTX, objURL string, duration time.Duration) (readPresignedURL string, err error) {
	return "", e.ErrNotImpl
}

// PutPresignedURL to generate an object upload with a ttl permision
func (gim *gcsImpl) PutPresignedURL(ctx ctx.CTX, objURL string, mime e.ContentType, duration time.Duration, metaData map[string]string) (presignedURL string, err error) {
	return "", e.ErrNotImpl
}

// Upload to upload object
func (gim *gcsImpl) Upload(ctx ctx.CTX, ct e.ContentType, objpath string, objraw []byte, objmetadata map[string]string) (URL string, readPresignedURL string, err error) {
	return "", "", e.ErrNotImpl
}

// Override to override exists obj, this function will compare original content and new content
func (gim *gcsImpl) Override(ctx ctx.CTX, ct e.ContentType, objPath string, objNewRaw []byte, objmetadata map[string]string) (objURL string, err error) {
	return "", e.ErrNotImpl
}

func (gim *gcsImpl) Delete(ctx ctx.CTX, contentType e.ContentType, objPathes []string) (result bool, err error) {
	return false, e.ErrNotImpl
}

// Health to tell every platform service latency
func (gim *gcsImpl) Health(ctx ctx.CTX) (status e.HealthStatus, err error) {
	return e.HealthStatus{}, e.ErrNotImpl
}

// Close to close client
func (gim *gcsImpl) Close() {}
