package usecases

import (
	"context"
	"phota/internal/entities"
	"phota/internal/usecases/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type postService struct {
	db *pgx.Conn
}

func (p *postService) Post(session_id string, ctx context.Context, desription string) (*uuid.UUID, error) {
	var user_service userService
	user_id, err := user_service.CheckSession(session_id, ctx)
	if err != nil {
		return nil, err
	}
	post_service := repository.NewPostRepository(p.db)
	post_id, err := post_service.Post(*user_id, ctx, desription)
	return post_id, err
}

func (p *postService) GettingPost(
	post_id uuid.UUID,
	session_id string,
	ctx context.Context,
) (*entities.Post, error) {
	var user_service userService
	user_id, err := user_service.CheckSession(session_id, ctx)
	if err != nil {
		return nil, err
	}
	post_service := repository.NewPostRepository(p.db)
	post, err := post_service.GettingPost(user_id, post_id, ctx)
	return post, nil
}
func (p *postService) PostChanging(
	visibility string,
	description string,
	post_id uuid.UUID,
	session_id string,
	ctx context.Context,
) error {
	var user_service userService
	user_id, err := user_service.CheckSession(session_id, ctx)
	if err != nil {
		return err
	}
	post_service := repository.NewPostRepository(p.db)
	err = post_service.PostChanging(user_id, visibility, description, post_id, ctx)
	return err
}
func (p *postService) Like(post_id uuid.UUID, session_id string, ctx context.Context) error {
	var user_service userService
	user_id, err := user_service.CheckSession(session_id, ctx)
	if err != nil {
		return err
	}
	post_service := repository.NewPostRepository(p.db)
	err = post_service.Like(user_id, post_id, ctx)
	return err
}
func (p *postService) GetLikes(session_id string, post_id uuid.UUID, ctx context.Context) (*[]uuid.UUID, error) {
	var user_service userService
	user_id, err := user_service.CheckSession(session_id, ctx)
	if err != nil {
		return nil, err
	}
	post_service := repository.NewPostRepository(p.db)
	likes, err := post_service.GetLikes(user_id, post_id, ctx)
	if err != nil {
		return nil, err
	}
	return likes, nil

}
func (p *postService) Unlike(session_id string, post_id uuid.UUID, ctx context.Context) error {
	var user_service userService
	user_id, err := user_service.CheckSession(session_id, ctx)
	if err != nil {
		return err
	}
	post_service := repository.NewPostRepository(p.db)
	err = post_service.Unlike(user_id, post_id, ctx)
	return err
}
