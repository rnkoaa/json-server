package client

import (
	"context"
	"fmt"
	"github.com/rnkoaa/json-server/pkg/domain"
)

type PhotoService interface {
	Get(ctx context.Context, id int) (*domain.Photo, *Response, error)
	GetByAlbum(context.Context, int) ([]*domain.Photo, *Response, error)
	List(ctx context.Context, id int) ([]*domain.Photo, *Response, error)
}

type PhotoServiceOp struct {
	client *Client
}

var (
	photosBaseURL   = "/photos"
	photosURLFormat = "/photos/%d"
)

var _ PhotoService = &PhotoServiceOp{}

func (u *PhotoServiceOp) Get(ctx context.Context, id int) (*domain.Photo, *Response, error) {
	path := fmt.Sprintf(photosURLFormat, id)
	var photo domain.Photo

	res, err := get(ctx, u.client, path, &photo)
	if err != nil {
		return nil, nil, err
	}
	return &photo, res, nil
}

func (u *PhotoServiceOp) GetByAlbum(ctx context.Context, userId int) ([]*domain.Photo, *Response, error) {
	path := fmt.Sprintf(albumsURLFormat, userId)
	path = fmt.Sprintf("%s/photos", path)
	var photos []*domain.Photo

	res, err := get(ctx, u.client, path, &photos)
	if err != nil {
		return nil, nil, err
	}
	return photos, res, nil
}

func (u *PhotoServiceOp) List(ctx context.Context, id int) ([]*domain.Photo, *Response, error) {
	var photos []*domain.Photo

	res, err := get(ctx, u.client, photosBaseURL, &photos)
	if err != nil {
		return nil, nil, err
	}
	return photos, res, nil
}
