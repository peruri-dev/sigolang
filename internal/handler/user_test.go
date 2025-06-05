package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"sigolang/internal/model"
	"sigolang/mocks"

	"github.com/danielgtaylor/huma/v2/humatest"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UserHandlerTestSuite struct {
	suite.Suite
	api     humatest.TestAPI
	mockSvc *mocks.AllServices
}

var MockContext = mock.MatchedBy(func(c context.Context) bool { return true })

func TestUserHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerTestSuite))
}

func (suite *UserHandlerTestSuite) SetupSuite() {
	suite.mockSvc = mocks.NewAllServices(suite.T())

	_, api := humatest.New(suite.T())
	suite.api = api

	// Register routes...
	h := &Handler{
		suite.mockSvc,
	}
	h.RegisterUser(api)
}

func (suite *UserHandlerTestSuite) TestListUsersReturnEmpty() {
	suite.mockSvc.On("UserList", MockContext).Return([]model.User{}, nil)
	// Make a GET request
	resp := suite.api.Get("/api/users")
	suite.Equalf(http.StatusNoContent, resp.Code, "unexpected status code %d", resp.Code)
}

func (suite *UserHandlerTestSuite) TestCreateUser() {
	suite.mockSvc.On("UserCreate", MockContext, "Alice", []string{"al@ice.com"}).Return(
		&model.User{
			Name: "Alex",
		}, nil,
	)
	suite.mockSvc.On("SendNotification", MockContext).Return(nil)

	// Make a PUT request
	resp := suite.api.Post("/api/users",
		"X-Authorization: abc123",
		map[string]any{
			"name": "Alice",
			"emails": []string{
				"al@ice.com",
			},
		})

	suite.Equalf(http.StatusCreated, resp.Code, "unexpected status code %d", resp.Code)

	// Convert the JSON response to a map
	var response map[string]any

	err := json.Unmarshal(resp.Body.Bytes(), &response)
	suite.NoError(err)

	suite.Equal("Alex", response["name"])
}
