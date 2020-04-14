package main

import (
	"context"
	"encoding/json"
	"fmt"
	restClient "github.com/rnkoaa/json-server/cmd/client/client"
	"github.com/rnkoaa/json-server/pkg/domain"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	client *restClient.Client
)

type Result struct {
	Error    error
	Response interface{}
}

func main() {
	count := 14
	ctx, cancel := context.WithCancel(context.Background())

	client = restClient.NewClient(http.DefaultClient, "")
	doWork(ctx, client, count)

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

func findUsers(ctx context.Context, client *restClient.Client, count int) []*domain.User {
	users := make([]*domain.User, 0)
	var req = make([]int, 0)
	for i := 1; i <= count; i++ {
		req = append(req, i)
	}

	terminated := doFetchUsersRequests(ctx, client, req)

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

		if i.Response != nil {
			if item, ok := i.Response.(*domain.User); ok {
				users = append(users, item)
			}
		}
	}
	return users
}

func printJson(v interface{}) {
	prettyJSON, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Fatal("Failed to generate json", err)
	}
	fmt.Printf("%s\n", string(prettyJSON))
}

func doWork(ctx context.Context, client *restClient.Client, count int) {
	var wg sync.WaitGroup
	defer wg.Wait()
	users := findUsers(ctx, client, count)
	for userIdx, u := range users {
		wg.Add(1)
		todos := fetchTodosByUser(ctx, &wg, client, u.ID)
		users[userIdx].Todos = todos

		wg.Add(1)
		albums := fetchAlbumsByUser(ctx, &wg, client, u.ID)
		for albumIdx, a := range albums {
			wg.Add(1)
			photos := fetchPhotosByAlbum(ctx, &wg, client, a.ID)
			albums[albumIdx].Photos = photos
		}
		users[userIdx].Albums = albums

		wg.Add(1)
		posts := fetchPostsByUser(ctx, &wg, client, u.ID)
		for postIdx, p := range posts {
			wg.Add(1)
			comments := fetchCommentsByPost(ctx, &wg, client, p.ID)
			posts[postIdx].Comments = comments
		}
		users[userIdx].Posts = posts
	}
	printJson(users)
}

func fetchTodosByUser(ctx context.Context, wg *sync.WaitGroup, c *restClient.Client, userId int) []*domain.Todo {
	defer wg.Done()
	results := make([]*domain.Todo, 0)

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
		if i.Response != nil {
			if item, ok := i.Response.([]*domain.Todo); ok {
				results = append(results, item...)
			}
		}

	}
	return results
}

func fetchPostsByUser(ctx context.Context, wg *sync.WaitGroup, c *restClient.Client, userId int) []*domain.Post {
	defer wg.Done()
	results := make([]*domain.Post, 0)
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
		if i.Response != nil {
			if item, ok := i.Response.([]*domain.Post); ok {
				results = append(results, item...)
			}
		}
	}
	return results
}

func fetchAlbumsByUser(ctx context.Context, wg *sync.WaitGroup, c *restClient.Client, userId int) []*domain.Album {
	defer wg.Done()
	results := make([]*domain.Album, 0)
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
		if i.Response != nil {
			if item, ok := i.Response.([]*domain.Album); ok {
				results = append(results, item...)
			}
		}
	}
	return results
}

func fetchCommentsByPost(ctx context.Context, wg *sync.WaitGroup, c *restClient.Client, postId int) []*domain.Comment {
	defer wg.Done()
	results := make([]*domain.Comment, 0)
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
		if i.Response != nil {
			if item, ok := i.Response.([]*domain.Comment); ok {
				results = append(results, item...)
			}
		}
	}
	return results
}

func fetchPhotosByAlbum(ctx context.Context, wg *sync.WaitGroup, c *restClient.Client, postId int) []*domain.Photo {
	defer wg.Done()
	results := make([]*domain.Photo, 0)
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
		if i.Response != nil {
			if item, ok := i.Response.([]*domain.Photo); ok {
				results = append(results, item...)
			}
		}
	}
	return results
}
