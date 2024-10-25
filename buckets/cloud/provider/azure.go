package provider

import (
	"time"

	pifc "github.com/hsuanshao/go-tools/buckets/cloud/provider/interface"
	e "github.com/hsuanshao/go-tools/buckets/entity"
	"github.com/hsuanshao/go-tools/ctx"
)

/**
azure sdk https://github.com/Azure/azure-sdk-for-go
*/

func NewAzure(ctx ctx.CTX, conf *e.Config) (azure pifc.ObjectServiceProvider, err error) {
	return &msImpl{}, e.ErrNotImpl
}

type msImpl struct {
	bucket string
	region string
}

// GetObjectList to fetch object list
func (gim *msImpl) GetObjectList(ctx ctx.CTX, prefix, delim string) (objURLs []string, err error) {
	return nil, e.ErrNotImpl
}

// ReadObjectContent to read object content
func (gim *msImpl) ReadObjectContent(ctx ctx.CTX, objectPath string) (objRaw []byte, metadata map[string]string, err error) {
	return nil, nil, e.ErrNotImpl
}

// IsObjectExists to check object existence by given url
func (gim *msImpl) IsObjectExists(ctx ctx.CTX, objURL string) (existed bool, err error) {
	return false, e.ErrNotImpl
}

func (gim *msImpl) GenReadPresignedURL(ctx ctx.CTX, objURL string, duration time.Duration) (readPresignedURL string, err error) {
	return "", e.ErrNotImpl
}

// PutPresignedURL to generate an object upload with a ttl permision
func (gim *msImpl) PutPresignedURL(ctx ctx.CTX, objURL string, mime e.ContentType, duration time.Duration, metaData map[string]string) (presignedURL string, err error) {
	return "", e.ErrNotImpl
}

// Upload to upload object
func (gim *msImpl) Upload(ctx ctx.CTX, ct e.ContentType, objpath string, objraw []byte, objmetadata map[string]string) (URL string, readPresignedURL string, err error) {
	return "", "", e.ErrNotImpl
}

// Override to override exists obj, this function will compare original content and new content
func (gim *msImpl) Override(ctx ctx.CTX, ct e.ContentType, objPath string, objNewRaw []byte, objmetadata map[string]string) (objURL string, err error) {
	return "", e.ErrNotImpl
}

func (gim *msImpl) Delete(ctx ctx.CTX, contentType e.ContentType, objPathes []string) (result bool, err error) {
	return false, e.ErrNotImpl
}

// Health to tell every platform service latency
func (gim *msImpl) Health(ctx ctx.CTX) (status e.HealthStatus, err error) {
	return e.HealthStatus{}, e.ErrNotImpl
}

// Close to close client
func (gim *msImpl) Close() {}
