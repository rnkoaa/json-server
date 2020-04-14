package main

import (
	"context"
	"fmt"
	"github.com/rnkoaa/json-server/pkg/domain"
)

type CommentService interface {
	Get(ctx context.Context, id int) (*domain.Comment, *Response, error)
	GetByPost(context.Context, int) ([]*domain.Comment, *Response, error)
	List(ctx context.Context, id int) ([]*domain.Comment, *Response, error)
}

type CommentServiceOp struct {
	client *Client
}

var (
	commentsBaseURL   = "/comments"
	commentsURLFormat = "/comments/%d"
)

var _ CommentService = &CommentServiceOp{}

func (u *CommentServiceOp) Get(ctx context.Context, id int) (*domain.Comment, *Response, error) {
	path := fmt.Sprintf(commentsURLFormat, id)
	var comment domain.Comment

	res, err := get(ctx, u.client, path, &comment)
	if err != nil {
		return nil, nil, err
	}
	return &comment, res, nil
}

func (u *CommentServiceOp) GetByPost(ctx context.Context, commentId int) ([]*domain.Comment, *Response, error) {
	path := fmt.Sprintf(postsURLFormat, commentId)
	path = fmt.Sprintf("%s/comments", path)
	var comments []*domain.Comment

	res, err := get(ctx, u.client, path, &comments)
	if err != nil {
		return nil, nil, err
	}
	return comments, res, nil
}

func (u *CommentServiceOp) List(ctx context.Context, id int) ([]*domain.Comment, *Response, error) {
	var comments []*domain.Comment

	res, err := get(ctx, u.client, commentsBaseURL, &comments)
	if err != nil {
		return nil, nil, err
	}
	return comments, res, nil
}
