package usecases

import (
	"context"
	"phota/internal/entities"

	"github.com/google/uuid"
)

type postService struct {
	repo         entities.PostRepository
	user_service entities.UserService
}

func NewPostService(p entities.PostRepository, u entities.UserService) entities.PostService {
	return &postService{p, u}
}

func (p *postService) Post(user_id *uuid.UUID, ctx context.Context, desription string) (*uuid.UUID, error) {
	post_id, err := p.repo.Post(*user_id, ctx, desription)
	return post_id, err
}

func (p *postService) GettingPost(
	post_id uuid.UUID,
	user_id *uuid.UUID,
	ctx context.Context,
) (*entities.Post, error) {
	post, err := p.repo.GettingPost(user_id, post_id, ctx)
	if err != nil {
		return nil, err
	}
	return post, nil
}
func (p *postService) PostChanging(
	visibility string,
	description string,
	post_id uuid.UUID,
	user_id *uuid.UUID,
	ctx context.Context,
) error {
	err := p.repo.PostChanging(user_id, visibility, description, post_id, ctx)
	return err
}
func (p *postService) Like(post_id uuid.UUID, user_id *uuid.UUID, ctx context.Context) error {
	err := p.repo.Like(user_id, post_id, ctx)
	return err
}
func (p *postService) GetLikes(user_id *uuid.UUID, post_id uuid.UUID, ctx context.Context) (*[]uuid.UUID, error) {
	likes, err := p.repo.GetLikes(user_id, post_id, ctx)
	if err != nil {
		return nil, err
	}
	return likes, nil

}
func (p *postService) Unlike(user_id *uuid.UUID, post_id uuid.UUID, ctx context.Context) error {
	err := p.repo.Unlike(user_id, post_id, ctx)
	return err
}
