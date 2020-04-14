package main

import (
	"context"
	"fmt"
	"github.com/rnkoaa/json-server/pkg/domain"
)

type PostService interface {
	Get(ctx context.Context, id int) (*domain.Post, *Response, error)
	GetByUser(context.Context, int) ([]*domain.Post, *Response, error)
	List(ctx context.Context, id int) ([]*domain.Post, *Response, error)
	GetComments(ctx context.Context, id int) ([]*domain.Comment, *Response, error)
}

type PostServiceOp struct {
	client *Client
}

var (
	postsBaseURL   = "/posts"
	postsURLFormat = "/posts/%d"
)

var _ PostService = &PostServiceOp{}

func (p *PostServiceOp) GetComments(ctx context.Context, id int) ([]*domain.Comment, *Response, error) {
	return p.client.Comment.GetByPost(ctx, id)
}

func (u *PostServiceOp) Get(ctx context.Context, id int) (*domain.Post, *Response, error) {
	path := fmt.Sprintf(postsURLFormat, id)
	var post domain.Post

	res, err := get(ctx, u.client, path, &post)
	if err != nil {
		return nil, nil, err
	}
	return &post, res, nil
}

func (u *PostServiceOp) GetByUser(ctx context.Context, userId int) ([]*domain.Post, *Response, error) {
	path := fmt.Sprintf(userURLFormat, userId)
	path = fmt.Sprintf("%s/posts", path)
	var posts []*domain.Post

	res, err := get(ctx, u.client, path, &posts)
	if err != nil {
		return nil, nil, err
	}
	return posts, res, nil
}

func (u *PostServiceOp) List(ctx context.Context, id int) ([]*domain.Post, *Response, error) {
	var posts []*domain.Post

	res, err := get(ctx, u.client, postsBaseURL, &posts)
	if err != nil {
		return nil, nil, err
	}
	return posts, res, nil
}
