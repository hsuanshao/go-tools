package provider

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	awsS3 "github.com/aws/aws-sdk-go/service/s3"
	"github.com/gabriel-vasile/mimetype"
	"github.com/sirupsen/logrus"

	"github.com/hsuanshao/go-tools/aws/s3"
	pifc "github.com/hsuanshao/go-tools/buckets/cloud/provider/interface"
	e "github.com/hsuanshao/go-tools/buckets/entity"
	"github.com/hsuanshao/go-tools/ctx"
)

/*
Due to AWS-SDK-Go-V2 Does not provide interface,
and both v1, v2 sdk are still keep in maintain status
for system stable consideration, here chose aws-sdk-go-v1 as
client sdk
*/

// NewS3 ...
func NewS3(ctx ctx.CTX, conf *e.Config) (s3srv pifc.ObjectServiceProvider, err error) {
	region, bucket := conf.Region, conf.Bucket

	var s3Srv s3.Bucket

	s3Srv = nil

	if conf.Option != nil {
		if conf.Option.AssumeRoleARN != nil && strings.TrimSpace(*conf.Option.AssumeRoleARN) != "" {
			s3Srv = s3.NewS3withAssumeRole(ctx, region, *conf.Option.AssumeRoleARN)
		}
	}

	if s3Srv == nil {
		s3Srv = s3.NewS3(ctx, region)
	}

	return &s3impl{
		s3Srv:   s3Srv,
		bucket:  bucket,
		region:  region,
		isMinio: false,
	}, nil
}

type s3impl struct {
	s3Srv   s3.Bucket
	bucket  string
	region  string
	isMinio bool
}

var (
	// awsS3ObjUrl is object url pattern to aws s3
	// rul is https://{bucket name}.s3.amazonaws.com/{object key}
	awsS3ObjUrl = "https://%s.s3.amazonaws.com/%s"

	// awsS3RegionObjUrl is the pattern while bucket has ".",
	// pattern is https://s3-{region}.amazonaws.com/{bucket name}/{object key}
	awsS3RegionObjUrl = "https://s3-%s.amazonaws.com/%s/%s"
)

// GetObjectList to fetch object list
func (s3Im *s3impl) GetObjectList(ctx ctx.CTX, prefix, delim string) (objURLs []string, err error) {
	// NOTE: please well manage object storage apporch, because pull
	// object list, return up to 1000.
	// to pull object list also cost some money, therefore, to put
	// object without planning is danger, and costly
	listInput := awsS3.ListObjectsV2Input{Bucket: aws.String(s3Im.bucket)}

	if prefix != "" && strings.TrimSpace(prefix) != "" {
		listInput.Prefix = aws.String(prefix)
	}

	if delim != "" && strings.TrimSpace(delim) != "" {
		listInput.Delimiter = aws.String(delim)
	}

	output, err := s3Im.s3Srv.ListObjectsV2(&listInput)
	if err != nil {
		ctx.WithFields(logrus.Fields{"err": err, "bucket": s3Im.bucket, "prefix": prefix, "delim": delim}).Error("list objects from s3 failed")
		return nil, e.ErrFetchObjList
	}

	if len(output.Contents) == 0 {
		ctx.WithFields(logrus.Fields{"prefix": prefix, "delim": delim, "listInput": listInput}).Warn("list object with zero object output content")
	}

	for idx, obj := range output.Contents {
		if obj.Key == nil {
			ctx.WithFields(logrus.Fields{"idx": idx, "obj": obj}).Warn("one of object key is nil")
			continue
		}

		url := s3Im.getObjURL(*obj.Key)
		objURLs = append(objURLs, url)
	}

	return objURLs, nil
}

func (s3Im *s3impl) getObjURL(objKey string) (url string) {
	prefixChk := strings.HasPrefix(objKey, "/")
	if prefixChk {
		objKey = objKey[1:]
	}
	hasDot := strings.ContainsAny(s3Im.bucket, ".")
	switch hasDot {
	case true:
		url = fmt.Sprintf(awsS3RegionObjUrl, s3Im.region, s3Im.bucket, objKey)
	case false:
		url = fmt.Sprintf(awsS3ObjUrl, s3Im.bucket, objKey)
	}
	return url
}

