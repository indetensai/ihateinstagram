package usecases

import (
	"context"
	"phota/internal/entities"

	"github.com/google/uuid"
)

type postService struct {
	repo             entities.PostRepository
	UserService      entities.UserService
	FollowingService entities.FollowingService
}

func NewPostService(
	p entities.PostRepository,
	u entities.UserService,
	f entities.FollowingService,
) entities.PostService {
	return &postService{repo: p, UserService: u, FollowingService: f}
}

func (p *postService) Post(user_id uuid.UUID, ctx context.Context, desription string) (*uuid.UUID, error) {
	post_id, err := p.repo.CreatePost(user_id, ctx, desription)
	return post_id, err
}

func (p *postService) GetPost(
	post_id uuid.UUID,
	user_id uuid.UUID,
	ctx context.Context,
) (*entities.Post, error) {
	post, err := p.repo.GetPostByID(user_id, post_id, ctx)
	if err != nil {
		return nil, err
	}
	switch post.Visibility {
	case "followers":
		if !p.FollowingService.IsFollowing(post.UserID, user_id, context.Background()) || post.UserID != user_id {
			return nil, entities.ErrNotAuthorized
		}
	case "private":
		if post.UserID != user_id {
			return nil, entities.ErrNotAuthorized
		}
	}
	return post, nil
}

func (p *postService) ChangePost(content entities.ChangePostParams, ctx context.Context) error {
	post, err := p.GetPost(content.PostID, content.UserID, context.Background())
	if err != nil {
		return err
	}
	if post.UserID != content.UserID {
		return entities.ErrNotAuthorized
	}
	return p.repo.ChangePost(content, ctx)
}

func (p *postService) Like(post_id uuid.UUID, user_id uuid.UUID, ctx context.Context) error {
	post, err := p.GetPost(post_id, user_id, context.Background())
	if err != nil {
		return err
	}
	switch post.Visibility {
	case "followers":
		if !p.FollowingService.IsFollowing(post.UserID, user_id, context.Background()) || post.UserID != user_id {
			return entities.ErrNotAuthorized
		}
	case "private":
		if post.UserID != user_id {
			return entities.ErrNotAuthorized
		}
	}
	return p.repo.CreateLike(user_id, post_id, ctx)
}

func (p *postService) GetLikes(post_id uuid.UUID, user_id uuid.UUID, ctx context.Context) ([]entities.AppUser, error) {
	post, err := p.GetPost(post_id, user_id, context.Background())
	if err != nil {
		return nil, err
	}
	switch post.Visibility {
	case "followers":
		if !p.FollowingService.IsFollowing(post.UserID, user_id, context.Background()) || post.UserID != user_id {
			return nil, entities.ErrNotAuthorized
		}
	case "private":
		if post.UserID != user_id {
			return nil, entities.ErrNotAuthorized
		}
	}
	likes, err := p.repo.GetLikes(post_id, ctx)
	if err != nil {
		return nil, err
	}
	return likes, nil
}

func (p *postService) Unlike(user_id uuid.UUID, post_id uuid.UUID, ctx context.Context) error {
	post, err := p.GetPost(post_id, user_id, context.Background())
	if err != nil {
		return err
	}
	switch post.Visibility {
	case "followers":
		if !p.FollowingService.IsFollowing(post.UserID, user_id, context.Background()) || post.UserID != user_id {
			return entities.ErrNotAuthorized
		}
	case "private":
		if post.UserID != user_id {
			return entities.ErrNotAuthorized
		}
	}
	return p.repo.DeleteLike(user_id, post_id, ctx)
}

func (p *postService) DeletePost(post_id uuid.UUID, user_id uuid.UUID, ctx context.Context) error {
	post, err := p.GetPost(post_id, user_id, context.Background())
	if err != nil {
		return err
	}
	if post.UserID != user_id {
		return entities.ErrNotAuthorized
	}
	return p.repo.DeletePost(post_id, ctx)
}
