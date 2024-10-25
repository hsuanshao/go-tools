package randm

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestSuite(t *testing.T) {
	suite.Run(t, new(randSuite))
}

type randSuite struct {
	suite.Suite
	node int64
	m    Method
}

func (s *randSuite) SetupSuite() {
	s.node = 1
	s.m = NewDice(s.node)
}

func (s *randSuite) TearDownSuite() {
	s.node = 1
	s.m = NewDice(s.node)
}

func (s *randSuite) TestGenerateRID() {
	//	rand.Seed(int64(10))
	testcase := []struct {
		Case         string
		MockFunc     func()
		ExecuteTimes uint
	}{
		{
			Case:         "Node 1",
			ExecuteTimes: 1000,
		},
	}

	//var td time.Duration
	for idx, c := range testcase {
		rids := make([]RID, c.ExecuteTimes)

		for i := 0; i < int(c.ExecuteTimes); i++ {
			// st := rand.Intn(9)

			// td = time.Duration(st) * time.Millisecond
			// time.Sleep(td)

			rids[i] = s.m.GenerateRID()
		}

		// condition 1 check
		for i := 0; i < int(c.ExecuteTimes)-1; i++ {
			fmt.Println(rids[i])
			if !(rids[i] < rids[i+1] && rids[i] != rids[i+1]) {
				s.Error(fmt.Errorf("case %d: i RID %v, i+1 RID %v, doesn't match RID definition", idx, rids[i], rids[i+1]))
			}
		}
	}
}

func (s *randSuite) TestGenRandomString() {
	testcase := []struct {
		Case      string
		StringLen uint
	}{
		{
			Case:      "Node 1, string len 25",
			StringLen: 25,
		},
	}

	for idx, c := range testcase {
		randStr := s.m.GenRandomString(c.StringLen)
		s.Equal(int(c.StringLen), len(randStr), fmt.Sprintf("Case %d: generate rand string length not match", idx))
	}
}

func (s *randSuite) TestIsValidateRID() {
	tCases := []struct {
		Case     string
		InputRID int64
		ExpRes   bool
	}{
		{
			Case:     "normal case",
			InputRID: 312119580761789396,
			ExpRes:   true,
		},
		{
			Case:     "future (4000/01/01 00:00:00) generate rid",
			InputRID: 3624795013042410472,
			ExpRes:   false,
		},
	}

	for i, c := range tCases {
		res := s.m.IsValidateRID(c.InputRID)
		s.Equal(c.ExpRes, res, fmt.Sprintf("Case[%d] Is validate RID result is not as expected", i))
	}
}