func (s3Im *s3impl) processObjURL(url string) (region, bucket, fileKey string) {
	region = s3Im.region
	bucket = s3Im.bucket

	prefixOk := strings.HasPrefix(url, "https://")
	if prefixOk {
		url = url[8:]
	}
	var urlV1ok, urlV2ok, urlOK bool
	urlOK = false
	urlV1ok = strings.Contains(url, ".s3.amazonaws.com/")
	if !urlV1ok {
		urlV2ok = strings.Contains(url, "amazonaws.com/")
	}
	var splitArr []string
	if urlV2ok || urlV1ok {
		urlOK = true
		if urlV1ok {
			splitArr = strings.Split(url, ".s3.amazonaws.com/")
			bucket = splitArr[0]
			splitArr = strings.Split(splitArr[1], "/")
		}

		if urlV2ok {
			splitArr = strings.Split(url, ".amazonaws.com/")
			region = splitArr[0][3:]
			splitArr = strings.Split(splitArr[1], "/")
			bucket = splitArr[0]
		}

		l := len(splitArr[1:]) - 1
		for n, s := range splitArr[1:] {
			fileKey += s
			if n < l {
				fileKey += "/"
			}
		}
	}

	if !prefixOk && !urlOK {
		fileKey = url
	}

	return region, bucket, fileKey
}

// ReadObjectContent to read object content
func (s3Im *s3impl) ReadObjectContent(ctx ctx.CTX, objectPath string) (objRaw []byte, metadata map[string]string, err error) {
	objRegion, objBucket, objKey := s3Im.processObjURL(objectPath)
	if objRegion != s3Im.region {
		ctx.WithFields(logrus.Fields{"object region": objRegion}).Warn("given object path region is not in permission region")
		return nil, nil, e.ErrWithoutPermissionToAccess
	}

	if objBucket != s3Im.bucket {
		ctx.WithFields(logrus.Fields{"object bucket": objBucket}).Warn("given object path bucket is not in expected bucket")
		return nil, nil, e.ErrWithoutPermissionToAccess
	}

	response, err := s3Im.s3Srv.GetObject(&awsS3.GetObjectInput{Bucket: aws.String(s3Im.bucket), Key: aws.String(objKey)})
	if err != nil {
		ctx.WithFields(logrus.Fields{"err": err, "region": s3Im.region, "bucket": s3Im.bucket, "key": objKey, "response": response}).Error("get object from s3 failed")
		return nil, nil, e.ErrGetObjFromS3
	}

	objMeta := map[string]string{}
	for key, val := range response.Metadata {
		if val != nil {
			objMeta[key] = *val
		}
	}
	contenType := ""
	if response.ContentType != nil {
		contenType = *response.ContentType
	}
	objSizeInByte := int64(0)
	if response.ContentLength != nil {
		objSizeInByte = *response.ContentLength
	}
	ctx.WithFields(logrus.Fields{"content type": contenType, "metadata": objMeta, "objBytes": objSizeInByte}).Info("check object context")

	contentByte, err := io.ReadAll(response.Body)
	if err != nil {
		ctx.WithField("err", err).Error("read response body get error")
		return nil, nil, e.ErrReadObjContent
	}
	defer response.Body.Close()

	return contentByte, objMeta, nil
}

// IsObjectExists to check object existence by given url
func (s3Im *s3impl) IsObjectExists(ctx ctx.CTX, objURL string) (existed bool, err error) {
	objRegion, objBucket, objKey := s3Im.processObjURL(objURL)
	if objRegion != s3Im.region {
		ctx.WithFields(logrus.Fields{"object region": objRegion}).Warn("given object path region is not in permission region")
		return false, e.ErrWithoutPermissionToAccess
	}

	if objBucket != s3Im.bucket {
		ctx.WithFields(logrus.Fields{"object bucket": objBucket}).Warn("given object path bucket is not in expected bucket")
		return false, e.ErrWithoutPermissionToAccess
	}

	opt, err := s3Im.s3Srv.GetObjectAttributes(&awsS3.GetObjectAttributesInput{Bucket: aws.String(s3Im.bucket), Key: aws.String(objKey)})
	if err != nil {
		ctx.WithFields(logrus.Fields{"err": err, "bucket": s3Im.bucket, "objKey": objKey, "objURL": objURL}).Info("get object attribute get error return")
		// NOTE: here doesn't return error is due to if given objURL not exists, get error to this function is a kind of positive result, kept log for tracing, learning
		return false, nil
	}

	exists := false
	if opt.LastModified != nil {
		exists = true
	}

	return exists, nil
}

