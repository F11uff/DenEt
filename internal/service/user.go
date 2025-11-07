package service

import (
	"context"
	"errors"
	"denet/internal/model"
	"denet/internal/repository"
)

type UserService interface {
	CompleteTask(ctx context.Context, userID, taskID string) error
	SetReferrer(ctx context.Context, userID, referrerID string) error
	GetUserStatus(ctx context.Context, userID string) (*model.UserStatus, error)
	GetLeaderboard(ctx context.Context, limit int) ([]model.LeaderboardUser, error)
}

type userService struct {
	uow repository.UnitOfWork
}

func NewUserService(uow repository.UnitOfWork) UserService {
	return &userService{uow: uow}
}

func (s *userService) CompleteTask(ctx context.Context, userID, taskID string) error {
	user, err := s.uow.Users().GetByID(ctx, userID)
	if err != nil {
		return err
	}
	task, err := s.uow.Tasks().GetByID(ctx, taskID)
	if err != nil {
		return err
	}
	completed, err := s.uow.UserTasks().IsTaskCompleted(ctx, userID, taskID)
	if err != nil {
		return err
	}
	if completed {
		return errors.New("task already completed")
	}
	return s.uow.Transactions().WithTransaction(ctx, func(ctx context.Context) error {
		if err := s.uow.UserTasks().CompleteTask(ctx, userID, taskID); err != nil {
			return err
		}
		newBalance := user.Balance + task.Points
		return s.uow.Users().UpdateBalance(ctx, userID, newBalance)
	})
}

func (s *userService) SetReferrer(ctx context.Context, userID, referrerID string) error {
	if userID == referrerID {
		return errors.New("user cannot refer themselves")
	}
	return s.uow.Users().SetReferrer(ctx, userID, referrerID)
}

func (s *userService) GetUserStatus(ctx context.Context, userID string) (*model.UserStatus, error) {
	user, err := s.uow.Users().GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	completedTasks, err := s.uow.UserTasks().GetCompletedTasks(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &model.UserStatus{
		User:           user,
		CompletedTasks: completedTasks,
		TotalPoints:    user.Balance,
	}, nil
}

func (s *userService) GetLeaderboard(ctx context.Context, limit int) ([]model.LeaderboardUser, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	return s.uow.Users().GetLeaderboard(ctx, limit)
}