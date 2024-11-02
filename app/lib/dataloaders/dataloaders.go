package dataloaders

import (
	models "app/models/generated"
	"context"
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/vikstrous/dataloadgen"
)

type ctxKey string

const (
	loadersKey = ctxKey("dataloaders")
)

// userReader reads Users from a database
type userReader struct {
	db *sql.DB
}

func (u *userReader) getUsers(ctx context.Context, userIDs []string) ([]*models.User, []error) {
	intUserIDs := make([]int, len(userIDs))
	for n := range userIDs {
		intUserIDs[n], _ = strconv.Atoi(userIDs[n])
	}
	users, err := models.Users(models.UserWhere.ID.IN(intUserIDs)).All(ctx, u.db)
	if err != nil {
		return nil, []error{err}
	}
	return users, nil
}

type Loaders struct {
	UserLoader *dataloadgen.Loader[string, *models.User]
}

func NewLoaders(db *sql.DB) *Loaders {
	ur := &userReader{db: db}
	return &Loaders{
		UserLoader: dataloadgen.NewLoader(ur.getUsers, dataloadgen.WithWait(time.Millisecond)),
	}
}

func Middleware(next http.Handler, db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loader := NewLoaders(db)
		r = r.WithContext(context.WithValue(r.Context(), loadersKey, loader))
		next.ServeHTTP(w, r)
	})
}

func For(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey).(*Loaders)
}

func GetUser(ctx context.Context, userID string) (*models.User, error) {
	loaders := For(ctx)
	return loaders.UserLoader.Load(ctx, userID)
}

func GetUsers(ctx context.Context, userIDs []string) ([]*models.User, error) {
	loaders := For(ctx)
	return loaders.UserLoader.LoadAll(ctx, userIDs)
}
