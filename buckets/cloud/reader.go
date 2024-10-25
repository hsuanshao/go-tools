package cloud

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/hsuanshao/go-tools/buckets/cloud/provider"
	cpi "github.com/hsuanshao/go-tools/buckets/cloud/provider/interface"
	e "github.com/hsuanshao/go-tools/buckets/entity"
	ifc "github.com/hsuanshao/go-tools/buckets/interface"
	"github.com/hsuanshao/go-tools/ctx"
)

// NewReader is the constructor to launch bucket reader
func NewReader(ctx ctx.CTX, conf *e.Config) (reader ifc.BucketReader, err error) {
	var blobService cpi.ObjectServiceProvider

	switch conf.CSP {
	case e.AWS:
		blobService, err = provider.NewS3(ctx, conf)
	case e.Minio:
		blobService, err = provider.NewMinio(ctx, conf)
	case e.Azure:
		blobService, err = provider.NewAzure(ctx, conf)
	case e.GCP:
		blobService, err = provider.NewGCS(ctx, conf)
	default:
		err = e.ErrInvalidateCSP
	}

	if err != nil {
		ctx.WithFields(logrus.Fields{"err": err, "csp": conf.CSP}).Error("initial blob services by given config failed")
		return nil, err
	}

	return &readImpl{
		cloud: blobService,
	}, nil
}

type readImpl struct {
	cloud cpi.ObjectServiceProvider
}

// ListObjURLs to list out objects storage URLs by given prefix and delim
func (rim *readImpl) ListObjURLs(ctx ctx.CTX, prefix, delim string) (objURLs []string, err error) {
	objectURLs, err := rim.cloud.GetObjectList(ctx, prefix, delim)
	if err != nil {
		ctx.WithFields(logrus.Fields{"err": err, "prefix": prefix, "delim": delim}).Error("get object list failed")
		return nil, e.ErrFetchObjList
	}
	return objectURLs, nil
}

// ReadObjectContent to get object content by given object url
func (rim *readImpl) ReadObjectContent(ctx ctx.CTX, objURL string) (objRawByte []byte, metadata map[string]string, err error) {
	contentBytes, objmetadata, err := rim.cloud.ReadObjectContent(ctx, objURL)
	if err != nil {
		ctx.WithFields(logrus.Fields{"err": err, "objURL": objURL}).Error("unable read object content")
		return nil, nil, e.ErrReadObjContent
	}
	return contentBytes, objmetadata, nil
}

// GenerateReadPresignedURL to generate a "read" permision within a "period" (by given duration),
// with a specific object storage path, if given object path/URL not fould it will return an error
func (rim *readImpl) GenerateReadPresignedURL(ctx ctx.CTX, duration time.Duration, objURL string) (presignedURL string, err error) {
	readPersignedURL, err := rim.cloud.GenReadPresignedURL(ctx, objURL, duration)
	if err != nil {
		ctx.WithFields(logrus.Fields{"err": err, "objURL": objURL, "duration": duration}).Error("generate presigned url failed")
		return "", e.ErrGenReadPresignedURL
	}

	return readPersignedURL, nil
}

// Close to Close BucketReader agent
func (rim *readImpl) Close() {
	rim.cloud.Close()
}