// GenReadPresginedURL to generate a private object with ttl permission
func (s3Im *s3impl) GenReadPresignedURL(ctx ctx.CTX, objURL string, duration time.Duration) (readPresignedURL string, err error) {
	objRegion, objBucket, objKey := s3Im.processObjURL(objURL)
	if objRegion != s3Im.region {
		ctx.WithFields(logrus.Fields{"object region": objRegion}).Warn("given object path region is not in permission region")
		return "", e.ErrWithoutPermissionToAccess
	}

	if objBucket != s3Im.bucket {
		ctx.WithFields(logrus.Fields{"object bucket": objBucket}).Warn("given object path bucket is not in expected bucket")
		return "", e.ErrWithoutPermissionToAccess
	}

	req, opt := s3Im.s3Srv.GetObjectRequest(&awsS3.GetObjectInput{
		Bucket: aws.String(s3Im.bucket),
		Key:    aws.String(objKey),
	})

	if opt.ContentLength != nil && *opt.ContentLength == 0 {
		// NOTE: the purpose here to capture opt information is try to capture while opt.ContentLength returns zero, try to capture some other information to know more information from s3 response
		ctType := ""
		if opt.ContentType != nil {
			ctType = *opt.ContentType
		}
		ctx.WithFields(logrus.Fields{"objURL": objURL, "bucket": objBucket, "objRegion": objRegion, "content-length": *opt.ContentLength, "opt-content-type": ctType}).Info("fetch object content length length is zero length")
	}

	readPresignedURL, err = req.Presign(15 * time.Minute)
	if err != nil {
		ctx.WithFields(logrus.Fields{"err": err, "bucket": s3Im.bucket, "objectKey": objKey, "objPath": objURL}).Error("try genereate read presigned url failed")
		return "", e.ErrGenReadPresignedURL
	}

	return readPresignedURL, nil
}

// PutPresignedURL to generate an object upload with a ttl permision
func (s3Im *s3impl) PutPresignedURL(ctx ctx.CTX, objURL string, mime e.ContentType, duration time.Duration, metaData map[string]string) (presignedURL string, err error) {
	objRegion, objBucket, objKey := s3Im.processObjURL(objURL)
	if objRegion != s3Im.region {
		ctx.WithFields(logrus.Fields{"object region": objRegion}).Warn("given object path region is not in permission region")
		return "", e.ErrWithoutPermissionToAccess
	}

	if objBucket != s3Im.bucket {
		ctx.WithFields(logrus.Fields{"object bucket": objBucket}).Warn("given object path bucket is not in expected bucket")
		return "", e.ErrWithoutPermissionToAccess
	}

	objExists, _ := s3Im.IsObjectExists(ctx, objURL)
	if objExists {
		ctx.WithFields(logrus.Fields{"objURL": objURL, "metadata": metaData, "content-type": mime}).Warn("try to generate by given obj url, however the url already exists an objects")
		return "", e.ErrObjectPathHasItem
	}

	req, _ := s3Im.s3Srv.PutObjectRequest(&awsS3.PutObjectInput{
		Bucket:      aws.String(s3Im.bucket),
		Key:         aws.String(objKey),
		Body:        strings.NewReader("EXPECTED CONTENTS"),
		ContentType: aws.String(mime.String()),
	})

	presPutURLstr, err := req.Presign(15 * time.Minute)
	if err != nil {
		ctx.WithFields(logrus.Fields{"err": err, "objBucket": objBucket, "objKey": objKey}).Error("generate put presign url failed")
		return "", e.ErrGenPutPresignedURL
	}
	return presPutURLstr, nil
}

