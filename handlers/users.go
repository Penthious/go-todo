package handlers

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"todo/domain"
)

type authResponse struct {
	Token *domain.JWTToken `json:"token"`
}

func (s *Server) registerUser() http.HandlerFunc {
	var payload domain.RegisterPayload

	return validatePayload(func(w http.ResponseWriter, r *http.Request) {
		user, err := s.domain.Register(payload)

		if err != nil {
			badRequestResponse(w, err)
			return
		}
		jwtToken, err := user.GenToken()
		if err != nil {
			badRequestResponse(w, err)
			return
		}
		jsonResponse(w, &authResponse{Token: jwtToken}, http.StatusCreated)
	}, &payload)
}

func (s *Server) getUserFromCTX(r http.Request) *domain.User {
	currentUser := r.Context().Value("currentUser").(*domain.User)

	return currentUser
}
func (s *Server) withUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := domain.ParseToken(r)

		if err != nil {
			unauthorizedResponse(w)
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

			user := domain.User{}
			mapstructure.Decode(claims["user"], &user)

			DBuser, err := s.domain.GetUserByID(user.ID)
			if err != nil {
				unauthorizedResponse(w)
				return
			}
			ctx := context.WithValue(r.Context(), "currentUser", DBuser)

			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			unauthorizedResponse(w)
			return
		}
	})
}
