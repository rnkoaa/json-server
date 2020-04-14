package client

import (
	"context"
	"fmt"
	"github.com/rnkoaa/json-server/pkg/domain"
)

type TodoService interface {
	Get(ctx context.Context, id int) (*domain.Todo, *Response, error)
	GetByUser(context.Context, int) ([]*domain.Todo, *Response, error)
	List(ctx context.Context, id int) ([]*domain.Todo, *Response, error)
}

type TodoServiceOp struct {
	client *Client
}


var (
	todosBaseURL   = "/todos"
	todosURLFormat = "/todos/%d"
)

var _ TodoService = &TodoServiceOp{}

func (u *TodoServiceOp) Get(ctx context.Context, id int) (*domain.Todo, *Response, error) {
	path := fmt.Sprintf(todosURLFormat, id)
	var todo domain.Todo

	res, err := get(ctx, u.client, path, &todo)
	if err != nil {
		return nil, nil, err
	}
	return &todo, res, nil
}


func (u *TodoServiceOp) GetByUser(ctx context.Context, userId int) ([]*domain.Todo, *Response, error) {
	path := fmt.Sprintf(userURLFormat, userId)
	path = fmt.Sprintf("%s/todos", path)
	var todos []*domain.Todo

	res, err := get(ctx, u.client, path, &todos)
	if err != nil {
		return nil, nil, err
	}
	return todos, res, nil
}

func (u *TodoServiceOp) List(ctx context.Context, id int) ([]*domain.Todo, *Response, error) {
	var todos []*domain.Todo

	res, err := get(ctx, u.client, todosBaseURL, &todos)
	if err != nil {
		return nil, nil, err
	}
	return todos, res, nil
}
