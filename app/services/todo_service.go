package services

import (
	"app/graph/model"
	models "app/models/generated"
	"app/view"
	"context"
	"database/sql"
	"strconv"

	"app/validator"

	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type TodoService interface {
	CreateTodo(ctx context.Context, requestParams model.CreateTodoInput, userID int) (*models.Todo, error)
	FetchTodoLists(ctx context.Context, userID int) ([]*models.Todo, error)
	FetchTodo(ctx context.Context, id int, userID int) (*models.Todo, error)
	UpdateTodo(ctx context.Context, id int, requestParams model.UpdateTodoInput, userID int) (*models.Todo, error)
	DeleteTodo(ctx context.Context, id int, userID int) (string, error)
}

type todoService struct {
	db *sql.DB
}

func NewTodoService(db *sql.DB) TodoService {
	return &todoService{db}
}

func (ts *todoService) CreateTodo(ctx context.Context, requestParams model.CreateTodoInput, userID int) (*models.Todo, error) {
	// NOTE: バリデーションチェック
	validationErrors := validator.ValidateCreateTodo(requestParams)
	if validationErrors != nil {
		return &models.Todo{}, view.NewBadRequestView(validationErrors)
	}

	todo := &models.Todo{}
	todo.Title = requestParams.Title
	todo.Content = null.String{String: requestParams.Content, Valid: true}
	todo.UserID = userID

	// NOTE: Create処理
	err := todo.Insert(ctx, ts.db, boil.Infer())
	if err != nil {
		return &models.Todo{}, view.NewInternalServerErrorView(err)
	}
	return todo, nil
}

func (ts *todoService) FetchTodoLists(ctx context.Context, userID int) ([]*models.Todo, error) {
	todos, err := models.Todos(qm.Where("user_id = ?", userID)).All(ctx, ts.db)
	if err != nil {
		return models.TodoSlice{}, view.NewNotFoundView(err)
	}

	return todos, nil
}

func (ts *todoService) FetchTodo(ctx context.Context, id int, userID int) (*models.Todo, error) {
	todo, err := models.Todos(qm.Where("id = ? AND user_id = ?", id, userID)).One(ctx, ts.db)
	if err != nil {
		return &models.Todo{}, view.NewNotFoundView(err)
	}

	return todo, nil
}

func (ts *todoService) UpdateTodo(ctx context.Context, id int, requestParams model.UpdateTodoInput, userID int) (*models.Todo, error) {
	todo, err := models.Todos(qm.Where("id = ? AND user_id = ?", id, userID)).One(ctx, ts.db)
	if err != nil {
		return &models.Todo{}, view.NewNotFoundView(err)
	}

	// NOTE: バリデーションチェック
	validationErrors := validator.ValidateUpdateTodo(requestParams)
	if validationErrors != nil {
		return &models.Todo{}, view.NewBadRequestView(validationErrors)
	}

	todo.Title = requestParams.Title
	todo.Content = null.String{String: requestParams.Content, Valid: true}

	// NOTE: Update処理
	_, updateError := todo.Update(ctx, ts.db, boil.Infer())
	if updateError != nil {
		return &models.Todo{}, view.NewInternalServerErrorView(updateError)
	}
	return todo, nil
}

func (ts *todoService) DeleteTodo(ctx context.Context, id int, userID int) (string, error) {
	todo, err := models.Todos(qm.Where("id = ? AND user_id = ?", id, userID)).One(ctx, ts.db)
	if err != nil {
		return strconv.Itoa(id), view.NewNotFoundView(err)
	}

	_, deleteError := todo.Delete(ctx, ts.db)
	if deleteError != nil {
		return strconv.Itoa(id), view.NewInternalServerErrorView(deleteError)
	}
	return strconv.Itoa(id), nil
}
