package client

import (
	"context"
	"fmt"
	"github.com/rnkoaa/json-server/pkg/domain"
)

type AlbumService interface {
	Get(ctx context.Context, id int) (*domain.Album, *Response, error)
	GetByUser(context.Context, int) ([]*domain.Album, *Response, error)
	GetPhotos(context.Context, int) ([]*domain.Photo, *Response, error)
	List(ctx context.Context, id int) ([]*domain.Album, *Response, error)
}

type AlbumServiceOp struct {
	client *Client
}

var (
	albumsBaseURL   = "/albums"
	albumsURLFormat = "/albums/%d"
)

var _ AlbumService = &AlbumServiceOp{}

func (a *AlbumServiceOp) Get(ctx context.Context, id int) (*domain.Album, *Response, error) {
	path := fmt.Sprintf(albumsURLFormat, id)
	var album domain.Album

	res, err := get(ctx, a.client, path, &album)
	if err != nil {
		return nil, nil, err
	}
	return &album, res, nil
}

func (a *AlbumServiceOp) GetPhotos(ctx context.Context, albumId int) ([]*domain.Photo, *Response, error) {
	return a.client.Photo.GetByAlbum(ctx, albumId)
}

func (a *AlbumServiceOp) GetByUser(ctx context.Context, userId int) ([]*domain.Album, *Response, error) {
	path := fmt.Sprintf(userURLFormat, userId)
	path = fmt.Sprintf("%s/albums", path)
	var albums []*domain.Album

	res, err := get(ctx, a.client, path, &albums)
	if err != nil {
		return nil, nil, err
	}
	return albums, res, nil
}

func (a *AlbumServiceOp) List(ctx context.Context, id int) ([]*domain.Album, *Response, error) {
	var albums []*domain.Album

	res, err := get(ctx, a.client, albumsBaseURL, &albums)
	if err != nil {
		return nil, nil, err
	}
	return albums, res, nil
}
