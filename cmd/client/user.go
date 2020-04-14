package main

import (
	"context"
	"fmt"
	"github.com/rnkoaa/json-server/pkg/domain"
)

var (
	userURLFormat = "/users/%d"
	userBaseURLFormat = "/users/%d"
)

type UserService interface {
	Get(ctx context.Context, id int) (*domain.User, *Response, error)
	List(ctx context.Context, id int) ([]*domain.User, *Response, error)
	GetPosts(ctx context.Context, id int) ([]*domain.Post, *Response, error)
	GetAlbums(ctx context.Context, id int) ([]*domain.Album, *Response, error)
	GetTodos(ctx context.Context, id int) ([]*domain.Todo, *Response, error)
}

type UserServiceOp struct {
	client *Client
}

var _ UserService = &UserServiceOp{}

func (u *UserServiceOp) GetPosts(ctx context.Context, id int) ([]*domain.Post, *Response, error) {
	return u.client.Post.GetByUser(ctx, id)
}

func (u *UserServiceOp) GetAlbums(ctx context.Context, id int) ([]*domain.Album, *Response, error) {
	return u.client.Album.GetByUser(ctx, id)
}

func (u *UserServiceOp) GetTodos(ctx context.Context, id int) ([]*domain.Todo, *Response, error) {
	return u.client.Todo.GetByUser(ctx, id)
}

func (u *UserServiceOp) Get(ctx context.Context, id int) (*domain.User, *Response, error) {
	path := fmt.Sprintf(userURLFormat, id)
	var user domain.User

	res, err := get(ctx, u.client, path, &user)
	if err != nil {
		return nil, nil, err
	}
	return &user, res, nil

}

func (u *UserServiceOp) List(ctx context.Context, id int) ([]*domain.User, *Response, error) {
	var users []*domain.User

	res, err := get(ctx, u.client, todosBaseURL, &users)
	if err != nil {
		return nil, nil, err
	}
	return users, res, nil
}
