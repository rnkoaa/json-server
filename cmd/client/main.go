package main

import (
	"context"
	"fmt"
	"github.com/rnkoaa/json-server/pkg/strings"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	//httpClient *http.Client
	client *Client
)

type Result struct {
	Error    error
	Response interface{}
}

func doRequests(ctx context.Context, urls []int) <-chan Result {
	responses := make(chan Result)
	go func() {
		defer close(responses)
		for _, userId := range urls {
			res, _, err := client.User.Get(ctx, userId)
			result := Result{Error: err, Response: res}
			select {
			case <-ctx.Done():
				return
			case responses <- result:
			}
		}
	}()

	return responses
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	client = NewClient(http.DefaultClient, "")
	fetchAllUsers(ctx, client)
	//fetchAllTodos(ctx, client)
	fetchTodosByUser(ctx, client, 1)
	fetchPostsByUser(ctx, client, 1)
	fetchAlbumsByUser(ctx, client, 1)
	fetchCommentsByPost(ctx, client, 1)
	fetchPhotosByAlbum(ctx, client, 1)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	// block until either an error or OS-level signals
	// to shutdown gracefully
	select {
	case <-sigChan:
		fmt.Printf("Shutdown signal received... closing server")
		cancel()
	}
}

func fetchTodosByUser(ctx context.Context, c *Client, userId int) {
	terminated := doFetchUserTodoRequests(ctx, c, userId)
	errCount := 0
	for i := range terminated {
		if i.Error != nil {
			fmt.Println(i.Error)
			errCount++
			if errCount >= 3 {
				fmt.Println("Too many errors, breaking!")
				break
			}
			continue
		}
		fmt.Println(strings.Stringify(i.Response))
	}
}

func fetchPostsByUser(ctx context.Context, c *Client, userId int) {
	terminated := doFetchUserPostRequests(ctx, c, userId)
	errCount := 0
	for i := range terminated {
		if i.Error != nil {
			fmt.Println(i.Error)
			errCount++
			if errCount >= 3 {
				fmt.Println("Too many errors, breaking!")
				break
			}
			continue
		}
		fmt.Println(strings.Stringify(i.Response))
	}
}

func fetchAlbumsByUser(ctx context.Context, c *Client, userId int) {
	terminated := doFetchUserAlbumRequests(ctx, c, userId)
	errCount := 0
	for i := range terminated {
		if i.Error != nil {
			fmt.Println(i.Error)
			errCount++
			if errCount >= 3 {
				fmt.Println("Too many errors, breaking!")
				break
			}
			continue
		}
		fmt.Println(strings.Stringify(i.Response))
	}
}

func fetchCommentsByPost(ctx context.Context, c *Client, postId int) {
	terminated := doFetchPostCommentsRequests(ctx, c, postId)
	errCount := 0
	for i := range terminated {
		if i.Error != nil {
			fmt.Println(i.Error)
			errCount++
			if errCount >= 3 {
				fmt.Println("Too many errors, breaking!")
				break
			}
			continue
		}
		fmt.Println(strings.Stringify(i.Response))
	}
}

func fetchPhotosByAlbum(ctx context.Context, c *Client, postId int) {
	terminated := doFetchPhotosByAlbumRequests(ctx, c, postId)
	errCount := 0
	for i := range terminated {
		if i.Error != nil {
			fmt.Println(i.Error)
			errCount++
			if errCount >= 3 {
				fmt.Println("Too many errors, breaking!")
				break
			}
			continue
		}
		fmt.Println(strings.Stringify(i.Response))
	}
}
