package lib

import (
	"app/graph"
	"app/graph/generated"
	"app/lib/auth"
	"app/lib/dataloaders"
	"app/services"
	"app/view"
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/vektah/gqlparser/v2/gqlerror"

	"github.com/rs/cors"
)

func GetGraphQLHttpHandler(db *sql.DB) http.Handler {
	// NOTE: service
	authService := services.NewAuthService(db)
	todoService := services.NewTodoService(db)

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(authService, todoService)}))

	srv.SetErrorPresenter(func(ctx context.Context, e error) *gqlerror.Error {
		err := graphql.DefaultErrorPresenter(ctx, e)

		var errorCode int64
		var error error

		var re view.ViewError
		if errors.As(e, &re) {
			errorCode = re.Code
			error = re.Message
		}

		err.Extensions = map[string]interface{}{
			"code":  errorCode,
			"error": error,
		}

		return err
	})
	
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:8080"},
		AllowCredentials: true,
	})

	dataloaderSrv := dataloaders.Middleware(c.Handler(srv), db)
	return auth.Middleware(dataloaderSrv, db)
}
