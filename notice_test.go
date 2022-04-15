package notice_test

import (
	"appupdater/common/notice"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type SuiteNotice struct {
	suite.Suite
	n *notice.Idempotent
}

func (s *SuiteNotice) SetupSuite() {
	s.n = notice.NewIdempotent(time.Second, time.Second*3)
	// assert.Equal(s.T(), "store", name)
	// assert.NoError(s.T(), s.d.Start(r))
	// s.T().Log(name)
}

// suite所有用例结束后执行
func (s *SuiteNotice) TearDownSuite() {
}

func (s *SuiteNotice) BeforeTest(suiteName, testName string) {
}

// 每个用例结束前执行，在TearDownTest之前执行
func (s *SuiteNotice) AfterTest(suiteName, testName string) {
}

func (s *SuiteNotice) TestListen001_beforeReset() {
	go func() {
		<-s.n.C
		s.T().Log("get timer 001")
	}()
	time.Sleep(500 * time.Millisecond)
	s.n.Reset()
	time.Sleep(1500 * time.Millisecond)
	s.T().Log("out")
}
func (s *SuiteNotice) TestListen002_afterReset() {
	s.n.Reset()
	s.T().Log("wait timer")
	time.Sleep(500 * time.Millisecond)
	<-s.n.C
	s.T().Log("get timer")
}
func (s *SuiteNotice) TestListen003_afterReset() {
	s.n.Reset()
	s.T().Log("wait timer")
	time.Sleep(1500 * time.Millisecond)
	<-s.n.C
	s.T().Log("get timer")
}
func (s *SuiteNotice) TestListen004_afterReset() {
	go func() {
		<-s.n.C
		s.T().Log("get timer xx")
	}()
	s.n.Reset()
	time.Sleep(500 * time.Millisecond)
	s.T().Log("touch timer")
	s.n.Reset()
	time.Sleep(500 * time.Millisecond)
	s.T().Log("touch timer")
	s.n.Reset()
	time.Sleep(2500 * time.Millisecond)
	s.T().Log("out")
}

func (s *SuiteNotice) TestListen005_dealline() {
	go func() {
		<-s.n.C
		s.T().Log("get timer xx")
	}()
	s.n.Reset()
	time.Sleep(2 * time.Second)
	s.T().Log("touch timer")
	s.n.Reset()
	time.Sleep(2 * time.Second)
	s.T().Log("touch timer")
	s.n.Reset()
	time.Sleep(2 * time.Second)
	s.T().Log("touch timer")
	s.n.Reset()
	time.Sleep(2 * time.Second)
	s.T().Log("touch timer")
	s.n.Reset()
	s.n.Reset()
	s.n.Reset()
	time.Sleep(2 * time.Second)
	s.T().Log("touch timer")
	s.n.Reset()
	s.T().Log("out")
}

func TestSuite_main(t *testing.T) {
	suite.Run(t, new(SuiteNotice))
}
