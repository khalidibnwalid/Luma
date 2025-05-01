package middlewares

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/khalidibnwalid/Luma/core"
)

type key string

const CtxUserIDKey key = "auth.JWT_USER_ID"

// TODO token refreshing
func JwtAuthBuilder(secret string) core.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(core.JwtSessionCookieName)
			if err != nil {
				unahuthorized(w)
				return
			}

			tokenString := strings.Replace(cookie.Value, "Bearer ", "", 1)
			if tokenString == "" {
				unahuthorized(w)
				return
			}

			_, claims, ok := core.ValidateJwtToken(secret, tokenString)
			if !ok || tokenExpired(claims["exp"].(float64)) {
				unahuthorized(w)
				return
			}

			userId, _ := claims.GetSubject()
			uuidId, err := uuid.Parse(userId)
			// technically, this is not needed, since it is inside the JWT
			if err != nil {
				unahuthorized(w)
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), CtxUserIDKey, uuidId))
			next.ServeHTTP(w, r)
		})
	}
}

func unahuthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
}

func tokenExpired(exp float64) bool {
	expInt := int64(exp)
	return time.Now().Unix() > expInt
}
