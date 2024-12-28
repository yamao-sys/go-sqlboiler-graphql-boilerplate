package resolvers

import (
	"app/lib"
	models "app/models/generated"
	"app/test/factories"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type TestUserResolverSuite struct {
	WithDBSuite
}

var (
	testUserGraphQLServerHandler http.Handler
)

func (s *TestUserResolverSuite) SetupTest() {
	s.SetDBCon()

	// NOTE: テスト対象のサーバのハンドラを設定
	testUserGraphQLServerHandler = lib.GetGraphQLHttpHandler(DBCon)
}

func (s *TestUserResolverSuite) TearDownTest() {
	s.CloseDB()
}

func (s *TestUserResolverSuite) TestSignUp() {
	res := httptest.NewRecorder()
	query := map[string]interface{}{
		"query": `mutation {
            signUp(input: {
                name: "test name 1",
                email: "test@example.com",
                password: "password"
            }) {
                user {
					id,
					name,
					email,
					nameAndEmail
				},
				validationErrors {
					name,
					email,
					password
				}
            }
        }`,
	}

	signUpRequestBody, _ := json.Marshal(query)
	req := httptest.NewRequest(http.MethodPost, "/query", strings.NewReader(string(signUpRequestBody)))
	req.Header.Set("Content-Type", "application/json")
	testUserGraphQLServerHandler.ServeHTTP(res, req)

	assert.Equal(s.T(), 200, res.Code)
	responseBody := make(map[string]interface{})
	_ = json.Unmarshal(res.Body.Bytes(), &responseBody)
	assert.Contains(s.T(), responseBody["data"], "signUp")

	// NOTE: ユーザが作成されていることを確認
	isExistUser, _ := models.Users(
		qm.Where("name = ? AND email = ?", "test name 1", "test@example.com"),
	).Exists(ctx, DBCon)
	assert.True(s.T(), isExistUser)
}

func (s *TestUserResolverSuite) TestSignUp_ValidationError() {
	res := httptest.NewRecorder()
	query := map[string]interface{}{
		"query": `mutation {
            signUp(input: {
                name: "test name 1",
                email: "",
                password: "password"
            }) {
                user {
					id,
					name,
					email,
					nameAndEmail
				},
				validationErrors {
					name,
					email,
					password
				}
            }
        }`,
	}

	signUpRequestBody, _ := json.Marshal(query)
	req := httptest.NewRequest(http.MethodPost, "/query", strings.NewReader(string(signUpRequestBody)))
	req.Header.Set("Content-Type", "application/json")
	testUserGraphQLServerHandler.ServeHTTP(res, req)

	assert.Equal(s.T(), 200, res.Code)
	assert.Contains(s.T(), res.Body.String(), "\"validationErrors\":{\"name\":[],\"email\":[\"Emailは必須入力です。\"],\"password\":[]}")

	// NOTE: ユーザが作成されていないことを確認
	isExistUser, _ := models.Users(
		qm.Where("name = ? AND email = ?", "test name 1", "test@example.com"),
	).Exists(ctx, DBCon)
	assert.False(s.T(), isExistUser)
}

func (s *TestUserResolverSuite) TestSignIn() {
	// NOTE: テスト用ユーザの作成
	user := factories.UserFactory.MustCreateWithOption(map[string]interface{}{"Email": "test@example.com"}).(*models.User)
	if err := user.Insert(ctx, DBCon, boil.Infer()); err != nil {
		s.T().Fatalf("failed to create test user %v", err)
	}

	res := httptest.NewRecorder()
	query := map[string]interface{}{
		"query": `mutation {
            signIn(input: {
                email: "test@example.com",
                password: "password"
            }) {
                validationError
            }
        }`,
	}

	signInRequestBody, _ := json.Marshal(query)
	req := httptest.NewRequest(http.MethodPost, "/query", strings.NewReader(string(signInRequestBody)))
	req.Header.Set("Content-Type", "application/json")
	testUserGraphQLServerHandler.ServeHTTP(res, req)

	assert.Equal(s.T(), 200, res.Code)
	assert.NotEmpty(s.T(), res.Header().Get("Set-Cookie"))
}

func (s *TestUserResolverSuite) TestCheckSignedIn_Ok() {
	s.SetAuthUser()
	s.SignIn()

	res := httptest.NewRecorder()
	query := map[string]interface{}{
		"query": `query checkSignedIn {
			checkSignedIn {
				isSignedIn
			}
		}`,
	}

	checkSignedInRequestBody, _ := json.Marshal(query)
	req := httptest.NewRequest(http.MethodPost, "/query", strings.NewReader(string(checkSignedInRequestBody)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", "token="+token)
	testUserGraphQLServerHandler.ServeHTTP(res, req)

	assert.Equal(s.T(), 200, res.Code)
	assert.Contains(s.T(), res.Body.String(), "\"isSignedIn\":true")
}

func (s *TestUserResolverSuite) TestCheckSignedIn_NotOk() {
	res := httptest.NewRecorder()
	query := map[string]interface{}{
		"query": `query checkSignedIn {
			checkSignedIn {
				isSignedIn
			}
		}`,
	}

	checkSignedInRequestBody, _ := json.Marshal(query)
	req := httptest.NewRequest(http.MethodPost, "/query", strings.NewReader(string(checkSignedInRequestBody)))
	req.Header.Set("Content-Type", "application/json")
	testUserGraphQLServerHandler.ServeHTTP(res, req)

	assert.Equal(s.T(), 200, res.Code)
	assert.Contains(s.T(), res.Body.String(), "\"isSignedIn\":false")
}

func TestUserResolver(t *testing.T) {
	// テストスイートを実施
	suite.Run(t, new(TestUserResolverSuite))
}
