package main

import (
	"context"
	"fmt"
	"github.com/rnkoaa/json-server/pkg/strings"
)

func doFetchUserTodoRequests(ctx context.Context, client *Client, userId int) <-chan Result {
	responses := make(chan Result)
	go func() {
		defer close(responses)
		res, _, err := client.Todo.GetByUser(ctx, userId)
		result := Result{Error: err, Response: res}
		select {
		case <-ctx.Done():
			return
		case responses <- result:
		}
	}()

	return responses
}
func doFetchUserPostRequests(ctx context.Context, client *Client, userId int) <-chan Result {
	responses := make(chan Result)
	go func() {
		defer close(responses)
		res, _, err := client.Post.GetByUser(ctx, userId)
		result := Result{Error: err, Response: res}
		select {
		case <-ctx.Done():
			return
		case responses <- result:
		}
	}()

	return responses
}

func doFetchUserAlbumRequests(ctx context.Context, client *Client, userId int) <-chan Result {
	responses := make(chan Result)
	go func() {
		defer close(responses)
		res, _, err := client.Album.GetByUser(ctx, userId)
		result := Result{Error: err, Response: res}
		select {
		case <-ctx.Done():
			return
		case responses <- result:
		}
	}()

	return responses
}

func doFetchPostCommentsRequests(ctx context.Context, client *Client, userId int) <-chan Result {
	responses := make(chan Result)
	go func() {
		defer close(responses)
		res, _, err := client.Comment.GetByPost(ctx, userId)
		result := Result{Error: err, Response: res}
		select {
		case <-ctx.Done():
			return
		case responses <- result:
		}
	}()

	return responses
}

func doFetchPhotosByAlbumRequests(ctx context.Context, client *Client, albumId int) <-chan Result {
	responses := make(chan Result)
	go func() {
		defer close(responses)
		res, _, err := client.Album.GetPhotos(ctx, albumId)
		result := Result{Error: err, Response: res}
		select {
		case <-ctx.Done():
			return
		case responses <- result:
		}
	}()

	return responses
}

func doFetchUsersRequests(ctx context.Context, client *Client, urls []int) <-chan Result {
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

func doFetchTodoRequest(ctx context.Context, client *Client, urls []int) <-chan Result {
	responses := make(chan Result)
	go func() {
		defer close(responses)
		for _, todoId := range urls {
			res, _, err := client.Todo.Get(ctx, todoId)
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

func fetchAllUsers(ctx context.Context, client *Client) {
	userIdSize := 10

	var req = make([]int, 0)
	for i := 1; i <= userIdSize; i++ {
		req = append(req, i)
	}

	terminated := doFetchUsersRequests(ctx, client, req)

	// go func() {
	// 	// Cancel the operation after 1 second.
	// 	time.Sleep(1 * time.Second)
	// 	fmt.Println("Canceling doWork goroutine...")
	// 	cancel()
	// }()

	// 	// var result User
	// 	// json.NewDecoder(resp.Body).Decode(&result)
	// 	fmt.Println(string(b))
	// }

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

func fetchAllTodos(ctx context.Context, client *Client) {
	todoSize := 10

	var req = make([]int, 0)
	for i := 1; i <= todoSize; i++ {
		req = append(req, i)
	}

	terminated := doFetchTodoRequest(ctx, client, req)

	// go func() {
	// 	// Cancel the operation after 1 second.
	// 	time.Sleep(1 * time.Second)
	// 	fmt.Println("Canceling doWork goroutine...")
	// 	cancel()
	// }()

	// 	// var result User
	// 	// json.NewDecoder(resp.Body).Decode(&result)
	// 	fmt.Println(string(b))
	// }

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
