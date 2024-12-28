package services

import (
	"app/graph/model"
	models "app/models/generated"
	"app/test/factories"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type TestTodoServiceSuite struct {
	WithDBSuite
}

var (
	user            *models.User
	testTodoService TodoService
)

func (s *TestTodoServiceSuite) SetupTest() {
	s.SetDBCon()

	// NOTE: テスト用ユーザの作成
	user = factories.UserFactory.MustCreateWithOption(map[string]interface{}{"Email": "test@example.com"}).(*models.User)
	if err := user.Insert(ctx, DBCon, boil.Infer()); err != nil {
		s.T().Fatalf("failed to create user %v", err)
	}

	testTodoService = NewTodoService(DBCon)
}

func (s *TestTodoServiceSuite) TearDownTest() {
	s.CloseDB()
}

func (s *TestTodoServiceSuite) TestCreateTodo() {
	requestParams := model.CreateTodoInput{Title: "test title 1", Content: "test content 1"}

	res, err := testTodoService.CreateTodo(ctx, requestParams, user.ID)

	var expectedValidationErrors []string
	assert.NotEmpty(s.T(), res.ID)
	assert.Equal(s.T(), expectedValidationErrors, res.ValidationErrors.Title)
	assert.Equal(s.T(), expectedValidationErrors, res.ValidationErrors.Content)
	assert.Nil(s.T(), err)

	// NOTE: Todoリストが作成されていることを確認
	isExistTodo, _ := models.Todos(
		qm.Where("title = ?", "test title 1"),
	).Exists(ctx, DBCon)
	assert.True(s.T(), isExistTodo)
}

func (s *TestTodoServiceSuite) TestCreateTodo_ValidationError() {
	requestParams := model.CreateTodoInput{Title: "", Content: "test content 1"}

	res, err := testTodoService.CreateTodo(ctx, requestParams, user.ID)

	assert.Nil(s.T(), err)
	assert.Empty(s.T(), res.ID)
	expectedTitleValidationErrors := []string{"タイトルは必須入力です。"}
	var expectedContentValidationErrors []string
	assert.Equal(s.T(), expectedTitleValidationErrors, res.ValidationErrors.Title)
	assert.Equal(s.T(), expectedContentValidationErrors, res.ValidationErrors.Content)

	// NOTE: Todoリストが作成されていないことを確認
	isExistTodo, _ := models.Users(
		qm.Where("title = ?", "test title 1"),
	).Exists(ctx, DBCon)
	assert.False(s.T(), isExistTodo)
}

func (s *TestTodoServiceSuite) TestFetchTodoLists() {
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

	todos, err := testTodoService.FetchTodoLists(ctx, user.ID)

	assert.Nil(s.T(), err)
	assert.Len(s.T(), todos, 2)
}

func (s *TestTodoServiceSuite) TestFetchTodo() {
	testTodo := models.Todo{Title: "test title 1", Content: null.String{String: "test content 1", Valid: true}, UserID: user.ID}
	if err := testTodo.Insert(ctx, DBCon, boil.Infer()); err != nil {
		s.T().Fatalf("failed to create test todos %v", err)
	}
	testTodo.Reload(ctx, DBCon)

	todo, err := testTodoService.FetchTodo(ctx, testTodo.ID, user.ID)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), testTodo.Title, todo.Title)
}

func (s *TestTodoServiceSuite) TestUpdateTodo() {
	testTodo := models.Todo{Title: "test title 1", Content: null.String{String: "test content 1", Valid: true}, UserID: user.ID}
	if err := testTodo.Insert(ctx, DBCon, boil.Infer()); err != nil {
		s.T().Fatalf("failed to create test todos %v", err)
	}

	requestParams := model.UpdateTodoInput{Title: "test updated title 1", Content: "test updated content 1"}
	res, err := testTodoService.UpdateTodo(ctx, testTodo.ID, requestParams, user.ID)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), strconv.Itoa(testTodo.ID), res.ID)
	var expectedValidationErrors []string
	assert.Equal(s.T(), expectedValidationErrors, res.ValidationErrors.Title)
	assert.Equal(s.T(), expectedValidationErrors, res.ValidationErrors.Content)

	// NOTE: TODOが更新されていることの確認
	if err := testTodo.Reload(ctx, DBCon); err != nil {
		s.T().Fatalf("failed to reload test todos %v", err)
	}
	assert.Equal(s.T(), "test updated title 1", testTodo.Title)
	assert.Equal(s.T(), null.String{String: "test updated content 1", Valid: true}, testTodo.Content)
}

func (s *TestTodoServiceSuite) TestUpdateTodo_ValidationError() {
	testTodo := models.Todo{Title: "test title 1", Content: null.String{String: "test content 1", Valid: true}, UserID: user.ID}
	if err := testTodo.Insert(ctx, DBCon, boil.Infer()); err != nil {
		s.T().Fatalf("failed to create test todos %v", err)
	}

	requestParams := model.UpdateTodoInput{Title: "", Content: "test updated content 1"}
	res, err := testTodoService.UpdateTodo(ctx, testTodo.ID, requestParams, user.ID)

	assert.Nil(s.T(), err)
	assert.Empty(s.T(), res.ID)
	expectedTitleValidationErrors := []string{"タイトルは必須入力です。"}
	var expectedContentValidationErrors []string
	assert.Equal(s.T(), expectedTitleValidationErrors, res.ValidationErrors.Title)
	assert.Equal(s.T(), expectedContentValidationErrors, res.ValidationErrors.Content)

	// NOTE: Todoが更新されていないこと
	if err := testTodo.Reload(ctx, DBCon); err != nil {
		s.T().Fatalf("failed to reload test todos %v", err)
	}
	assert.Equal(s.T(), "test title 1", testTodo.Title)
	assert.Equal(s.T(), null.String{String: "test content 1", Valid: true}, testTodo.Content)
}

func (s *TestTodoServiceSuite) TestDeleteTodo() {
	testTodo := models.Todo{Title: "test title 1", Content: null.String{String: "test content 1", Valid: true}, UserID: user.ID}
	if err := testTodo.Insert(ctx, DBCon, boil.Infer()); err != nil {
		s.T().Fatalf("failed to create test todos %v", err)
	}

	id, err := testTodoService.DeleteTodo(ctx, testTodo.ID, user.ID)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), strconv.Itoa(testTodo.ID), id)
	// NOTE: TODOが削除されていることの確認
	reloadErr := testTodo.Reload(ctx, DBCon)
	assert.NotNil(s.T(), reloadErr)
}

func TestTodoService(t *testing.T) {
	// テストスイートを実行
	suite.Run(t, new(TestTodoServiceSuite))
}
