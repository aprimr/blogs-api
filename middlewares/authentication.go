package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/aprimr/blogs-api/utils"
	"github.com/aprimr/blogs-api/validation"
)

func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Get auth header and trim prefix
		authHeader := r.Header.Get("Authorization")
		if validation.IsEmptyString(authHeader) {
			utils.SendError(w, "JWT token missing", http.StatusUnauthorized)
			return
		}
		jwtToken := strings.TrimPrefix(authHeader, "Bearer ")

		// Verify JWT Token
		jwtClaims, err := utils.VerifyToken(jwtToken)
		if err != nil {
			utils.SendError(w, "JWT token mismatched", http.StatusUnauthorized)
			return
		}

		// Send jwtClaims with next request
		ctx := context.WithValue(r.Context(), "uid", jwtClaims.Uid)
		ctx = context.WithValue(ctx, "email", jwtClaims.Email)

		// call next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
