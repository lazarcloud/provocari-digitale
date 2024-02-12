package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/lazarcloud/provocari-digitale/api/auth/jwt"
	"github.com/lazarcloud/provocari-digitale/api/globals"
	"github.com/lazarcloud/provocari-digitale/api/utils"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.RespondWithUnauthorized(w, "Missing Authorization header")
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			utils.RespondWithUnauthorized(w, "Invalid Authorization header")
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == "" {
			utils.RespondWithUnauthorized(w, "Missing token")
			return
		}

		parsedJWT, err := jwt.ParseJWT(tokenStr)
		if err != nil {
			utils.RespondWithUnauthorized(w, "Error decrypting token: "+err.Error())
			return
		}

		ctx := context.WithValue(r.Context(), globals.ContextAccessRoleKey, parsedJWT.AccessRole)

		ctx = context.WithValue(ctx, globals.ContextUserIdKey, parsedJWT.UserId)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func GetRole(r *http.Request) string {
	val := r.Context().Value(globals.ContextAccessRoleKey)
	if val != nil {
		return val.(string)
	}
	return ""
}

func GetUserId(r *http.Request) string {
	val := r.Context().Value(globals.ContextUserIdKey)
	if val != nil {
		return val.(string)
	}
	return ""
}

func ServiceAccountMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role := GetRole(r)
		if role != globals.AuthRoleService {
			utils.RespondWithUnauthorized(w, "Invalid role, service only route")
			return
		}
		next.ServeHTTP(w, r)
	})
}
