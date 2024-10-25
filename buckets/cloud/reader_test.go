package cloud

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	cpi "github.com/hsuanshao/go-tools/buckets/cloud/provider/interface/mocks"
	e "github.com/hsuanshao/go-tools/buckets/entity"
	ifc "github.com/hsuanshao/go-tools/buckets/interface"
	"github.com/hsuanshao/go-tools/ctx"
)

var (
	mockObjReadSrv = new(cpi.ObjectServiceProvider)
)

type readSuite struct {
	suite.Suite
	CloudReader ifc.BucketReader
}

func TestNewReader(t *testing.T) {
	mockConfig := e.Config{
		CSP:    e.AWS,
		Bucket: "btq.ag",
		Region: "ap-northeast-1",
	}

	_, err := NewReader(mockCTX, &mockConfig)
	if err != nil {
		t.Error(err)
	}

	mockConf2 := e.Config{
		CSP:    e.AWS,
		Bucket: "btq.ag",
		Region: "",
	}
	_, err = NewReader(mockCTX, &mockConf2)
	if err != e.ErrS3RegionisEmptyStr {
		t.Error(err)
	}

	mockConfig3 := e.Config{
		CSP:    e.AWS,
		Bucket: "",
		Region: "ap-northeast-1",
	}

	_, err = NewReader(mockCTX, &mockConfig3)
	if err != e.ErrS3BucketisEmptyStr {
		t.Error(err)
	}
}

func TestReaderSute(t *testing.T) {
	ctx.SetDebugLevel()
	suite.Run(t, new(readSuite))
}

func (rs *readSuite) SetupTest() {
	rs.CloudReader = &readImpl{
		cloud: mockObjReadSrv,
	}
}

func (rs *readSuite) TearDownSuite() {
	rs.CloudReader = &readImpl{
		cloud: mockObjReadSrv,
	}
	mockCTX = ctx.Background()
}

func (rs *readSuite) TestListObjURLs() {
	// bktObjLists := []string{
	// 	"indexer/config/decoder_conf.yml",
	// 	"indexer/config/sync_conf.yml",
	// 	"indexer/config/logger_conf.yml",
	// 	"indexer/config/queue_conf.yml",
	// 	"indexer/i18n/en_US.json",
	// 	"indexer/i18n/zh_TW.json",
	// 	"official_website/i18n/en_US.json",
	// 	"official_website/i18n/zh_TW.json",
	// }
	testcases := []struct {
		Case     string
		MockFunc func()
		Prefix   string
		Delim    string
		ExpRes   []string
		ExpErr   error
	}{
		{
			Case: "normal case",
			MockFunc: func() {
				mockObjReadSrv.On("GetObjectList", mockCTX, "indexer", "i18n").Return([]string{
					"https://s3-ap-northeast-1.amazonaws.com/btq.ag/indexer/config/decoder_conf.yml",
					"https://s3-ap-northeast-1.amazonaws.com/btq.ag/indexer/config/sync_conf.yml",
					"https://s3-ap-northeast-1.amazonaws.com/btq.ag/indexer/config/logger_conf.yml",
					"https://s3-ap-northeast-1.amazonaws.com/btq.ag/indexer/config/queue_conf.yml",
				}, nil).Once()
			},
			Prefix: "indexer",
			Delim:  "i18n",
			ExpRes: []string{
				"https://s3-ap-northeast-1.amazonaws.com/btq.ag/indexer/config/decoder_conf.yml",
				"https://s3-ap-northeast-1.amazonaws.com/btq.ag/indexer/config/sync_conf.yml",
				"https://s3-ap-northeast-1.amazonaws.com/btq.ag/indexer/config/logger_conf.yml",
				"https://s3-ap-northeast-1.amazonaws.com/btq.ag/indexer/config/queue_conf.yml",
			},
			ExpErr: nil,
		},
		{
			Case: "not found case",
			MockFunc: func() {
				mockObjReadSrv.On("GetObjectList", mockCTX, "pq-scale", "").Return([]string{}, nil).Once()
			},
			Prefix: "pq-scale",
			Delim:  "",
			ExpRes: []string{},
			ExpErr: nil,
		},
		{
			Case: "fetch get error return",
			MockFunc: func() {
				mockObjReadSrv.On("GetObjectList", mockCTX, "indexer", "indexer").Return(nil, e.ErrFetchObjList).Once()
			},
			Prefix: "indexer",
			Delim:  "indexer",
			ExpRes: nil,
			ExpErr: e.ErrFetchObjList,
		},
	}

	for idx, c := range testcases {
		mockCTX = ctx.WithValue(mockCTX, "case no", idx)
		mockCTX = ctx.WithValue(mockCTX, "name", fmt.Sprintf("Test ListObjURLs, Case %v", c.Case))
		c.MockFunc()

		objURLs, err := rs.CloudReader.ListObjURLs(mockCTX, c.Prefix, c.Delim)

		rs.Equal(c.ExpErr, err, fmt.Sprintf("Case[%v]%v: test return error is not as expected", idx, c.Case))
		if c.ExpErr == nil {
			rs.Equal(len(c.ExpRes), len(objURLs), fmt.Sprintf("Case[%v]%v: test return obj url list count is not as expected", idx, c.Case))
			if len(c.ExpRes) == len(objURLs) {
				for oidx, url := range objURLs {
					rs.Equal(c.ExpRes[oidx], url, fmt.Sprintf("Case[%v]%v: test return obj url is not as prespective, ExpURL idx is %v", idx, c.Case, oidx))
				}
			}
		}
	}
}

