package secrettool

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/hsuanshao/go-tools/ctx"
)

var (
	mockCTX = ctx.Background()
)

type testSuite struct {
	suite.Suite
	SecretTool Utility
}

func TestSecretTool(t *testing.T) {
	suite.Run(t, new(testSuite))
}

func (ts *testSuite) SetupSuite() {
	ts.SecretTool = NewSecretTool()
}

type encryptInfo struct {
	Raw     string `json:"raw"`
	Encrypt string `json:"encrypt"`
	Key     string `json:"key"`
}

func (ts *testSuite) TestEncrypt() {
	tCases := []struct {
		CTX        ctx.CTX
		RawMessage string
		ExpErr     error
	}{
		{
			CTX:        mockCTX,
			RawMessage: "hello world, BTQ",
			ExpErr:     nil,
		},
		{
			CTX:        mockCTX,
			RawMessage: "        ",
			ExpErr:     ErrEmptyMessageForEncrypt,
		},
		{
			CTX:        mockCTX,
			RawMessage: "dbindxeradmin",
			ExpErr:     nil,
		},
		{
			CTX:        mockCTX,
			RawMessage: "kya_WEZ3aet_wrj8nhn",
			ExpErr:     nil,
		},
	}

	for ci, tc := range tCases {
		encryptedMsg, publicKey, err := ts.SecretTool.Encrypt(tc.CTX, tc.RawMessage)
		if err == nil {
			f, _ := os.Create("./testdata/case_" + fmt.Sprintf("%d", ci) + ".json")
			s := encryptInfo{
				Raw:     tc.RawMessage,
				Encrypt: encryptedMsg,
				Key:     publicKey,
			}
			sb, _ := json.Marshal(s)
			f.Write(sb)
			f.Close()
		}

		if encryptedMsg == "" && tc.ExpErr != nil {
			ts.Error(fmt.Errorf("encrypt message should not as empty string"))
		}

		if publicKey == "" && tc.ExpErr != nil {
			ts.Error(fmt.Errorf("public key should not as empty public key return"))
		}

		ts.Equal(tc.ExpErr, err, fmt.Sprintf("Case[%d] compare error", ci))
	}
}

func (ts *testSuite) TestDecrypt() {
	var testErr error

	fs, _ := os.ReadDir("./testdata")

	for ci, fe := range fs {
		s, _ := fe.Info()
		r, _ := os.ReadFile("./testdata/" + s.Name())

		dst := encryptInfo{}
		err := json.Unmarshal(r, &dst)
		if err != nil {
			mockCTX.WithField("err", err).Error("json unmarshal hmm")
			continue
		}

		rawMsg, err := ts.SecretTool.Decrypt(mockCTX, dst.Encrypt, dst.Key)
		ts.Equal(dst.Raw, rawMsg, fmt.Sprintf("Case[%d] compare encrypted message", ci))
		ts.Equal(testErr, err, fmt.Sprintf("Case[%d] compare error", ci))

		testErr = nil
	}
}