// Upload to upload object
func (s3Im *s3impl) Upload(ctx ctx.CTX, ct e.ContentType, objpath string, objraw []byte, objmetadata map[string]string) (URL string, readPresignedURL string, err error) {
	objRegion, objBucket, objKey := s3Im.processObjURL(objpath)
	if objRegion != s3Im.region {
		ctx.WithFields(logrus.Fields{"object region": objRegion}).Warn("given object path region is not in permission region")
		return "", "", e.ErrWithoutPermissionToAccess
	}

	if objBucket != s3Im.bucket {
		ctx.WithFields(logrus.Fields{"object bucket": objBucket}).Warn("given object path bucket is not in expected bucket")
		return "", "", e.ErrWithoutPermissionToAccess
	}

	opt, err := s3Im.s3Srv.GetObjectAttributes(&awsS3.GetObjectAttributesInput{Bucket: aws.String(s3Im.bucket), Key: aws.String(objKey)})
	if err != nil {
		ctx.WithFields(logrus.Fields{"err": err, "bucket": s3Im.bucket, "objectKey": objKey}).Info("fetch object attribuite get error return")
		// NOTE: not return error here, it is due to we are expected it should not have any object at given bucket and obj key
	}

	exists := false
	if opt.LastModified != nil {
		exists = true
	}

	if exists {
		ctx.WithFields(logrus.Fields{"objRegion": objRegion, "objBucket": objBucket, "objKey": objKey, "lastModifiedTime": opt.LastModified.Format(time.RFC3339), "given objPath": objpath}).Warn("given obj path already has object")

		return "", "", e.ErrObjectPathHasItem
	}

	objContentType := mimetype.Detect(objraw)
	if objContentType.String() != ct.String() {
		ctx.WithFields(logrus.Fields{"expect content type": ct, "detect content type": objContentType}).Error("upload new object content type is not match input parameter")
		return "", "", e.ErrUploadNotMatchContentType
	}

	bodyBytes := bytes.NewReader(objraw)
	putOutput, err := s3Im.s3Srv.PutObject(&awsS3.PutObjectInput{
		Bucket:   aws.String(s3Im.bucket),
		Key:      aws.String(objKey),
		Body:     bodyBytes,
		Metadata: aws.StringMap(objmetadata),
	})
	if err != nil {
		ctx.WithFields(logrus.Fields{"err": err, "bucket": s3Im.bucket, "objectPath": objpath, "obj": objKey}).Error("upload object to s3 bucket failed")
		return "", "", e.ErrUploadObjToS3
	}

	verID := ""
	// NOTE: VersionId not nil only while the bucket been settle has versioning control
	if putOutput.VersionId != nil {
		verID = *putOutput.VersionId
	}
	ctx.WithField("objVerID", verID).Info("check object version id")

	objURL := s3Im.getObjURL(objKey)
	readPresignedURL, err = s3Im.GenReadPresignedURL(ctx, objURL, 5*time.Minute)
	if err != nil {
		ctx.WithFields(logrus.Fields{"err": err, "objURL": objURL}).Warn("try to genenerate read presigned url failed from just upload object url")
		// NOTE: return nil err here, is due to some cloud object storeage service might has latency issue to generate presigned url (read)
		return objURL, "", nil
	}

	return objURL, readPresignedURL, nil
}

