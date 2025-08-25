package posts

import (
	"context"
	"strings"
	"errors"
	"github.com/rs/zerolog/log"
	"database/sql"

	"forum/internal/model/posts"
)

func (r *repository) CreatePost(ctx context.Context, model posts.PostModel) error {
	query := `INSERT INTO posts(user_id, post_title, post_content, post_hashtags, created_at, updated_at, created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, model.UserID, model.PostTitle, model.PostContent, model.PostHashtags, model.CreatedAt, model.UpdatedAt, model.CreatedBy, model.UpdatedBy)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) GetAllPost(ctx context.Context, limit, offset int) (posts.GetAllPostResponse, error) {
	query := `SELECT p.id, p.user_id, u.username, p.post_title, p.post_content, p.post_hashtags 
	FROM posts p JOIN users u ON p.user_id = u.id ORDER BY p.updated_at DESC LIMIT ? OFFSET ?`

	response := posts.GetAllPostResponse{}
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return response, err
	}
	defer rows.Close()

	data := make([]posts.Post, 0)
	for rows.Next() {
		var (
			model    posts.PostModel
			username string
		)
		err = rows.Scan(&model.ID, &model.UserID, &username, &model.PostTitle, &model.PostContent, &model.PostHashtags)
		if err != nil {
			return response, err
		}
		data = append(data, posts.Post{
			ID:           model.ID,
			UserID:       model.UserID,
			Username:     username,
			PostTitle:    model.PostTitle,
			PostContent:  model.PostContent,
			PostHashtags: strings.Split(model.PostHashtags, ","),
		})
	}
	response.Data = data
	response.Pagination = posts.Pagination{
		Limit:  limit,
		Offset: offset,
	}
	return response, nil
}

func (r *repository) GetPostByID(ctx context.Context, id int64) (*posts.Post, error) {
    query := `
        SELECT p.id, p.user_id, u.username, p.post_title, p.post_content, p.post_hashtags
        FROM posts p 
        JOIN users u ON p.user_id = u.id 
        WHERE p.id = ?
    `

    var (
        model    posts.PostModel
        username string
    )
    row := r.db.QueryRowContext(ctx, query, id)

    err := row.Scan(&model.ID, &model.UserID, &username, &model.PostTitle, &model.PostContent, &model.PostHashtags)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            log.Warn().Msgf("GetPostByID: Post with ID %d not found", id)
            return nil, nil // Return nil if no rows are found
        }
        log.Error().Err(err).Msg("GetPostByID: Failed to fetch post from database")
        return nil, err
    }

    return &posts.Post{
        ID:           model.ID,
        UserID:       model.UserID,
        Username:     username,
        PostTitle:    model.PostTitle,
        PostContent:  model.PostContent,
        PostHashtags: strings.Split(model.PostHashtags, ","),
    }, nil
}

func (r *repository) UpdatePost(ctx context.Context, model posts.PostModel) error {
	query := `UPDATE posts SET post_title = ?, post_content = ?, post_hashtags = ?, updated_at = ?, updated_by = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, model.PostTitle, model.PostContent, model.PostHashtags, model.UpdatedAt, model.UpdatedBy, model.ID)
	if err != nil {
		return err
	}
	return nil
}