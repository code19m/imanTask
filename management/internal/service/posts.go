package service

import (
	"context"
	"database/sql"
	"fmt"
	"management/internal/domain"
)

func (s *service) GetPosts(ctx context.Context, offset int32, limit int32) ([]domain.Post, int32, error) {
	// Warning: Unsafe | Not production ready
	// Used string formatter instead of prepared statements because of
	// current sqlite driver doesn't support prepared statements for limit and offset keywords
	sqlStatement := fmt.Sprintf("SELECT id, user_id, title, body FROM posts LIMIT %d OFFSET %d", limit, offset)

	rows, err := s.db.QueryContext(ctx, sqlStatement)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var posts []domain.Post
	for rows.Next() {
		var post domain.Post
		if err := rows.Scan(&post.Id, &post.UserID, &post.Title, &post.Body); err != nil {
			return nil, 0, err
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	var count int32
	if err := s.db.QueryRowContext(ctx, "SELECT count(*) FROM posts").Scan(&count); err != nil {
		return nil, 0, err
	}

	return posts, count, nil
}

func (s *service) GetPost(ctx context.Context, id int32) (domain.Post, error) {
	var post domain.Post
	if err := s.db.QueryRowContext(ctx, "SELECT id, user_id, title, body FROM posts WHERE id = ?", id).
		Scan(&post.Id, &post.UserID, &post.Title, &post.Body); err != nil {

		if err == sql.ErrNoRows {
			return domain.Post{}, domain.ErrPostNotFound
		}
		return domain.Post{}, err
	}
	return post, nil
}

func (s *service) UpdatePost(ctx context.Context, payload domain.Post) error {
	_, err := s.db.ExecContext(ctx, "UPDATE posts SET user_id = ?, title = ?, body = ? WHERE id = ?",
		payload.UserID, payload.Title, payload.Body, payload.Id)
	return err
}

func (s *service) DeletePost(ctx context.Context, id int32) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM posts WHERE id = ?", id)
	return err
}
