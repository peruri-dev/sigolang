package service

import (
	"context"
	"testing"

	"sigolang/config/testconfig"
	"sigolang/db/seeds"
	"sigolang/lib/db"

	"github.com/stretchr/testify/suite"
)

type UserServiceTestSuite struct {
	suite.Suite
	svc *Services
}

func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}

func (suite *UserServiceTestSuite) SetupSuite() {
	c := testconfig.ReloadTestConfig()

	dbConn, err := db.Open(c)
	if err != nil {
		suite.T().Skipf("please enable db_sqlite: %s", err.Error())
	}

	suite.svc = &Services{
		DB: dbConn,
	}
	seeds.ResetSchema(context.Background(), dbConn)
}

func (suite *UserServiceTestSuite) TestListUsers() {
	users, err := suite.svc.UserList(context.Background())
	suite.NoError(err)

	suite.GreaterOrEqual(len(users), 2)
}

func (suite *UserServiceTestSuite) TestCreateUser() {
	user, err := suite.svc.UserCreate(context.Background(), "Bob", []string{"b@bbb.com"})
	suite.NoError(err)

	suite.Equal("Bob", user.Name)
}
