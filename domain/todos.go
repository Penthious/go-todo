package domain

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v9"
	"time"
)

type Todo struct {
	ID        int64       `json:"id"`
	Title     string      `json:"title"`
	Completed bool        `json:"completed" pg:",use_zero"`
	UserID    int64       `json:"user_id"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
	DeletedAt pg.NullTime `json:"deletedAt" pg:",soft_delete"`
}

func (t *Todo) BeforeUpdate(ctx context.Context) (context.Context, error) {
	t.UpdatedAt = time.Now()
	return ctx, nil
}

type CreateTodoPayload struct {
	Title string `json:"title"`
}

func (todo *CreateTodoPayload) IsValid() (bool, map[string]string) {
	v := NewValidator()

	v.ValueIsRequired("title", todo.Title)
	v.MustBeLongerThan("title", todo.Title, 3)

	return v.IsValid(), v.errors
}

func (d *Domain) CreateTodo(payload CreateTodoPayload, user *User) (*Todo, error) {
	data := &Todo{
		Title:     payload.Title,
		Completed: false,
		UserID:    user.ID,
	}

	fmt.Printf("DATA: %v", data)
	todo, err := d.DB.TodoRepo.Create(data)

	if err != nil {
		return nil, err
	}

	return todo, nil
}

func (d *Domain) GetTodoByID(id int64) (*Todo, error) {
	todo, err := d.DB.TodoRepo.GetByID(id)

	if err != nil {
		return nil, err
	}

	return todo, nil
}

func (d *Domain) DeleteTodo(todo *Todo) error {

	err := d.DB.TodoRepo.Delete(todo)

	if err != nil {
		return err
	}

	return nil
}

type UpdateTodoPayload struct {
	Title     *string `json:"title"`
	Completed *bool   `json:"completed"`
}

func (todo *UpdateTodoPayload) IsValid() (bool, map[string]string) {
	v := NewValidator()

	if *todo.Title != "" {
		v.MustBeLongerThan("title", *todo.Title, 3)
	}

	return v.IsValid(), v.errors
}

func (d *Domain) UpdateTodo(payload UpdateTodoPayload, todo *Todo) (*Todo, error) {

	if *payload.Title != "" {
		todo.Title = *payload.Title
	}
	todo.Completed = *payload.Completed
	todo, err := d.DB.TodoRepo.Update(todo)

	if err != nil {
		return nil, err
	}

	return todo, err
}

func (t *Todo) IsOwner(user *User) bool {
	return t.UserID == user.ID
}