// Override to override exists obj, this function will compare original content and new content
func (s3Im *s3impl) Override(ctx ctx.CTX, ct e.ContentType, objPath string, objNewRaw []byte, objmetadata map[string]string) (objURL string, err error) {
	objRegion, objBucket, objKey := s3Im.processObjURL(objPath)
	if objRegion != s3Im.region {
		ctx.WithFields(logrus.Fields{"object region": objRegion}).Warn("given object path region is not in permission region")
		return "", e.ErrWithoutPermissionToAccess
	}

	if objBucket != s3Im.bucket {
		ctx.WithFields(logrus.Fields{"object bucket": objBucket}).Warn("given object path bucket is not in expected bucket")
		return "", e.ErrWithoutPermissionToAccess
	}
	// fetch original object information
	response, err := s3Im.s3Srv.GetObject(&awsS3.GetObjectInput{Bucket: aws.String(s3Im.bucket), Key: aws.String(objKey)})
	if err != nil {
		ctx.WithFields(logrus.Fields{"err": err, "region": s3Im.region, "bucket": s3Im.bucket, "key": objKey}).Error("get object from s3 failed")
		return "", e.ErrFetchOriginObjFromS3ByGivenObjPath
	}

	if response.LastModified == nil {
		ctx.WithFields(logrus.Fields{"err": err, "region": s3Im.region, "bucket": s3Im.bucket, "key": objKey, "last modified time": response.LastModified}).Warn("get object from s3 but without obj modified time")
		return "", e.ErrFetchOriginObjFromS3ByGivenObjPath
	}

	originContentType := ""
	if response.ContentType != nil {
		originContentType = *response.ContentType
	}

	var latestModifiedTime time.Time
	if response.LastModified != nil {
		latestModifiedTime = *response.LastModified
		latestModifiedTime.UTC()
	}
	originTimeStr := latestModifiedTime.Format(time.RFC3339)
	ctx.WithFields(logrus.Fields{"origin-content-type": originContentType, "latest-modified-utc-time": originTimeStr}).Info("[record] origin object related information")

	optObjURL, _, err := s3Im.Upload(ctx, ct, objKey, objNewRaw, objmetadata)
	if err != nil {
		ctx.WithFields(logrus.Fields{"err": err, "content-type": ct, "objPath": objPath}).Error("upload new object content failed")
		return "", e.ErrOverrideObject
	}
	return optObjURL, nil
}

// Delete for delete object item from blob object storage
func (s3Im *s3impl) Delete(ctx ctx.CTX, contentType e.ContentType, objPathes []string) (result bool, err error) {
	objIndentifiedKey := []*awsS3.ObjectIdentifier{}
	region := s3Im.region
	bucket := s3Im.bucket
	for _, objPath := range objPathes {
		objRegion, objBucket, objKey := s3Im.processObjURL(objPath)
		if objRegion != region {
			ctx.WithFields(logrus.Fields{"object region": objRegion}).Warn("given object path region is not in permission region")
			return false, e.ErrWithoutPermissionToAccess
		}

		if objBucket != bucket {
			ctx.WithFields(logrus.Fields{"object bucket": objBucket}).Warn("given object path bucket is not in expected bucket")
			return false, e.ErrWithoutPermissionToAccess
		}
		deleteKey := awsS3.ObjectIdentifier{
			Key: aws.String(objKey),
		}
		objIndentifiedKey = append(objIndentifiedKey, &deleteKey)
	}

	resOpt, err := s3Im.s3Srv.DeleteObjects(&awsS3.DeleteObjectsInput{
		Bucket: aws.String(s3Im.bucket),
		Delete: &awsS3.Delete{
			Objects: objIndentifiedKey,
			Quiet:   aws.Bool(false),
		},
	})
	if err != nil {
		ctx.WithFields(logrus.Fields{"err": err}).Error("delete object from s3 failed")
		return false, e.ErrDeleteObject
	}

	if len(resOpt.Deleted) != len(objIndentifiedKey) {
		ctx.WithFields(logrus.Fields{"response output delete object length": len(resOpt.Deleted), "input objects": len(objIndentifiedKey)}).Warn("delete object count doesn't matching")
	}

	return true, nil
}

// Health to tell every platform service latency
func (s3Im *s3impl) Health(ctx ctx.CTX) (status e.HealthStatus, err error) {
	now := time.Now()
	endpoint := "https://health." + s3Im.region + "amazonaws.com"
	timeout := time.Duration(30 * time.Second)
	_, err = net.DialTimeout("tcp", endpoint, timeout)
	if err != nil {
		return e.HealthStatus{Cloud: e.AWS, Latency: timeout}, e.ErrS3HealthTimeOut
	}

	spendTime := time.Since(now)
	return e.HealthStatus{
		Cloud:   e.AWS,
		Latency: spendTime,
	}, nil
}

// Close to close client
func (s3Im *s3impl) Close() {
	// NOTE: s3Im.s3Srv is a aws sesssion, not a s3 client connection
	// it has no needed to close it.
}
