package services

import (
	"app/graph/model"
	models "app/models/generated"
	"app/test/factories"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type TestAuthServiceSuite struct {
	WithDBSuite
}

var testAuthService AuthService

func (s *TestAuthServiceSuite) SetupTest() {
	s.SetDBCon()

	testAuthService = NewAuthService(DBCon)
}

func (s *TestAuthServiceSuite) TearDownTest() {
	s.CloseDB()
}

func (s *TestAuthServiceSuite) TestSignUp() {
	requestParams := model.SignUpInput{Name: "test name 1", Email: "test@example.com", Password: "password"}

	_, err := testAuthService.SignUp(ctx, requestParams)

	assert.Nil(s.T(), err)

	// NOTE: ユーザが作成されていることを確認
	isExistUser, err := models.Users(
		qm.Where("name = ? AND email = ?", "test name 1", "test@example.com"),
	).Exists(ctx, DBCon)
	if err != nil {
		s.T().Fatalf("failed to create user %v", err)
	}
	assert.True(s.T(), isExistUser)
}

func (s *TestAuthServiceSuite) TestSignUp_ValidationError() {
	requestParams := model.SignUpInput{Name: "test name 1", Email: "", Password: "password"}

	res, _ := testAuthService.SignUp(ctx, requestParams)

	assert.Len(s.T(), res.ValidationErrors.Name, 0)
	assert.Len(s.T(), res.ValidationErrors.Email, 1)
	assert.Len(s.T(), res.ValidationErrors.Password, 0)

	// NOTE: ユーザが作成されていないことを確認
	isExistUser, _ := models.Users(
		qm.Where("name = ?", "test name 1"),
	).Exists(ctx, DBCon)
	assert.False(s.T(), isExistUser)
}

func (s *TestAuthServiceSuite) TestSignIn() {
	// NOTE: テスト用ユーザの作成
	user := factories.UserFactory.MustCreateWithOption(map[string]interface{}{"Email": "test@example.com"}).(*models.User)
	if err := user.Insert(ctx, DBCon, boil.Infer()); err != nil {
		s.T().Fatalf("failed to create test user %v", err)
	}

	requestParams := model.SignInInput{Email: "test@example.com", Password: "password"}

	token, validationError, err := testAuthService.SignIn(ctx, requestParams)

	assert.NotNil(s.T(), token)
	assert.Nil(s.T(), validationError)
	assert.Nil(s.T(), err)
}

func (s *TestAuthServiceSuite) TestSignIn_NotFoundError() {
	// NOTE: テスト用ユーザの作成
	user := factories.UserFactory.MustCreateWithOption(map[string]interface{}{"Email": "test@example.com"}).(*models.User)
	if err := user.Insert(ctx, DBCon, boil.Infer()); err != nil {
		s.T().Fatalf("failed to create test user %v", err)
	}

	requestParams := model.SignInInput{Email: "test_1@example.com", Password: "password"}

	token, validationError, err := testAuthService.SignIn(ctx, requestParams)

	assert.Empty(s.T(), token)
	assert.NotNil(s.T(), validationError)
	assert.Nil(s.T(), err)
}

func TestAuthService(t *testing.T) {
	// テストスイートを実行
	suite.Run(t, new(TestAuthServiceSuite))
}
