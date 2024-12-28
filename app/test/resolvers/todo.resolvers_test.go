package resolvers

import (
	"app/lib"
	models "app/models/generated"
	"app/test/factories"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type TestTodoResolverSuite struct {
	WithDBSuite
}

var (
	testTodoGraphQLServerHandler http.Handler
)

func (s *TestTodoResolverSuite) SetupTest() {
	s.SetDBCon()

	// NOTE: テスト対象のサーバのハンドラを設定
	testTodoGraphQLServerHandler = lib.GetGraphQLHttpHandler(DBCon)
}

func (s *TestTodoResolverSuite) TearDownTest() {
	s.CloseDB()
}

func (s *TestTodoResolverSuite) TestCreateTodo_Unauthorized() {
	res := httptest.NewRecorder()
	query := map[string]interface{}{
		"query": `mutation {
            createTodo(input: {
                title: "test title 1",
                content: "",
            }) {
                validationErrors {
					title
					content
				}
            }
        }`,
	}

	signUpRequestBody, _ := json.Marshal(query)
	req := httptest.NewRequest(http.MethodPost, "/query", strings.NewReader(string(signUpRequestBody)))
	req.Header.Set("Content-Type", "application/json")
	testTodoGraphQLServerHandler.ServeHTTP(res, req)

	assert.Equal(s.T(), 200, res.Code)
	responseBody := make(map[string]([1]map[string]map[string]interface{}))
	_ = json.Unmarshal(res.Body.Bytes(), &responseBody)
	assert.Equal(s.T(), float64(401), responseBody["errors"][0]["extensions"]["code"])

	// NOTE: Todoリストが作成されていないことを確認
	isExistTodo, _ := models.Todos(
		qm.Where("title = ?", "test title 1"),
	).Exists(ctx, DBCon)
	assert.False(s.T(), isExistTodo)
}

func (s *TestTodoResolverSuite) TestCreateTodo() {
	s.SetAuthUser()
	s.SignIn()

	res := httptest.NewRecorder()
	query := map[string]interface{}{
		"query": `mutation {
            createTodo(input: {
                title: "test title 1",
                content: "",
            }) {
                validationErrors {
					title
					content
				}
            }
        }`,
	}

	signUpRequestBody, _ := json.Marshal(query)
	req := httptest.NewRequest(http.MethodPost, "/query", strings.NewReader(string(signUpRequestBody)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", "token="+token)
	testTodoGraphQLServerHandler.ServeHTTP(res, req)

	assert.Equal(s.T(), 200, res.Code)
	assert.Contains(s.T(), res.Body.String(), "\"validationErrors\":{\"title\":[],\"content\":[]}")

	// NOTE: Todoリストが作成されていることを確認
	isExistTodo, _ := models.Todos(
		qm.Where("title = ?", "test title 1"),
	).Exists(ctx, DBCon)
	assert.True(s.T(), isExistTodo)
}

func (s *TestTodoResolverSuite) TestCreateTodo_ValidationError() {
	s.SetAuthUser()
	s.SignIn()

	res := httptest.NewRecorder()
	query := map[string]interface{}{
		"query": `mutation {
            createTodo(input: {
                title: "",
                content: "",
            }) {
                validationErrors {
					title
					content
				}
            }
        }`,
	}

	signUpRequestBody, _ := json.Marshal(query)
	req := httptest.NewRequest(http.MethodPost, "/query", strings.NewReader(string(signUpRequestBody)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", "token="+token)
	testTodoGraphQLServerHandler.ServeHTTP(res, req)

	assert.Equal(s.T(), 200, res.Code)
	assert.Contains(s.T(), res.Body.String(), "\"validationErrors\":{\"title\":[\"タイトルは必須入力です。\"],\"content\":[]}")

	// NOTE: Todoリストが作成されていないことを確認
	isExistTodo, _ := models.Todos(
		qm.Where("title = ?", "test title 1"),
	).Exists(ctx, DBCon)
	assert.False(s.T(), isExistTodo)
}

func (s *TestTodoResolverSuite) TestFetchTodo_Unauthorized() {
	s.SetAuthUser()
	testTodo := models.Todo{Title: "test title 1", Content: null.String{String: "test content 1", Valid: true}, UserID: user.ID}
	if err := testTodo.Insert(ctx, DBCon, boil.Infer()); err != nil {
		s.T().Fatalf("failed to create test todos %v", err)
	}

	res := httptest.NewRecorder()
	id := strconv.Itoa(testTodo.ID)
	query := map[string]interface{}{
		"query": `query {
            fetchTodo(id: ` + id + `) {
                id,
                title,
                content,
                createdAt,
				updatedAt
            }
        }`,
	}

	signUpRequestBody, _ := json.Marshal(query)
	req := httptest.NewRequest(http.MethodPost, "/query", strings.NewReader(string(signUpRequestBody)))
	req.Header.Set("Content-Type", "application/json")
	testTodoGraphQLServerHandler.ServeHTTP(res, req)

	assert.Equal(s.T(), 200, res.Code)
	responseBody := make(map[string]([1]map[string]map[string]interface{}))
	_ = json.Unmarshal(res.Body.Bytes(), &responseBody)
	assert.Equal(s.T(), float64(401), responseBody["errors"][0]["extensions"]["code"])
}

func (s *TestTodoResolverSuite) TestFetchTodo() {
	s.SetAuthUser()
	s.SignIn()

	testTodo := models.Todo{Title: "test title 1", Content: null.String{String: "test content 1", Valid: true}, UserID: user.ID}
	if err := testTodo.Insert(ctx, DBCon, boil.Infer()); err != nil {
		s.T().Fatalf("failed to create test todos %v", err)
	}

	res := httptest.NewRecorder()
	id := strconv.Itoa(testTodo.ID)
	query := map[string]interface{}{
		"query": `query {
            fetchTodo(id: ` + id + `) {
                id,
                title,
                content,
                createdAt,
				updatedAt,
				user {
					id
				}
            }
        }`,
	}

	signUpRequestBody, _ := json.Marshal(query)
	req := httptest.NewRequest(http.MethodPost, "/query", strings.NewReader(string(signUpRequestBody)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", "token="+token)
	testTodoGraphQLServerHandler.ServeHTTP(res, req)

	assert.Equal(s.T(), 200, res.Code)
	responseBody := make(map[string](map[string]map[string]interface{}))
	_ = json.Unmarshal(res.Body.Bytes(), &responseBody)
	assert.Contains(s.T(), responseBody["data"], "fetchTodo")
	assert.Equal(s.T(), float64(testTodo.ID), responseBody["data"]["fetchTodo"]["id"])
	assert.Contains(s.T(), responseBody["data"]["fetchTodo"]["user"], "id")
}

func (s *TestTodoResolverSuite) TestFetchTodo_NotFound() {
	s.SetAuthUser()
	s.SignIn()
	testTodo := models.Todo{Title: "test title 1", Content: null.String{String: "test content 1", Valid: true}, UserID: user.ID}
	if err := testTodo.Insert(ctx, DBCon, boil.Infer()); err != nil {
		s.T().Fatalf("failed to create test todos %v", err)
	}

	res := httptest.NewRecorder()
	id := strconv.Itoa(testTodo.ID + 1)
	query := map[string]interface{}{
		"query": `query {
            fetchTodo(id: ` + id + `) {
                id,
                title,
                content,
                createdAt,
				updatedAt,
				user {
					id
				}
            }
        }`,
	}

	signUpRequestBody, _ := json.Marshal(query)
	req := httptest.NewRequest(http.MethodPost, "/query", strings.NewReader(string(signUpRequestBody)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", "token="+token)
	testTodoGraphQLServerHandler.ServeHTTP(res, req)

	assert.Equal(s.T(), 200, res.Code)
	responseBody := make(map[string]([1]map[string]map[string]interface{}))
	_ = json.Unmarshal(res.Body.Bytes(), &responseBody)
	assert.Equal(s.T(), float64(404), responseBody["errors"][0]["extensions"]["code"])
}

func (s *TestTodoResolverSuite) TestFetchTodoLists_Unauthorized() {
	s.SetAuthUser()
	testTodo := models.Todo{Title: "test title 1", Content: null.String{String: "test content 1", Valid: true}, UserID: user.ID}
	if err := testTodo.Insert(ctx, DBCon, boil.Infer()); err != nil {
		s.T().Fatalf("failed to create test todos %v", err)
	}

	res := httptest.NewRecorder()
	query := map[string]interface{}{
		"query": `query {
            fetchTodoLists {
                id,
                title,
                content,
                createdAt,
				updatedAt,
				user {
					id
				}
            }
        }`,
	}

	signUpRequestBody, _ := json.Marshal(query)
	req := httptest.NewRequest(http.MethodPost, "/query", strings.NewReader(string(signUpRequestBody)))
	req.Header.Set("Content-Type", "application/json")
	testTodoGraphQLServerHandler.ServeHTTP(res, req)

	assert.Equal(s.T(), 200, res.Code)
	responseBody := make(map[string]([1]map[string]map[string]interface{}))
	_ = json.Unmarshal(res.Body.Bytes(), &responseBody)
	assert.Equal(s.T(), float64(401), responseBody["errors"][0]["extensions"]["code"])
}

func (s *TestTodoResolverSuite) TestFetchTodoLists() {
	s.SetAuthUser()
	s.SignIn()

	var todosSlice models.TodoSlice
	todosSlice = append(todosSlice, &models.Todo{
		Title:   "test title 1",
		Content: null.String{String: "test content 1", Valid: true},
		UserID:  user.ID,
	})
	todosSlice = append(todosSlice, &models.Todo{
		Title:   "test title 2",
		Content: null.String{String: "test content 2", Valid: true},
		UserID:  user.ID,
	})
	_, err := todosSlice.InsertAll(ctx, DBCon, boil.Infer())
	if err != nil {
		s.T().Fatalf("failed to create TestFetchTodoLists Data: %v", err)
	}

	res := httptest.NewRecorder()
	query := map[string]interface{}{
		"query": `query {
            fetchTodoLists {
                id,
                title,
                content,
                createdAt,
				updatedAt,
				user {
					id
				}
            }
        }`,
	}

	signUpRequestBody, _ := json.Marshal(query)
	req := httptest.NewRequest(http.MethodPost, "/query", strings.NewReader(string(signUpRequestBody)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", "token="+token)
	testTodoGraphQLServerHandler.ServeHTTP(res, req)

	assert.Equal(s.T(), 200, res.Code)
	responseBody := make(map[string](map[string]([]map[string]interface{})))
	_ = json.Unmarshal(res.Body.Bytes(), &responseBody)
	assert.Len(s.T(), responseBody["data"]["fetchTodoLists"], 2)
	assert.Contains(s.T(), responseBody["data"]["fetchTodoLists"][0]["user"], "id")
}

func (s *TestTodoResolverSuite) TestFetchTodoLists_CountZero() {
	s.SetAuthUser()
	s.SignIn()

	res := httptest.NewRecorder()
	query := map[string]interface{}{
		"query": `query {
            fetchTodoLists {
                id,
                title,
                content,
                createdAt,
				updatedAt,
				user {
					id
				}
            }
        }`,
	}

	signUpRequestBody, _ := json.Marshal(query)
	req := httptest.NewRequest(http.MethodPost, "/query", strings.NewReader(string(signUpRequestBody)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", "token="+token)
	testTodoGraphQLServerHandler.ServeHTTP(res, req)

	assert.Equal(s.T(), 200, res.Code)
	responseBody := make(map[string](map[string]([]map[string]interface{})))
	_ = json.Unmarshal(res.Body.Bytes(), &responseBody)
	assert.Len(s.T(), responseBody["data"]["fetchTodoLists"], 0)
}

func (s *TestTodoResolverSuite) TestUpdateTodo_Unauthorized() {
	s.SetAuthUser()
	testTodo := models.Todo{Title: "test title 1", Content: null.String{String: "test content 1", Valid: true}, UserID: user.ID}
	if err := testTodo.Insert(ctx, DBCon, boil.Infer()); err != nil {
		s.T().Fatalf("failed to create test todos %v", err)
	}

	res := httptest.NewRecorder()
	id := strconv.Itoa(testTodo.ID)
	query := map[string]interface{}{
		"query": `mutation {
            updateTodo(id: ` + id + `, input: {
                title: "test updated title 1",
                content: "test updated content 1",
            }) {
                id
                validationErrors {
					title
					content
				}
            }
        }`,
	}

	signUpRequestBody, _ := json.Marshal(query)
	req := httptest.NewRequest(http.MethodPost, "/query", strings.NewReader(string(signUpRequestBody)))
	req.Header.Set("Content-Type", "application/json")
	testTodoGraphQLServerHandler.ServeHTTP(res, req)

	assert.Equal(s.T(), 200, res.Code)
	responseBody := make(map[string]([1]map[string]map[string]interface{}))
	_ = json.Unmarshal(res.Body.Bytes(), &responseBody)
	assert.Equal(s.T(), float64(401), responseBody["errors"][0]["extensions"]["code"])
}

func (s *TestTodoResolverSuite) TestUpdateTodo() {
	s.SetAuthUser()
	s.SignIn()

	testTodo := models.Todo{Title: "test title 1", Content: null.String{String: "test content 1", Valid: true}, UserID: user.ID}
	if err := testTodo.Insert(ctx, DBCon, boil.Infer()); err != nil {
		s.T().Fatalf("failed to create test todos %v", err)
	}

	res := httptest.NewRecorder()
	id := strconv.Itoa(testTodo.ID)
	query := map[string]interface{}{
		"query": `mutation {
            updateTodo(id: ` + id + `, input: {
                title: "test updated title 1",
                content: "test updated content 1",
            }) {
                id
                validationErrors {
					title
					content
				}
            }
        }`,
	}

	signUpRequestBody, _ := json.Marshal(query)
	req := httptest.NewRequest(http.MethodPost, "/query", strings.NewReader(string(signUpRequestBody)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", "token="+token)
	testTodoGraphQLServerHandler.ServeHTTP(res, req)

	assert.Equal(s.T(), 200, res.Code)
	assert.Contains(s.T(), res.Body.String(), "\"validationErrors\":{\"title\":[],\"content\":[]}")

	// NOTE: TODOが更新されていることの確認
	if err := testTodo.Reload(ctx, DBCon); err != nil {
		s.T().Fatalf("failed to reload test todos %v", err)
	}
	assert.Equal(s.T(), "test updated title 1", testTodo.Title)
	assert.Equal(s.T(), null.String{String: "test updated content 1", Valid: true}, testTodo.Content)
}

func (s *TestTodoResolverSuite) TestUpdateTodo_NotFound() {
	s.SetAuthUser()
	s.SignIn()

	otherUser := factories.UserFactory.MustCreateWithOption(map[string]interface{}{"Email": "test_2@example.com"}).(*models.User)
	if err := otherUser.Insert(ctx, DBCon, boil.Infer()); err != nil {
		s.T().Fatalf("failed to create test user %v", err)
	}
	testTodo := models.Todo{Title: "test title 1", Content: null.String{String: "test content 1", Valid: true}, UserID: otherUser.ID}
	if err := testTodo.Insert(ctx, DBCon, boil.Infer()); err != nil {
		s.T().Fatalf("failed to create test todos %v", err)
	}

	res := httptest.NewRecorder()
	id := strconv.Itoa(testTodo.ID)
	query := map[string]interface{}{
		"query": `mutation {
            updateTodo(id: ` + id + `, input: {
                title: "test updated title 1",
                content: "test updated content 1",
            }) {
                id
                validationErrors {
					title
					content
				}
            }
        }`,
	}

	signUpRequestBody, _ := json.Marshal(query)
	req := httptest.NewRequest(http.MethodPost, "/query", strings.NewReader(string(signUpRequestBody)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", "token="+token)
	testTodoGraphQLServerHandler.ServeHTTP(res, req)

	assert.Equal(s.T(), 200, res.Code)
	responseBody := make(map[string]([1]map[string]map[string]interface{}))
	_ = json.Unmarshal(res.Body.Bytes(), &responseBody)
	assert.Equal(s.T(), float64(404), responseBody["errors"][0]["extensions"]["code"])
}

func (s *TestTodoResolverSuite) TestUpdateTodo_ValidationError() {
	s.SetAuthUser()
	s.SignIn()

	testTodo := models.Todo{Title: "test title 1", Content: null.String{String: "test content 1", Valid: true}, UserID: user.ID}
	if err := testTodo.Insert(ctx, DBCon, boil.Infer()); err != nil {
		s.T().Fatalf("failed to create test todos %v", err)
	}

	res := httptest.NewRecorder()
	id := strconv.Itoa(testTodo.ID)
	query := map[string]interface{}{
		"query": `mutation {
            updateTodo(id: ` + id + `, input: {
                title: "",
                content: "test updated content 1",
            }) {
                id
                validationErrors {
					title
					content
				}
            }
        }`,
	}

	signUpRequestBody, _ := json.Marshal(query)
	req := httptest.NewRequest(http.MethodPost, "/query", strings.NewReader(string(signUpRequestBody)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", "token="+token)
	testTodoGraphQLServerHandler.ServeHTTP(res, req)

	assert.Equal(s.T(), 200, res.Code)
	assert.Contains(s.T(), res.Body.String(), "\"validationErrors\":{\"title\":[\"タイトルは必須入力です。\"],\"content\":[]}")
}

func (s *TestTodoResolverSuite) TestDeleteTodo_Unauthorized() {
	s.SetAuthUser()
	testTodo := models.Todo{Title: "test title 1", Content: null.String{String: "test content 1", Valid: true}, UserID: user.ID}
	if err := testTodo.Insert(ctx, DBCon, boil.Infer()); err != nil {
		s.T().Fatalf("failed to create test todos %v", err)
	}

	res := httptest.NewRecorder()
	id := strconv.Itoa(testTodo.ID)
	query := map[string]interface{}{
		"query": `mutation {
            deleteTodo(id: ` + id + `)
        }`,
	}

	signUpRequestBody, _ := json.Marshal(query)
	req := httptest.NewRequest(http.MethodPost, "/query", strings.NewReader(string(signUpRequestBody)))
	req.Header.Set("Content-Type", "application/json")
	testTodoGraphQLServerHandler.ServeHTTP(res, req)

	assert.Equal(s.T(), 200, res.Code)
	responseBody := make(map[string]([1]map[string]map[string]interface{}))
	_ = json.Unmarshal(res.Body.Bytes(), &responseBody)
	assert.Equal(s.T(), float64(401), responseBody["errors"][0]["extensions"]["code"])
}

func (s *TestTodoResolverSuite) TestDeleteTodo() {
	s.SetAuthUser()
	s.SignIn()

	testTodo := models.Todo{Title: "test title 1", Content: null.String{String: "test content 1", Valid: true}, UserID: user.ID}
	if err := testTodo.Insert(ctx, DBCon, boil.Infer()); err != nil {
		s.T().Fatalf("failed to create test todos %v", err)
	}

	res := httptest.NewRecorder()
	id := strconv.Itoa(testTodo.ID)
	query := map[string]interface{}{
		"query": `mutation {
            deleteTodo(id: ` + id + `)
        }`,
	}

	signUpRequestBody, _ := json.Marshal(query)
	req := httptest.NewRequest(http.MethodPost, "/query", strings.NewReader(string(signUpRequestBody)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", "token="+token)
	testTodoGraphQLServerHandler.ServeHTTP(res, req)

	assert.Equal(s.T(), 200, res.Code)
	responseBody := make(map[string]interface{})
	_ = json.Unmarshal(res.Body.Bytes(), &responseBody)
	assert.Contains(s.T(), responseBody["data"], "deleteTodo")

	// NOTE: TODOが削除されていることの確認
	reloadErr := testTodo.Reload(ctx, DBCon)
	assert.NotNil(s.T(), reloadErr)
}

func (s *TestTodoResolverSuite) TestDeleteTodo_NotFound() {
	s.SetAuthUser()
	s.SignIn()

	otherUser := factories.UserFactory.MustCreateWithOption(map[string]interface{}{"Email": "test_2@example.com"}).(*models.User)
	if err := otherUser.Insert(ctx, DBCon, boil.Infer()); err != nil {
		s.T().Fatalf("failed to create test user %v", err)
	}
	testTodo := models.Todo{Title: "test title 1", Content: null.String{String: "test content 1", Valid: true}, UserID: otherUser.ID}
	if err := testTodo.Insert(ctx, DBCon, boil.Infer()); err != nil {
		s.T().Fatalf("failed to create test todos %v", err)
	}

	res := httptest.NewRecorder()
	id := strconv.Itoa(testTodo.ID)
	query := map[string]interface{}{
		"query": `mutation {
            deleteTodo(id: ` + id + `)
        }`,
	}

	signUpRequestBody, _ := json.Marshal(query)
	req := httptest.NewRequest(http.MethodPost, "/query", strings.NewReader(string(signUpRequestBody)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", "token="+token)
	testTodoGraphQLServerHandler.ServeHTTP(res, req)

	assert.Equal(s.T(), 200, res.Code)
	responseBody := make(map[string]([1]map[string]map[string]interface{}))
	_ = json.Unmarshal(res.Body.Bytes(), &responseBody)
	assert.Equal(s.T(), float64(404), responseBody["errors"][0]["extensions"]["code"])
}

func TestTodoResolver(t *testing.T) {
	// テストスイートを実施
	suite.Run(t, new(TestTodoResolverSuite))
}
