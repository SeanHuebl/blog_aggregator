// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: feeds.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const addFeed = `-- name: AddFeed :one
INSERT INTO feeds (id, name, url, user_id)
VALUES ($1, $2, $3, $4)
RETURNING id, name, url, created_at, updated_at, user_id, last_fetched_at
`

type AddFeedParams struct {
	ID     uuid.UUID
	Name   string
	Url    string
	UserID uuid.UUID
}

func (q *Queries) AddFeed(ctx context.Context, arg AddFeedParams) (Feed, error) {
	row := q.db.QueryRowContext(ctx, addFeed,
		arg.ID,
		arg.Name,
		arg.Url,
		arg.UserID,
	)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Url,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.LastFetchedAt,
	)
	return i, err
}

const getFeed = `-- name: GetFeed :one
SELECT id,
    name
FROM feeds
WHERE url = $1
`

type GetFeedRow struct {
	ID   uuid.UUID
	Name string
}

func (q *Queries) GetFeed(ctx context.Context, url string) (GetFeedRow, error) {
	row := q.db.QueryRowContext(ctx, getFeed, url)
	var i GetFeedRow
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const getFeeds = `-- name: GetFeeds :many
SELECT feeds.name,
    feeds.url,
    users.name
FROM feeds
    INNER JOIN users ON feeds.user_id = users.id
`

type GetFeedsRow struct {
	Name   string
	Url    string
	Name_2 string
}

func (q *Queries) GetFeeds(ctx context.Context) ([]GetFeedsRow, error) {
	rows, err := q.db.QueryContext(ctx, getFeeds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFeedsRow
	for rows.Next() {
		var i GetFeedsRow
		if err := rows.Scan(&i.Name, &i.Url, &i.Name_2); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getNextFeedToFetch = `-- name: GetNextFeedToFetch :one
SELECT id
FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1
`

func (q *Queries) GetNextFeedToFetch(ctx context.Context) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, getNextFeedToFetch)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const markFeedFetched = `-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
`

func (q *Queries) MarkFeedFetched(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, markFeedFetched, id)
	return err
}
