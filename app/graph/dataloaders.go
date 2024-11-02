package graph

import (
	models "app/models/generated"
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/graph-gophers/dataloader/v7"
)

type Loaders struct {
	UserByID *dataloader.Loader[string, *models.User]
}

// Middleware dataloaderはメモリキャッシュを返すらしい オプションでOFFにもできるらしい
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), loadersKey, &Loaders{
			UserByID: dataloader.NewBatchedLoader(func(ctx context.Context, userIds []string) []*dataloader.Result[*models.User] {
				fmt.Println("batch get users:", userIds)

				// ユーザIDのリストからユーザ情報を取得する
				// サンプル実装なので適当な値を返していますが、プロダクト実装では以下のようにしてください。
				//   - "SELECT * FROM users WHERE id IN (id1, id2, id3, ...)"のようなSQLでDBからユーザ情報を一括取得する
				//   - 他のサービスのBatch Read APIを呼ぶ
				// それでN+1問題を回避することができます。
				results := make([]*dataloader.Result[*models.User], len(userIds))
				for i, id := range userIds {
					intID, _ := strconv.Atoi(id)
					results[i] = &dataloader.Result[*models.User]{
						Data:  &models.User{ID: intID, Name: "user " + id},
						Error: nil,
					}
				}

				return results
			}),
		})
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

type contextKey int

var loadersKey contextKey

func CtxLoaders(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey).(*Loaders)
}
