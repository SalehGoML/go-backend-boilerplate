package middleware

import (
	utils "Salehaskarzadeh/internal/Utils"
	"Salehaskarzadeh/internal/storee"
	"Salehaskarzadeh/internal/tokens"
	"context"
	"net/http"
	"strings"
)

type UserMiddleware struct {
	UserStore storee.UserStore
}

type contextKey string

const UserContexKey = contextKey("user")

func SetUser(r *http.Request, user *storee.User) *http.Request {
	ctx := context.WithValue(r.Context(), UserContexKey, user)
	return r.WithContext(ctx)
}

func GetUser(r *http.Request) *storee.User {
	user, ok := r.Context().Value(UserContexKey).(*storee.User)
	if !ok {
		panic("missing user in request") // bad actor call
	}
	return user
}

func (um *UserMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// within this anonymouse function
		// we can interject any incoming requests to our server

		w.Header().Add("Vary", "Authorization")
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			r = SetUser(r, storee.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authHeader, " ") // Bearer <TOKEN>
		if len(headerParts) != 2 || strings.TrimSpace(headerParts[0]) != "Bearer" {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid authorization header"})
			return
		}

		token := headerParts[1]
		user, err := um.UserStore.GetUserToken(tokens.ScopeAuth, token)
		if err != nil {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid token"})
			return
		}
		if user == nil {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "token expired or invalid"})
			return
		}

		r = SetUser(r, user)
		next.ServeHTTP(w, r)

	})

}

func (um *UserMiddleware) RequireUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUser(r)

		if user.IsAnonymous() {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "you must be logged in to access this route"})
			return
		}

		next.ServeHTTP(w, r)
	})
}
