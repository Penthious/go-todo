package handlers

import (
	"context"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
	"todo/domain"
)

func (s *Server) createTodo() http.HandlerFunc {

	var payload domain.CreateTodoPayload
	return validatePayload(func(w http.ResponseWriter, r *http.Request) {
		user := s.getUserFromCTX(r)
		todo, err := s.domain.CreateTodo(payload, user)

		if err != nil {
			badRequestResponse(w, err)
			return
		}
		jsonResponse(w, todo, http.StatusCreated)
	}, &payload)
}

func (s *Server) updateTodo() http.HandlerFunc {
	var payload domain.UpdateTodoPayload
	return validatePayload(func(w http.ResponseWriter, r *http.Request) {
		todo, err := s.domain.UpdateTodo(payload, s.getTodoFromCTX(r))

		if err != nil {
			badRequestResponse(w, err)
			return
		}
		jsonResponse(w, todo, http.StatusOK)
	}, &payload)
}

func (s *Server) deleteTodo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		todo := s.getTodoFromCTX(r)
		err := s.domain.DB.TodoRepo.Delete(todo)
		if err != nil {
			badRequestResponse(w, err)
			return
		}

		jsonResponse(w, nil, http.StatusNoContent)
	}
}
func (s *Server) todoCTX(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//todo := new(domain.Todo)
		if todoID := chi.URLParam(r, "id"); todoID != "" {
			id, err := strconv.ParseInt(todoID, 0, 0)
			if err != nil {
				badRequestResponse(w, err)
				return
			}

			todo, err := s.domain.GetTodoByID(id)

			if err != nil {
				response := map[string]string{"error": domain.ErrNoResult.Error()}
				jsonResponse(w, response, http.StatusNotFound)
			}

			ctx := context.WithValue(r.Context(), "todo", todo)

			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}

func (s *Server) getTodoFromCTX(r *http.Request) *domain.Todo {
	todo := r.Context().Value("todo").(*domain.Todo)

	return todo
}
