package auth

import (
	models "app/models/generated"
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
)

type contextKey struct {
	uuid string
}

var (
	signInWriterKey = contextKey{"signInWriter"}
	userKey         = contextKey{"user"}
	authCookieKey   = "token"
)

func Middleware(next http.Handler, db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// NOTE: ログイン時のCookieをResponseWriterでセットするためのcontextをセット
		ctx := context.WithValue(r.Context(), signInWriterKey, w)
		r = r.WithContext(ctx)

		// NOTE: リクエストからCookieを取得
		tokenString, err := r.Cookie(authCookieKey)
		// NOTE: 未認証でもアクセス可のページ
		if err != nil || tokenString == nil {
			next.ServeHTTP(w, r)
			return
		}

		// NOTE: tokenに該当するユーザを取得する
		token, err := jwt.Parse(tokenString.Value, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(os.Getenv("JWT_TOKEN_KEY")), nil
		})
		if err != nil {
			next.ServeHTTP(w, r)
			fmt.Println("failt jwt parse")
			return
		}

		var userID int
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID = int(claims["user_id"].(float64))
		}
		if userID == 0 {
			next.ServeHTTP(w, r)
			fmt.Println("invalid token")
			return
		}
		user, _ := models.FindUser(ctx, db, userID)

		// NOTE: Contextにuserをセットする
		withUserContext := context.WithValue(r.Context(), userKey, user)
		r = r.WithContext(withUserContext)

		next.ServeHTTP(w, r)
	})
}

func GetUser(ctx context.Context) *models.User {
	user, _ := ctx.Value(userKey).(*models.User)
	return user
}

func SetAuthCookie(ctx context.Context, token string) {
	w, _ := ctx.Value(signInWriterKey).(http.ResponseWriter)

	week := 60 * 60 * 24 * 7

	cookie := http.Cookie{
		HttpOnly: true,
		MaxAge:   week * 2,
		Secure:   true,
		Name:     authCookieKey,
		Value:    token,
	}
	http.SetCookie(w, &cookie)
}
