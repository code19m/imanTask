package service

import (
	"collector/internal/domain"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"

	"golang.org/x/sync/errgroup"
)

const (
	collectPostsUrl = "https://gorest.co.in/public/v1/posts"
)

var (
	mapjon = make(map[int32]struct{})
	m      = new(sync.Mutex)
)

func (s *service) CollectPosts(ctx context.Context, startPage int) error {
	g, ctx := errgroup.WithContext(ctx)

	for i := startPage; i < startPage+50; i++ {
		i := i

		g.Go(func() error {
			posts, err := downloadPosts(ctx, i)
			if err != nil {
				return err
			}

			for _, post := range posts {
				m.Lock()
				mapjon[int32(post.Id)] = struct{}{}
				m.Unlock()
			}

			return createOrUpdatePosts(ctx, posts, s.db)
		})
	}

	err := g.Wait()
	if err != nil {
		return err
	}

	return nil
}

func downloadPosts(ctx context.Context, page int) ([]domain.Post, error) {
	url := fmt.Sprintf("%s?page=%d", collectPostsUrl, page)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("invalid status code")
	}

	body, _ := io.ReadAll(resp.Body)

	var apiResponse apiResponse
	err = json.Unmarshal(body, &apiResponse)
	return apiResponse.Data, err
}

func createOrUpdatePosts(ctx context.Context, posts []domain.Post, db *sql.DB) error {
	stmt, err := db.Prepare(`
		INSERT OR REPLACE INTO posts (id, user_id, title, body)
		VALUES (?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, post := range posts {
		_, err = stmt.Exec(post.Id, post.UserID, post.Title, post.Body, post.UserID, post.Title, post.Body)
		if err != nil {
			return err
		}
	}

	return nil
}

type apiResponse struct {
	Data []domain.Post `json:"data"`
}
