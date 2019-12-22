package handlers

import (
	"net/http"
	"todo/domain"
)

func (s *Server) createTodo() http.HandlerFunc {

	var payload domain.CreateTodoPayload
	return validatePayload(func(w http.ResponseWriter, r *http.Request) {
		user := s.getUserFromCTX(*r)
		todo, err := s.domain.CreateTodo(payload, user)

		if err != nil {
			badRequestResponse(w, err)
			return
		}
		jsonResponse(w, todo, http.StatusCreated)
	}, &payload)
}
