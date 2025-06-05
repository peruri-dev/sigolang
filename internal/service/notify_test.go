package service

import (
	"context"
	"testing"

	"sigolang/config/testconfig"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/suite"
)

type NotifyServiceTestSuite struct {
	suite.Suite
	svc *Services
}

func TestNotifyServiceTestSuite(t *testing.T) {
	suite.Run(t, new(NotifyServiceTestSuite))
}

func (suite *NotifyServiceTestSuite) SetupSuite() {
	_ = testconfig.ReloadTestConfig()

	suite.svc = &Services{
		Resty: resty.New(),
	}

	// block all HTTP requests
	httpmock.ActivateNonDefault(suite.svc.Resty.GetClient())
}

func (suite *NotifyServiceTestSuite) SetupTest() {
	// remove any mocks
	httpmock.Reset()
}

func (suite *NotifyServiceTestSuite) TearDownSuite() {
	httpmock.DeactivateAndReset()
}

func (suite *NotifyServiceTestSuite) TestSendNotification() {
	responder := httpmock.NewStringResponder(200, "")
	fakeUrl := "http://sigolang-example.com:8888/api/users/notify"
	httpmock.RegisterResponder("PUT", fakeUrl, responder)

	err := suite.svc.SendNotification(context.Background())

	suite.NoError(err)
}
