package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
	"todo/domain"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Server struct {
	domain *domain.Domain
}

func setupMiddleware(r *chi.Mux) {
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Compress(6, "application/json"))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(middleware.Timeout(60 * time.Second))
}

func NewServer(domain *domain.Domain) *Server {
	return &Server{domain: domain}
}

func SetupRouter(domain *domain.Domain) *chi.Mux {
	server := NewServer(domain)

	r := chi.NewRouter()
	setupMiddleware(r)
	server.setupEndpoints(r)

	return r
}

func jsonResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(statusCode)

	if data == nil {
		data = map[string]string{}
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}

func badRequestResponse(w http.ResponseWriter, err error) {
	response := map[string]string{"error": err.Error()}
	jsonResponse(w, response, http.StatusBadRequest)
}

func unauthorizedResponse(w http.ResponseWriter) {
	response := map[string]string{"error": "unauthorized"}
	jsonResponse(w, response, http.StatusUnauthorized)
}

func forbiddenResponse(w http.ResponseWriter) {
	response := map[string]string{"error": "forbidden"}
	jsonResponse(w, response, http.StatusForbidden)
}

type PayloadValidation interface {
	IsValid() (bool, map[string]string)
}

// validatePayload decode the body to json and validate each field
func validatePayload(next http.HandlerFunc, payload PayloadValidation) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		err := json.NewDecoder(r.Body).Decode(&payload)

		if err != nil {
			badRequestResponse(w, err)
			return
		}

		defer r.Body.Close()

		if valid, errs := payload.IsValid(); !valid {
			jsonResponse(w, errs, http.StatusBadRequest)

			return
		}

		ctx := context.WithValue(r.Context(), "payload", payload)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (s *Server) withOwner(subjectType string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			currentUser := s.getUserFromCTX(r)
			subject := r.Context().Value(subjectType).(domain.HaveOwner)

			if !subject.IsOwner(currentUser) {

				forbiddenResponse(w)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

}