func (rs *readSuite) TestReadObjectContent() {
	testcases := []struct {
		Case        string
		MockFunc    func()
		ObjURL      string
		ExpRes      []byte
		ExpMetadata map[string]string
		ExpErr      error
	}{}

	for idx, c := range testcases {
		mockCTX = ctx.WithValue(mockCTX, "case no", idx)
		mockCTX = ctx.WithValue(mockCTX, "name", fmt.Sprintf("Test ReadObjectContent, Case %v", c.Case))
		c.MockFunc()

		objRawByte, objMetadata, err := rs.CloudReader.ReadObjectContent(mockCTX, c.ObjURL)
		rs.Equal(c.ExpRes, objRawByte, fmt.Sprintf("Case[%v]%v: test return obj raw byte is not as prespective", idx, c.Case))
		rs.Equal(c.ExpMetadata, objMetadata, fmt.Sprintf("Case[%v]%v: test return obj metadata is not as prespective", idx, c.Case))
		rs.Equal(c.ExpErr, err, fmt.Sprintf("Case[%v]%v: test return error is not as prespective", idx, c.Case))
	}
}

func (rs *readSuite) TestGenerateReadPresignedURL() {
	testcases := []struct {
		Case     string
		MockFunc func()
		Duration time.Duration
		ObjURL   string
		ExpRes   string
		ExpErr   error
	}{
		{
			Case: "normal case",
			MockFunc: func() {
				mockObjReadSrv.On("GenReadPresignedURL", mockCTX, "/client_ag_result/massa/2022/1130/aggregation_result-asdf34sSWFWsx.json", 5*time.Minute).Return("https://s3-ap-northeast-1.amazonaws.com/btq.ag/client_ag_result/massa/2022/1130/aggregation_result-asdf34sSWFWsx.json?X-Amz-Algorithm=AW54-HMAC-SHA256&X-Amz-Date=20221130T150956&X-Amz-SignedHeaders=host&X-Amz-Expires=300&X-Amz-Credential=AKIAYIDZSD2%202211%2Fap-northeast-1%2Fs3%2Faws4_request&aws4_request&X-amz-Signature=062ddg453248sdf3298df98u3984dsasd34sd", nil).Once()
			},
			Duration: 5 * time.Minute,
			ObjURL:   "/client_ag_result/massa/2022/1130/aggregation_result-asdf34sSWFWsx.json",
			ExpRes:   "https://s3-ap-northeast-1.amazonaws.com/btq.ag/client_ag_result/massa/2022/1130/aggregation_result-asdf34sSWFWsx.json?X-Amz-Algorithm=AW54-HMAC-SHA256&X-Amz-Date=20221130T150956&X-Amz-SignedHeaders=host&X-Amz-Expires=300&X-Amz-Credential=AKIAYIDZSD2%202211%2Fap-northeast-1%2Fs3%2Faws4_request&aws4_request&X-amz-Signature=062ddg453248sdf3298df98u3984dsasd34sd",
			ExpErr:   nil,
		},
		{
			Case: "object path has no object exsits",
			MockFunc: func() {
				mockObjReadSrv.On("GenReadPresignedURL", mockCTX, "/client_upload/2022/1130/tx-asdf34sSWFWsx.json", 5*time.Minute).Return("", e.ErrGenReadPresignedURL).Once()
			},
			Duration: 5 * time.Minute,
			ObjURL:   "/client_upload/2022/1130/tx-asdf34sSWFWsx.json",
			ExpRes:   "",
			ExpErr:   e.ErrGenReadPresignedURL,
		},
	}

	for idx, c := range testcases {
		mockCTX = ctx.WithValue(mockCTX, "case no", idx)
		mockCTX = ctx.WithValue(mockCTX, "name", fmt.Sprintf("Test GenerateReadPresignedURL, Case %v", c.Case))
		c.MockFunc()

		objReadPresignedURL, err := rs.CloudReader.GenerateReadPresignedURL(mockCTX, c.Duration, c.ObjURL)
		rs.Equal(c.ExpRes, objReadPresignedURL, fmt.Sprintf("Case[%v]%v: test return presigned url is not as prespective", idx, c.Case))
		rs.Equal(c.ExpErr, err, fmt.Sprintf("Case[%v]%v: test return error is not as prespective", idx, c.Case))
	}
}
