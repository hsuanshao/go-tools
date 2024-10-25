package provider

import (
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	awsS3 "github.com/aws/aws-sdk-go/service/s3"

	s3api "github.com/hsuanshao/go-tools/aws/s3/mocks"
	cpIfc "github.com/hsuanshao/go-tools/buckets/cloud/provider/interface"
	e "github.com/hsuanshao/go-tools/buckets/entity"
	"github.com/hsuanshao/go-tools/ctx"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

/**
s3 & minio bucket unit test
*/

var (
	mockCTX       = ctx.Background()
	defaultRegion = "ap-northeast-1"
	defaultBucket = "btq-bucket"
	defaultS3api  = new(s3api.Bucket)
)

type s3TestSuite struct {
	suite.Suite
	CSP cpIfc.ObjectServiceProvider
}

func TestS3Suite(t *testing.T) {
	ctx.SetDebugLevel()
	suite.Run(t, new(s3TestSuite))
}

func (s *s3TestSuite) SetupTest() {
	s.CSP = &s3impl{
		s3Srv:  defaultS3api,
		region: defaultRegion,
		bucket: defaultBucket,
	}
}

// TearDownSuite please applied to in the end of any test case to reset
// test suite back to default status
func (s *s3TestSuite) TearDownSuite() {
	mockCTX = ctx.Background()
	s.CSP = &s3impl{
		s3Srv:  defaultS3api,
		region: defaultRegion,
		bucket: defaultBucket,
	}
}

// Test GetObjectList
func (s *s3TestSuite) TestGetObjectList() {
	testcase := []struct {
		Case     string
		MockFunc func()
		Prefix   string
		Delim    string
		ExpRes   []string
		ExpErr   error
	}{
		{
			Case: "normal case, has prefix, without delim",
			MockFunc: func() {
				defaultS3api.On("ListObjectsV2", &awsS3.ListObjectsV2Input{
					Bucket: aws.String(defaultBucket),
					Prefix: aws.String("config"),
				}).Return(&awsS3.ListObjectsV2Output{
					Contents: []*awsS3.Object{
						{
							Key: aws.String("/config/indexer_setup.json"),
						},
						{
							Key: aws.String("/config/buckets.json"),
						},
					},
				}, nil).Once()
			},
			Prefix: "config",
			Delim:  "",
			ExpRes: []string{
				"https://btq-bucket.s3.amazonaws.com/config/indexer_setup.json",
				"https://btq-bucket.s3.amazonaws.com/config/buckets.json",
			},
			ExpErr: nil,
		},
		{
			Case: "normal case, has prefix, exclusive indexer",
			MockFunc: func() {
				defaultS3api.On("ListObjectsV2", &awsS3.ListObjectsV2Input{
					Bucket:    aws.String(defaultBucket),
					Prefix:    aws.String("config"),
					Delimiter: aws.String("indexer"),
				}).Return(&awsS3.ListObjectsV2Output{
					Contents: []*awsS3.Object{
						{
							Key: aws.String("/config/buckets.json"),
						},
					},
				}, nil).Once()
			},
			Prefix: "config",
			Delim:  "indexer",
			ExpRes: []string{
				"https://btq-bucket.s3.amazonaws.com/config/buckets.json",
			},
			ExpErr: nil,
		},
		{
			Case: "use a new bucket client",
			MockFunc: func() {
				s.CSP = &s3impl{
					s3Srv:  defaultS3api,
					region: "ap-south-1",
					bucket: "btq.ag",
				}

				defaultS3api.On("ListObjectsV2", &awsS3.ListObjectsV2Input{
					Bucket: aws.String("btq.ag"),
					Prefix: aws.String("config"),
				}).Return(&awsS3.ListObjectsV2Output{
					Contents: []*awsS3.Object{
						{
							Key: aws.String("/config/indexer_setup.json"),
						},
						{
							Key: aws.String("/config/buckets.json"),
						},
					},
				}, nil).Once()
			},
			Prefix: "config",
			Delim:  "",
			ExpRes: []string{
				"https://s3-ap-south-1.amazonaws.com/btq.ag/config/indexer_setup.json",
				"https://s3-ap-south-1.amazonaws.com/btq.ag/config/buckets.json",
			},
			ExpErr: nil,
		},
		{
			Case: "use a new bucket client, but given wrong bucket",
			MockFunc: func() {
				s.CSP = &s3impl{
					s3Srv:  defaultS3api,
					region: "ap-south-1",
					bucket: defaultBucket,
				}

				defaultS3api.On("ListObjectsV2", &awsS3.ListObjectsV2Input{
					Bucket: aws.String(defaultBucket),
					Prefix: aws.String("config"),
				}).Return(nil, fmt.Errorf("NoSuchBucket, %v", defaultBucket)).Once()
			},
			Prefix: "config",
			Delim:  "",
			ExpRes: nil,
			ExpErr: e.ErrFetchObjList,
		},
	}

	for idx, c := range testcase {
		mockCTX = ctx.WithValue(mockCTX, "caes no", idx)
		mockCTX = ctx.WithValue(mockCTX, "case name", c.Case)
		// run mock func
		c.MockFunc()

		objURLs, err := s.CSP.GetObjectList(mockCTX, c.Prefix, c.Delim)
		s.Equal(c.ExpErr, err, fmt.Sprintf("Case[%v]%v: error not as expected", idx, c.Case))

		s.Equal(len(c.ExpRes), len(objURLs), fmt.Sprintf("Case[%v]%v: compare return urls slice length", idx, c.Case))

		for urlIDX, url := range c.ExpRes {
			s.Equal(url, objURLs[urlIDX], fmt.Sprintf("[C]return urls on %v, get not match url, expected URL is %v, but get %v", urlIDX, url, objURLs[urlIDX]))
		}

		// Note: tear down suite to help reset test suite return to original test setup
		s.TearDownSuite()
	}
}

func (s *s3TestSuite) TestReadObjectContent() {
	testcase := []struct {
		Case        string
		MockFunc    func()
		ObjURL      string
		ExpRes      []byte
		ExpMetadata map[string]string
		ExpErr      error
	}{
		{
			Case: "Normal case",
			MockFunc: func() {
				stringRead := strings.NewReader(`[{"cloud":"aws","purpose":"i18n","region":"ap-northeast-1","bucket":"btq.system.i18n"},{"cloud":"aws","purpose":"config","region":"ap-northeast-1","bucket":"btq.system.conf"}]`)
				stringRC := io.NopCloser(stringRead)
				defaultS3api.On("GetObject", &awsS3.GetObjectInput{Bucket: aws.String(defaultBucket), Key: aws.String("/config/bucket.json")}).Return(&awsS3.GetObjectOutput{
					ContentType: aws.String("application/json"),
					Metadata: map[string]*string{
						"for-system": aws.String("indexer"),
						"version":    aws.String("v1.0"),
					},
					ContentLength: aws.Int64(175),
					Body:          stringRC,
				}, nil).Once()
			},
			ObjURL: "/config/bucket.json",
			ExpRes: []byte{91, 123, 34, 99, 108, 111, 117, 100, 34, 58, 34, 97, 119, 115, 34, 44, 34, 112, 117, 114, 112, 111, 115, 101, 34, 58, 34, 105, 49, 56, 110, 34, 44, 34, 114, 101, 103, 105, 111, 110, 34, 58, 34, 97, 112, 45, 110, 111, 114, 116, 104, 101, 97, 115, 116, 45, 49, 34, 44, 34, 98, 117, 99, 107, 101, 116, 34, 58, 34, 98, 116, 113, 46, 115, 121, 115, 116, 101, 109, 46, 105, 49, 56, 110, 34, 125, 44, 123, 34, 99, 108, 111, 117, 100, 34, 58, 34, 97, 119, 115, 34, 44, 34, 112, 117, 114, 112, 111, 115, 101, 34, 58, 34, 99, 111, 110, 102, 105, 103, 34, 44, 34, 114, 101, 103, 105, 111, 110, 34, 58, 34, 97, 112, 45, 110, 111, 114, 116, 104, 101, 97, 115, 116, 45, 49, 34, 44, 34, 98, 117, 99, 107, 101, 116, 34, 58, 34, 98, 116, 113, 46, 115, 121, 115, 116, 101, 109, 46, 99, 111, 110, 102, 34, 125, 93},
			ExpMetadata: map[string]string{
				"for-system": "indexer",
				"version":    "v1.0",
			},
			ExpErr: nil,
		},
		{
			Case: "Object not exits case",
			MockFunc: func() {

				defaultS3api.On("GetObject", &awsS3.GetObjectInput{Bucket: aws.String(defaultBucket), Key: aws.String("/config/bucket-special.json")}).Return(&awsS3.GetObjectOutput{
					ContentType:   nil,
					Metadata:      nil,
					ContentLength: aws.Int64(0),
					Body:          nil,
				}, fmt.Errorf("any error")).Once()
			},
			ObjURL:      "/config/bucket-special.json",
			ExpRes:      nil,
			ExpMetadata: nil,
			ExpErr:      e.ErrGetObjFromS3,
		},
	}

	for idx, c := range testcase {
		mockCTX = ctx.WithValue(mockCTX, "caes no", idx)
		mockCTX = ctx.WithValue(mockCTX, "case name", c.Case)
		c.MockFunc()

		respBytes, metadata, err := s.CSP.ReadObjectContent(mockCTX, c.ObjURL)

		if !s.Equal(c.ExpRes, respBytes) {
			mockCTX.Error("expected result not match")
		}

		if !s.Equal(len(c.ExpMetadata), len(metadata)) {
			mockCTX.Error("expected metadata of object header length not match")
		}

		for key, val := range c.ExpMetadata {
			if !s.Equal(val, metadata[key]) {
				mockCTX.WithFields(logrus.Fields{"key": key, "metadata val": metadata[key]}).Error("the key in metadata return is not match expected result")
			}
		}

		s.Equal(c.ExpErr, err, fmt.Sprintf("Case[%v]%v, expected error not match", idx, c.Case))
		s.TearDownSuite()
	}
}

func (s *s3TestSuite) TestIsObjectExists() {
	testcase := []struct {
		Case     string
		MockFunc func()
		ObjURL   string
		ExpRes   bool
		ExpErr   error
	}{
		{
			Case: "exists case",
			MockFunc: func() {
				mockModifiedTimeStr := "2022-11-01T12:00:00-00:00"
				mockTimeObj, _ := time.Parse(time.RFC3339, mockModifiedTimeStr)
				defaultS3api.On("GetObjectAttributes", &awsS3.GetObjectAttributesInput{
					Bucket: aws.String(defaultBucket),
					Key:    aws.String("config/buckets.json"),
				}).Return(&awsS3.GetObjectAttributesOutput{
					LastModified: aws.Time(mockTimeObj),
					VersionId:    nil,
				}, nil).Once()
			},
			ObjURL: "https://s3-ap-northeast-1.amazonaws.com/btq-bucket/config/buckets.json",
			ExpRes: true,
			ExpErr: nil,
		},
		{
			Case:     "has no permission access (region issue)",
			MockFunc: func() {},
			ObjURL:   "https://s3-ap-south-1.amazonaws.com/btq.ag/config/buckets.json",
			ExpRes:   false,
			ExpErr:   e.ErrWithoutPermissionToAccess,
		},
		{
			Case:     "has no permission access (bucket issue)",
			MockFunc: func() {},
			ObjURL:   "https://s3-ap-northeast-1.amazonaws.com/btq.ag/config/buckets.json",
			ExpRes:   false,
			ExpErr:   e.ErrWithoutPermissionToAccess,
		},
	}

	for idx, c := range testcase {
		mockCTX = ctx.WithValue(mockCTX, "caes no", idx)
		mockCTX = ctx.WithValue(mockCTX, "case name", c.Case)
		c.MockFunc()

		objExsists, err := s.CSP.IsObjectExists(mockCTX, c.ObjURL)
		if !s.Equal(c.ExpRes, objExsists) {
			mockCTX.Error("expected result not match")
		}

		if !s.Equal(c.ExpErr, err) {
			mockCTX.Error("expected error not match")
		}

		s.TearDownSuite()
	}
}

func (s *s3TestSuite) TestGenReadPresignedURL() {
	testcases := []struct {
		Case     string
		MockFunc func()
		ObjURL   string
		Duration time.Duration
		ExpURL   string
		ExpErr   error
	}{
		{
			Case:     "Region is incorrect",
			MockFunc: func() {},
			ObjURL:   "https://s3-ap-east-1.amazonaws.com/btq.ag/indexer/config/decoder_conf.yml",
			Duration: 10 * time.Minute,
			ExpURL:   "",
			ExpErr:   e.ErrWithoutPermissionToAccess,
		},
		{
			Case:     "Bucket is incorrect",
			MockFunc: func() {},
			ObjURL:   "https://s3-ap-northeast-1.amazonaws.com/btq.ag/indexer/config/decoder_conf.yml",
			Duration: 10 * time.Minute,
			ExpURL:   "",
			ExpErr:   e.ErrWithoutPermissionToAccess,
		},
		// {
		// 	Case: "Normal case",
		// 	MockFunc: func() {
		// 		key, _ := hex.DecodeString("31bdadd96698c204aa9ce1448ea94ae1fb4a9a0b3c9d773b51bb1822666b8f22")
		// 		keyB64 := base64.StdEncoding.EncodeToString(key)
		// 		// This is our KMS response
		// 		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 			fmt.Fprintln(w, fmt.Sprintf("%s%s%s", `{"KeyId":"test-key-id","Plaintext":"`, keyB64, `"}`))
		// 		}))
		// 		defer ts.Close()

		// 		sess := unit.Session.Copy(&aws.Config{
		// 			MaxRetries:       aws.Int(0),
		// 			Endpoint:         aws.String(ts.URL),
		// 			DisableSSL:       aws.Bool(true),
		// 			S3ForcePathStyle: aws.Bool(true),
		// 			Region:           aws.String("ap-northeast-1"),
		// 		})

		// 		attemptTime, _ := time.Parse(time.RFC3339, "2022-11-30T12:00:00-07:30")
		// 		// this is a lambda serverless yaml file, please ignore the content, it is not important
		// 		yamlStr := `
		// 		service: btq-indexer-test

		// 		package:
		// 		  individually: true
		// 		  excludeDevDependencies: false
		// 		  exclude:
		// 			- bin/**

		// 		frameworkVersion: ">=3.0.0 <4.0.0"

		// 		provider:
		// 		  name: aws
		// 		  runtime: go1.x
		// 		  architecture: x86_64
		// 		  stage: ${opt:stage, self:custom.defaultStage}
		// 		  region: ${opt:region, self:custom.defaultRegion}
		// 		  timeout: 30
		// 		  environment:
		// 			  BTQ_RUNTIME_ENV: ${opt:stage, self:custom.defaultStage}
		// 		custom:
		// 		  defaultStage: 'dev'
		// 		  defaultRegion: 'ap-northeast-1'
		// 		  prune:
		// 			automatic: true
		// 			number: 3
		// 		`
		// 		strReadPtr := strings.NewReader(yamlStr)
		// 		ioCloser := io.NopCloser(strReadPtr)
		// 		// ioSeeker := ioutil.

		// 		mockNetURLStruct := url.URL{
		// 			Scheme: "https:",
		// 			Opaque: "",
		// 			User: nil,
		// 			Host: "ap-northeast-1.amazonaws.com",
		// 			RawPath: "ap-northeast-1.amazonaws.com",
		// 			OmitHost: false,
		// 			ForceQuery: false,
		// 			RawQuery: "",
		// 			Fragment: "",
		// 			RawFragment: "",
		// 		}

		// 		mockFetchObjReqBody := strings.NewReader(``)

		// 		mockHTTPReqStruct := http.Request{
		// 			Method: http.MethodGet,
		// 			URL: &mockNetURLStruct,
		// 			Proto: "HTTP/1.1",
		// 			ProtoMajor: 1,
		// 			ProtoMinor: 0,
		// 			Header: http.Header{},
		// 			Body: ,
		// 		}

		// 		mockHTTPResponse := http.Response{
		// 			ContentLength: 566,
		// 			Body: ioCloser,
		// 			Close: true,
		// 			Uncompressed: true,
		// 			Request: &mockHTTPReqStruct,
		// 		}
		// 		defaultS3api.On("GetObjectRequest", &awsS3.GetObjectInput{
		// 			Bucket: aws.String("btq-bucket"),
		// 			Key:    aws.String("indexer/config/decoder_conf.yml"),
		// 		}).Return(, &awsS3.GetObjectOutput{
		// 			Body:             ioCloser,
		// 			BucketKeyEnabled: aws.Bool(true),
		// 			CacheControl:     aws.String("request"),
		// 		}).Once()
		// 	},
		// 	ObjURL:   "https://s3-ap-northeast-1.amazonaws.com/btq-bucket/indexer/config/decoder_conf.yml",
		// 	Duration: 10 * time.Minute,
		// 	ExpURL:   "https://s3-ap-east-1.amazonaws.com/btq.ag/indexer/config/decoder_conf.yml?X-Amz-Algorithm=AW54-HMAC-SHA256&X-Amz-Date=20221130T150956&X-Amz-SignedHeaders=host&X-Amz-Expires=900&X-Amz-Credential=AKIAYIDZSD2%202211%2Fap-northeast-1%2Fs3%2Faws4_request&aws4_request&X-amz-Signature=062ddg453248sdf3298df98u3984dsasd34sd",
		// 	ExpErr:   nil,
		// },
	}

	for idx, c := range testcases {
		mockCTX = ctx.WithValue(mockCTX, "case no", idx)
		mockCTX = ctx.WithValue(mockCTX, "case name", c.Case)
		c.MockFunc()

		readURL, err := s.CSP.GenReadPresignedURL(mockCTX, c.ObjURL, c.Duration)

		s.Equal(c.ExpURL, readURL, fmt.Sprintf("Case[%v]%v: check output read presigned url but get unexpected url", idx, c.Case))
		s.Equal(c.ExpErr, err, fmt.Sprintf("Case[%v]%v: check error get inconsist result", idx, c.Case))

		s.TearDownSuite()
	}
}

func (s *s3TestSuite) TestPutPresignedURL() {
	testcases := []struct {
		Case     string
		MockFunc func()
		ObjURL   string
		Mime     e.ContentType
		Duration time.Duration
		Metadata map[string]string
		ExpURL   string
		ExpErr   error
	}{}

	for idx, c := range testcases {
		mockCTX = ctx.WithValue(mockCTX, "case no", idx)
		mockCTX = ctx.WithValue(mockCTX, "case name", c.Case)
		c.MockFunc()

		putObjURL, err := s.CSP.PutPresignedURL(mockCTX, c.ObjURL, c.Mime, c.Duration, c.Metadata)

		s.Equal(c.ExpURL, putObjURL, fmt.Sprintf("Case[%v]%v: check output put presigned url but get unexpected url", idx, c.Case))
		s.Equal(c.ExpErr, err, fmt.Sprintf("Case[%v]%v: check error get inconsist result", idx, c.Case))

		s.TearDownSuite()
	}
}

func (s *s3TestSuite) TestUpload() {
	testcases := []struct {
		Case                string
		MockFunc            func()
		Mime                e.ContentType
		ObjURL              string
		ObjByte             []byte
		ObjMetadata         map[string]string
		ExpURL              string
		ExpReadPresignedURL string
		ExpErr              error
	}{}

	for idx, c := range testcases {
		mockCTX = ctx.WithValue(mockCTX, "case no", idx)
		mockCTX = ctx.WithValue(mockCTX, "case name", c.Case)
		c.MockFunc()

		objURL, readPresignedURL, err := s.CSP.Upload(mockCTX, c.Mime, c.ObjURL, c.ObjByte, c.ObjMetadata)

		s.Equal(c.ExpURL, objURL, fmt.Sprintf("Case[%v]%v: check output obj url but get unexpected url", idx, c.Case))
		s.Equal(c.ExpReadPresignedURL, readPresignedURL, fmt.Sprintf("Case[%v]%v: check output obj read presigned url but get unexpected url", idx, c.Case))
		s.Equal(c.ExpErr, err, fmt.Sprintf("Case[%v]%v: check error get inconsist result", idx, c.Case))

		s.TearDownSuite()
	}
}

func (s *s3TestSuite) TestOverride() {
	testcases := []struct {
		Case        string
		MockFunc    func()
		Mime        e.ContentType
		ObjURL      string
		ObjByte     []byte
		ObjMetadata map[string]string
		ExpURL      string
		ExpErr      error
	}{}

	for idx, c := range testcases {
		mockCTX = ctx.WithValue(mockCTX, "case no", idx)
		mockCTX = ctx.WithValue(mockCTX, "case name", c.Case)
		c.MockFunc()

		objURL, err := s.CSP.Override(mockCTX, c.Mime, c.ObjURL, c.ObjByte, c.ObjMetadata)

		s.Equal(c.ExpURL, objURL, fmt.Sprintf("Case[%v]%v: check output obj url but get unexpected url", idx, c.Case))

		s.Equal(c.ExpErr, err, fmt.Sprintf("Case[%v]%v: check error get inconsist result", idx, c.Case))

		s.TearDownSuite()
	}
}

func (s *s3TestSuite) TestHealth() {}

// Mockup response ..........
var (
//	mockBukcetObjetPath = []string{
//		"/config/indexer_setup.json",
//		"/config/buckets.json",
//		"/i18n/exploer/en-US.json",
//		"/i18n/exploer/zh-TW.json",
//	}
)
